package cmd

import (
	"os"

	"github.com/spf13/cobra"
	"github.com/voidKandy/go-pyenv/pyenv"
)

var rootCmd = &cobra.Command{
	Use:     "prattl",
	Short:   "Prattl is a transcription tool",
	Long:    "A transcription tool built with Go and Python.\nComplete documentation is available at https://github.com/benleem/prattl",
	Version: "0.01",
}

func Execute(pyenv pyenv.PyEnv) {
	if err := rootCmd.Execute(); err != nil {
		os.Exit(1)
	}
}
