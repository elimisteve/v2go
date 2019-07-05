package cmd

import (
	"fmt"

	"github.com/spf13/cobra"
)

func init() {
	rootCmd.AddCommand(runCmd)
}

var runCmd = &cobra.Command{
	Use:     "run myprogram.v [ myprogram2.v ... ]",
	Example: "$ v2go run *.v",
	Short:   "Translate given .v file(s) to Go then run resulting Go binary",
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Printf("This will translate then run!\n")
	},
}
