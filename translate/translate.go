// Steve Phillips / elimisteve
// 2019.07.05

package translate

import (
	"fmt"
	"io/ioutil"
	"os/exec"
	"regexp"
	"strings"
)

func cleanFileList(args []string) (vFilenames []string) {
	for _, s := range args {
		if strings.HasSuffix(s, ".v") && !strings.HasSuffix(s, "_test.v") {
			vFilenames = append(vFilenames, s)
		}
	}
	return vFilenames
}

func IsValidVFile(fname string) bool {
	return strings.HasSuffix(fname, ".v") && !strings.HasSuffix(fname, "_test.v")
}

var (
	reJsonDecode          = regexp.MustCompile(`(\w+) := json.decode\((.*?), (\w+)\) or {`)
	reStringInterpolation = regexp.MustCompile(`[^\\]\$[\w\.]+|[^\\]\$\{.+?\}`)
	reForIn               = regexp.MustCompile(`for (\w+) in (.*?) {`)
	reForInWithKey        = regexp.MustCompile(`for (\w+), (\w+) in (.*?) {`)
	reQuestion            = regexp.MustCompile(`(\w+) := (.*?)\?$`)
)

func quotesAndStringInterp(l, vvar string) string {
	if !strings.Contains(l, "$") {
		return l
	}
	varName := vvar[len(" $"):]
	if varName[0] == '{' {
		varName = varName[len("{") : len(varName)-len("}")]
	}
	// "some $name interp" -> "some " + fmt.Sprintf("%v", name) + " interp"
	repl := vvar[:1] + `" + fmt.Sprintf("%v", ` + varName + `) + "`
	l = strings.Replace(l, vvar, repl, -1)
	l = strings.Replace(l, "('", `("`, -1)
	l = strings.Replace(l, "')", `")`, -1)
	l = strings.Replace(l, `("" + `, `(`, -1)
	l = strings.Replace(l, ` + ""`, ``, -1)
	return l
}

func TranslateVSource(in []byte) (out []byte, err error) {
	add := func(s string) { out = append(out, []byte(s)...) }

	lines := strings.Split(string(in), "\n")
	inComment := false
	inQuotes := false
	specifiedPkg := false
	start := true
	importIndex := -1
	inMain := false
	skippedFnMain := false
	shouldMaybeDeferClosingBrace := false

	for _, l := range lines {
		justDidStringInterp := false
		l = strings.TrimRight(l, " \r\n")
		if strings.HasPrefix(l, "//") {
			add(l + "\n")
			continue
		}
		if strings.HasPrefix(l, "/*") {
			add(l + "\n")
			inComment = true
			continue
		}
		if strings.HasSuffix(l, "*/") {
			add(l + "\n")
			inComment = false
			continue
		}
		if start && !inComment {
			if !specifiedPkg && strings.HasPrefix(l, "module") {
				l = strings.Replace(l, "module", "package", 1)
				add(l + "\n")
				specifiedPkg = true
				continue
			} else {
				if !specifiedPkg {
					add(`package main

`)
					specifiedPkg = true
				}
				importIndex = len(out) - 1
				add("func main() {\n")
				inMain = true
				shouldMaybeDeferClosingBrace = true
			}
			start = false
		}

		if start {
			continue
		}

		if !inComment && strings.HasPrefix(l, "import") {
			l = strings.Replace(l, "'", `"`, -1)
			if !strings.Contains(l, `"`) {
				toImport := l[len("import "):]
				if toImport == "http" {
					toImport = "net/http"
				}
				if toImport == "json" {
					toImport = "encoding/json"
				}
				// Add double-quotes around imports
				l = fmt.Sprintf(`import "%s"`, toImport)
			}
			out = append(out[:importIndex], []byte("\n"+l+"\n"+string(out[importIndex:]))...)
			continue
		}

		if !inComment && strings.HasPrefix(l, "fn main(") && inMain {
			skippedFnMain = true
			continue
		}

		if strings.HasPrefix(l, "struct") {
			typ := l[len("struct ") : len(l)-2]
			add(fmt.Sprintf("type %s struct {\n", typ))
			continue
		}

		if strings.HasPrefix(l, "fn ") {
			l = strings.Replace(l, "fn ", "func ", 1)
		}

		//
		// Generic logic
		//

		jsonDecode := reJsonDecode.FindAllStringSubmatch(l, -1)
		if len(jsonDecode) > 0 {
			newVar := jsonDecode[0][1]
			typ := strings.Replace(jsonDecode[0][2], " ", "", -1)
			data := jsonDecode[0][3]
			add(fmt.Sprintf("%s := *new(%s)\n", newVar, typ))
			add(fmt.Sprintf("if err := json.Unmarshal([]byte(%s), &%s); err != nil {\n", data, newVar))
			continue
		}

		allVvars := reStringInterpolation.FindAllStringSubmatch(l, -1)
		for _, vvars := range allVvars {
			if len(vvars) == 0 {
				continue
			}
			vvar := vvars[0]
			l = quotesAndStringInterp(l, vvar)
			justDidStringInterp = true
		}

		if !inComment && strings.Contains(l, "mut ") {
			l = strings.Replace(l, "mut ", "", -1)
		}

		allForIns := reForIn.FindAllStringSubmatch(l, -1)
		if len(allForIns) > 0 {
			for _, forIn := range allForIns {
				add(fmt.Sprintf("\tfor _, %s := range %s {\n", forIn[1], forIn[2]))
			}
			continue
		}

		allForInsWithKey := reForInWithKey.FindAllStringSubmatch(l, -1)
		if len(allForInsWithKey) > 0 {
			for _, forIn := range allForInsWithKey {
				add(fmt.Sprintf("\tfor %s, %s := range %s {\n", forIn[1], forIn[2], forIn[3]))
			}
			continue
		}

		allQs := reQuestion.FindAllStringSubmatch(l, -1)
		if len(allQs) > 0 {
			for _, q := range allQs {
				add(fmt.Sprintf("\t%s, err := %s\n\tif err != nil {\n\t\tpanic(err)\n\t}\n", q[1],
					quotesAndStringInterp(q[2], q[2])))
			}
			continue
		}

		l = strings.Replace(l, "print(", "fmt.Print(", -1)
		l = strings.Replace(l, "eprintln(", `fmt.Fprintf(os.Stderr, `, -1)
		l = strings.Replace(l, "exit(", "os.Exit(", -1)

		if ndx := strings.Index(l, "println("); ndx != -1 {
			l = strings.Replace(l, "println('", `fmt.Printf("`, -1)
			l = strings.Replace(l, "')", `\n")`, -1)
			l = strings.Replace(l, "println(", "fmt.Println(", -1)
		}


		// TODO: Properly handle multi-line strings

		lb := []byte(l)
		var tweakedlb []byte

		for i := range lb {
			last := i == len(lb)-1
			if lb[i] == '\'' || lb[i] == '"' {
				if inQuotes && !last && !justDidStringInterp &&
					lb[i+1] != ')' && lb[i+1] != ',' {

					tweakedlb = append(tweakedlb, '\\', '"')
				} else {
					tweakedlb = append(tweakedlb, '"')
					inQuotes = !inQuotes
				}
			} else if lb[i] == '`' {
				tweakedlb = append(tweakedlb, '\'')
			} else {
				tweakedlb = append(tweakedlb, lb[i])
			}
		}

		add(string(tweakedlb) + "\n")
		if start {
			start = false
		}
		justDidStringInterp = false
	}
	if shouldMaybeDeferClosingBrace && !skippedFnMain {
		out = append(out, '}')
	}
	out = append(out, '\n')
	return out, nil
}

func TranslateVFile(vfile string) (goFilename string, err error) {
	goFilename = vfile[:len(vfile)-1] + "go"
	in, err := ioutil.ReadFile(vfile)
	if err != nil {
		return "", err
	}
	out, err := TranslateVSource(in)
	if err != nil {
		return "", err
	}
	if err := ioutil.WriteFile(goFilename, out, 0644); err != nil {
		return goFilename, err
	}
	goimportsPath, err := exec.LookPath("goimports")
	if err != nil {
		fmt.Printf("WARNING: Couldn't find goimports\n")
		return goFilename, nil
	}
	err = exec.Command(goimportsPath, []string{"-w", goFilename}...).Run()
	if err != nil {
		return goFilename, err
	}
	return goFilename, nil
}

func TranslateVFiles(vFilenames []string) (goFilenames []string, err error) {
	goFilenames = make([]string, 0, len(vFilenames))
	for _, fname := range vFilenames {
		if !IsValidVFile(fname) {
			fmt.Printf("WARNING: Ignoring file %s\n", fname)
			continue
		}
		goFilename, err := TranslateVFile(fname)
		if err != nil {
			return goFilenames, err
		}
		goFilenames = append(goFilenames, goFilename)
	}
	return goFilenames, nil
}
