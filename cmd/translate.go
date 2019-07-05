package cmd

import (
	"log"

	"github.com/elimisteve/v2go/translate"
	"github.com/spf13/cobra"
)

func init() {
	rootCmd.AddCommand(translateCmd)
}

var translateCmd = &cobra.Command{
	Use:   "translate myprogram.v [ myprogram2.v ... ]",
	Short: "Translate given .v file(s) to Go  [default]",
	Run: func(cmd *cobra.Command, args []string) {
		_, err := translate.TranslateVFiles(args)
		if err != nil {
			log.Fatalf("Error translating files: %v\n", err)
		}
	},
}
