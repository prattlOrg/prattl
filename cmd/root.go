package cmd

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"
)

var RootCmd = &cobra.Command{
	Use:   "prattl",
	Short: "Prattl is a transcription tool",
	Long:  "A transcription tool built with Go and Python.\nComplete documentation is available at https://github.com/prattlOrg/prattl",
}

func Execute(version string) {
	if version == "dev" {
		RootCmd.Version = version
	} else {
		RootCmd.Version = fmt.Sprintf("v%v", version)
	}
	if err := RootCmd.Execute(); err != nil {
		os.Exit(1)
	}
}
