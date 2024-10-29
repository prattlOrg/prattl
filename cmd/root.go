package cmd

import (
	"os"

	"github.com/spf13/cobra"
)

var rootCmd = &cobra.Command{
	Use:     "prattl",
	Short:   "Prattl is a transcription tool",
	Long:    "A transcription tool built with Go and Python.\nComplete documentation is available at https://github.com/prattlOrg/prattl",
	Version: "v0.2.11",
}

func Execute() {
	if err := rootCmd.Execute(); err != nil {
		os.Exit(1)
	}
}
