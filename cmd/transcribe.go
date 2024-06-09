package cmd

import (
	"fmt"
	"os"

	"github.com/kluctl/go-embed-python/python"
	"github.com/spf13/cobra"
)

func init() {
	rootCmd.AddCommand(transcribeCmd)
}

var transcribeCmd = &cobra.Command{
	Use:   "transcribe",
	Short: "Transcribe the provided audio file (file path)",
	Long:  `This command transcribes the provided audiofile and prints the resulting string`,
	Args:  cobra.MaximumNArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		if len(args) == 0 {
			return fmt.Errorf("No file path provided")
		}
		err := transcribe(args[0])
		if err != nil {
			return err
		}
		return nil

	},
}

func transcribe(fp string) error {
	ep, err := python.NewEmbeddedPython("transcribe")
	if err != nil {
		return err
	}
	cmd := ep.PythonCmd("-c", "print('Transcribing audio file')")
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	err = cmd.Run()
	if err != nil {
		return err
	}
	return nil
}
