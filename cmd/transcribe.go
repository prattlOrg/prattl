package cmd

import (
	"bytes"
	"fmt"
	"os"
	"prattl/pysrc"

	"github.com/spf13/cobra"
	"github.com/voidKandy/go-pyenv/pyenv"
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
	// fileInfo, err := os.Stat(fp)
	// if err != nil {
	// 	return "", fmt.Errorf("%s\n", err)
	// }
	// fileSize := fmt.Sprintf(", size: %.2fmb", float64(fileInfo.Size())/1048576)
	// return fileInfo.Name() + fileSize, nil

	fileBytes, err := os.ReadFile(fp)
	if err != nil {
		return "", err
	}
	program, err := pysrc.ReturnFile("transcribe.py")
	if err != nil {
		return "", err
	}

	// BAD!!!!
	home, err := os.UserHomeDir()
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
	env := pyenv.PyEnv{
		ParentPath: home + "/.prattl/",
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
		return "", fmt.Errorf("error writing to stdin: " + stderr.String())
	}
	stdin.Close()
	if err = cmd.Wait(); err != nil {
		return "", fmt.Errorf("error waiting for command: " + stderr.String())
	}
	output := out.String()
	return output, nil
}
