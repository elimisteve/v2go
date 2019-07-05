package cmd

import (
	"fmt"
	"log"

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
		out, err := translate.TranslateAndRunFiles(args)
		if err != nil {
			log.Fatalf("Error translating and running files: %v\n", err)
		}
		fmt.Printf("%s", out)
	},
}
