// Steve Phillips / elimisteve
// 2019.07.05

package translate

import (
	"fmt"
	"os/exec"
)

func TranslateAndRunFiles(args []string) (out []byte, err error) {
	goFiles, err := TranslateVFiles(args)
	if err != nil {
		return nil, fmt.Errorf("Error translating .v files to .go: %v\n", err)
	}
	runArgs := append([]string{"run"}, goFiles...)
	out, err = exec.Command("go", runArgs...).Output()
	if err != nil {
		return nil, fmt.Errorf("Error running translated Go files: %v\n", err)
	}
	return out, nil
}
