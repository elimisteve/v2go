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

var reStringInterpolation = regexp.MustCompile(`[^\\]\$\w+|[^\\]\$\{\w+\}`)

func TranslateVSource(in []byte) (out []byte, err error) {
	add := func(s string) { out = append(out, []byte(s)...) }

	lines := strings.Split(string(in), "\n")
	inComment := false
	start := true
	for _, l := range lines {
		l = strings.TrimSpace(l)
		if strings.HasPrefix(l, "//") {
			continue
		}
		if strings.HasPrefix(l, "/*") {
			inComment = true
			continue
		}
		if strings.HasSuffix(l, "*/") {
			inComment = false
			continue
		}
		if start && !inComment {
			if strings.HasPrefix(l, "module") {
				l = strings.Replace(l, "module", "package", 1)
				add(l + "\n")
				continue
			} else {
				add(`func main() {
`)
				defer func() {
					add("}")
				}()
			}
			start = false
		}

		if start {
			continue
		}

		if !start && !inComment && strings.HasPrefix(l, "import") {
			// Imports
		}

		if ndx := strings.Index(l, "println("); ndx != -1 {
			l = strings.Replace(l, "println('", `fmt.Printf("`, -1)
			l = strings.Replace(l, "')", `\n")`, -1)
			allVvars := reStringInterpolation.FindAllStringSubmatch(l, -1)
			for _, vvars := range allVvars {
				if len(vvars) == 0 {
					continue
				}
				vvar := vvars[0]
				varName := vvar[len(" $"):]
				if varName[0] == '{' {
					varName = varName[len("{") : len(varName)-len("}")]
				}
				// "some $name interp" -> "some " + fmt.Sprintf("%v", name) + " interp"
				l = strings.Replace(
					l,
					vvar,
					vvars[0][:1]+`" + fmt.Sprintf("%v", `+varName+`) + "`,
					-1)
			}
		}

		// TODO: Properly handle multi-line strings

		lb := []byte(l)

		for i := range lb {
			if lb[i] == '\'' {
				lb[i] = '"'
			} else if lb[i] == '`' {
				lb[i] = '\''
			}
		}

		add(string(lb) + "\n")
		if start {
			start = false
		}
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
