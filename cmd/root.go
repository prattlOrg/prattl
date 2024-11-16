package cmd

import (
	"os"

	"github.com/spf13/cobra"
)

var (
	version = "dev"
	// commit  = "none"
	// date    = "unknown"
)

var RootCmd = &cobra.Command{
	Use:     "prattl",
	Short:   "Prattl is a transcription tool",
	Long:    "A transcription tool built with Go and Python.\nComplete documentation is available at https://github.com/prattlOrg/prattl",
	Version: version,
}

func Execute() {
	if err := RootCmd.Execute(); err != nil {
		os.Exit(1)
	}
}
