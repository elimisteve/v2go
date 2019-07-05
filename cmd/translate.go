package cmd

import (
	"fmt"

	"github.com/spf13/cobra"
)

func init() {
	rootCmd.AddCommand(translateCmd)
}

var translateCmd = &cobra.Command{
	Use:   "translate myprogram.v [ myprogram2.v ... ]",
	Short: "Translate given .v file(s) to Go  [default]",
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Printf("This will _just_ translate!\n")
	},
}
