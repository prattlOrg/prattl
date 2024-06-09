package cmd

import (
	"fmt"

	"github.com/spf13/cobra"
)

func init() {
	rootCmd.AddCommand(versionCmd)
}

var versionCmd = &cobra.Command{
	Use:   "version",
	Short: "Print the version number of Prattl",
	Long:  `All software has versions. This is Prattl's`,
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println("prattl v0.01 -- HEAD")
	},
}
