package cmd

import (
	"bytes"
	"fmt"
	"os"
	"time"

	"github.com/benleem/prattl/pysrc"
	"github.com/briandowns/spinner"
	"github.com/spf13/cobra"
)

func init() {
	rootCmd.AddCommand(transcribeCmd)
}

var transcribeCmd = &cobra.Command{
	Use:   "transcribe <filepath/to/audio.mp3>",
	Short: "Transcribe the provided audio file (file path)",
	Long:  `This command transcribes the provided audiofile and prints the resulting string`,
	Args:  cobra.MaximumNArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		if len(args) == 0 {
			return fmt.Errorf("%s", "no file path provided\n")
		}
		transcription, err := transcribe(args[0])
		if err != nil {
			return err
		}
		fmt.Println(transcription)
		return nil
	},
}

func transcribe(fp string) (string, error) {
	fileBytes, err := os.ReadFile(fp)

	s := spinner.New(spinner.CharSets[35], 100*time.Millisecond, spinner.WithWriter(os.Stderr))
	s.Prefix = "transcribing: "
	s.Suffix = "\n"
	s.Start()

	if err != nil {
		return "", err
	}
	program, err := pysrc.ReturnFile("transcribe.py")
	if err != nil {
		return "", err
	}

	env, err := pysrc.PrattlEnv()
	if err != nil {
		fmt.Printf("Error getting prattl env: %v\n", err)
		os.Exit(1)
	}
	cmd := env.ExecutePython("-c", program)
	var out bytes.Buffer
	var stderr bytes.Buffer
	stdin, err := cmd.StdinPipe()
	if err != nil {
		return "", fmt.Errorf("error instantiating pipe: " + err.Error())
	}
	cmd.Stdout = &out
	cmd.Stderr = &stderr
	if err = cmd.Start(); err != nil {
		return "", fmt.Errorf("error starting command: " + err.Error())
	}
	_, err = stdin.Write(fileBytes)
	if err != nil {
		return "", fmt.Errorf("error writing to stdin: " + err.Error())
	}
	stdin.Close()
	if err = cmd.Wait(); err != nil {
		return "", fmt.Errorf("error waiting for command: " + err.Error())
	}
	s.Stop()
	output := out.String()
	return output, nil
}
