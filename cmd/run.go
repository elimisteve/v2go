package cmd

import (
	"fmt"
	"log"
	"os/exec"

	"github.com/elimisteve/v2go/translate"
	"github.com/spf13/cobra"
)

func init() {
	rootCmd.AddCommand(runCmd)
}

var runCmd = &cobra.Command{
	Use:     "run myprogram.v [ myprogram2.v ... ]",
	Example: "$ v2go run test_v_files/hello_world.v",
	Short:   "Translate given .v file(s) to Go then run resulting Go binary",
	Run: func(cmd *cobra.Command, args []string) {
		goFiles, err := translate.TranslateVFiles(args)
		if err != nil {
			log.Fatalf("Error translating files: %v\n", err)
		}
		runArgs := append([]string{"run"}, goFiles...)
		out, err := exec.Command("go", runArgs...).Output()
		if err != nil {
			log.Fatalf("Error running translated Go files: %v\n", err)
		}
		fmt.Printf("%s", out)
	},
}
