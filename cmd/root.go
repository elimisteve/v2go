package cmd

import (
	"fmt"
	"log"
	"os"
	"os/exec"

	"github.com/elimisteve/v2go/translate"
	"github.com/spf13/cobra"
)

var rootCmd = &cobra.Command{
	Use:   "v2go",
	Args:  cobra.ArbitraryArgs,
	Short: "v2go translates V code to Go code",
	Long: `v2go is a V-to-Go translator/transpiler/compiler.

To learn more about V, the best-designed programming language of all time, visit <https://vlang.io>.`,
	Run: func(cmd *cobra.Command, args []string) {
		if len(args) == 0 {
			fmt.Printf("ERROR: Specify 1 or more .v files to translate.\n\n")
			cmd.UsageFunc()(cmd)
			os.Exit(1)
		}

		// Translate
		goFilenames, err := translate.TranslateVFiles(args)
		if err != nil {
			log.Fatalf("Erorr translating .v files: %v\n", err)
		}

		// $ go run <all_translated_files>
		goRunArgs := append([]string{"run"}, goFilenames...)
		if err := exec.Command("go", goRunArgs...).Run(); err != nil {
			log.Fatalf("Error running translated Go files `%#v` -- %v\n",
				goFilenames, err)
		}
	},
}

func Execute() {
	if err := rootCmd.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}
