package cmd

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"os"
	"strings"

	"github.com/benleem/prattl/embed"
	"github.com/benleem/prattl/pysrc"

	"github.com/spf13/cobra"
)

func init() {
	rootCmd.AddCommand(transcribeCmd)
}

var transcribeCmd = &cobra.Command{
	Use:   "transcribe <filepath/to/audio.mp3>",
	Short: "Transcribe the provided audio file (file path)",
	Long:  `This command transcribes the provided audiofile and prints the resulting string`,
	Args:  cobra.MinimumNArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		if len(args) == 0 {
			return fmt.Errorf("%s", "no file path provided\n")
		}

		transcriptionMap := make(map[string]string)
		transcriptions, err := transcribe(args)
		if err != nil {
			return err
		}

		for i, trans := range transcriptions {
			transcriptionMap[args[i]] = trans
			// _, err := io.WriteString(os.Stdout, trans+"\n")
			// if err != nil {
			// 	return fmt.Errorf("error writing to stdout: %v", err)
			// }
		}

		jsonOutput, err := json.Marshal(transcriptionMap)
		if err != nil {
			return fmt.Errorf("error marshaling to JSON: %v", err)
		}

		_, err = io.WriteString(os.Stdout, string(jsonOutput)+"\n")
		if err != nil {
			return fmt.Errorf("error writing to stdout: %v", err)
		}

		return nil
	},
}

func transcribe(fps []string) ([]string, error) {
	returnStrings := []string{}

	var allBytes []byte
	for i, fp := range fps {
		fileBytes, err := os.ReadFile(fp)
		if err != nil {
			return returnStrings, err
		}
		allBytes = append(allBytes, fileBytes...)

		if i < len(fps)-1 {
			allBytes = append(allBytes, embed.CodeBytes...)
		}

	}
	program, err := pysrc.ReturnFile("transcribe.py")
	if err != nil {
		return returnStrings, err
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
		return returnStrings, fmt.Errorf("error instantiating pipe: " + err.Error())
	}
	cmd.Stdout = &out
	cmd.Stderr = &stderr
	if err = cmd.Start(); err != nil {
		return returnStrings, fmt.Errorf("error starting command: " + err.Error())
	}
	_, err = stdin.Write(allBytes)
	if err != nil {
		return returnStrings, fmt.Errorf("error writing to stdin: " + stderr.String())
	}
	stdin.Close()
	if err = cmd.Wait(); err != nil {
		return returnStrings, fmt.Errorf("error waiting for command: " + stderr.String())
	}

	output := out.String()

	returnStrings = strings.Split(strings.ToLower(output), embed.SeparatorExpectedString)

	for _, str := range returnStrings {
		str = fmt.Sprintf("---%s---\n", str)
	}

	return returnStrings, nil
}
