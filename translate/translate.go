// Steve Phillips / elimisteve
// 2019.07.05

package translate

import (
	"fmt"
	"io/ioutil"
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

func TranslateVSource(in []byte) (out []byte, err error) {
	lines := strings.Split(string(in), "\n")
	for _, vline := range lines {
		goline := vline
		// TODO: add actual translation functionality!
		out = append(out, []byte(goline)...)
	}
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
	err = ioutil.WriteFile(goFilename, out, 0644)
	return goFilename, err
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
