package cmd

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"os"
	"strings"

	"github.com/prattlOrg/prattl/internal/embed"
	"github.com/prattlOrg/prattl/internal/pysrc"

	"github.com/spf13/cobra"
)

func init() {
	rootCmd.AddCommand(transcribeCmd)
}

var transcribeCmd = &cobra.Command{
	Use:   "transcribe <filepath/to/audio.mp3>",
	Short: "Transcribe the provided audio file (file path)",
	Long:  "This command transcribes the provided audiofile and prints the resulting string",
	Args:  cobra.MinimumNArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {

		if len(args) == 0 {
			return fmt.Errorf("%s", "no file path provided\n")
		}
		fmt.Fprintln(os.Stderr, "Transcribing..")


		transcriptionMap := make(map[string]string)
		transcriptions, err := transcribe(args)
		if err != nil {
			return err
		}

		for i, trans := range transcriptions {
			transcriptionMap[args[i]] = trans

		}

		jsonOutput, err := json.Marshal(transcriptionMap)
		if err != nil {
			return fmt.Errorf("error marshaling to JSON: %v", err)
		}

		clearLine()
		_, err = io.WriteString(os.Stdout, string(jsonOutput)+"\n")
		if err != nil {
			return fmt.Errorf("error writing to stdout: %v", err)
		}

		return nil
	},
}

func isPipeInput() bool {
	fileInfo, _ := os.Stdin.Stat()
	return fileInfo.Mode()&os.ModeCharDevice == 0
}

func transcribe(fps []string) ([]string, error) {
	returnStrings := []string{}

	var allBytes []byte
	for i, fp := range fps {
		fileBytes, err := os.ReadFile(fp)
		if err != nil {
			return returnStrings, fmt.Errorf("error reading file: %v", err)
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

	env, err := pysrc.GetPrattlEnv()
	if err != nil {
		return returnStrings, err
	}
	cmd, err := env.ExecutePython("-c", program)
	if err != nil {
		return returnStrings, err
	}
	var out bytes.Buffer
	var stderr bytes.Buffer
	stdin, err := cmd.StdinPipe()
	if err != nil {
		return returnStrings, fmt.Errorf("error instantiating pipe: %v", err)
	}
	cmd.Stdout = &out
	cmd.Stderr = &stderr
	if err = cmd.Start(); err != nil {
		return returnStrings, fmt.Errorf("error starting command: %v", err)
	}
	_, err = stdin.Write(allBytes)
	if err != nil {
		return returnStrings, fmt.Errorf("error writing to stdin: %v", err)
	}
	stdin.Close()
	if err = cmd.Wait(); err != nil {
		return returnStrings, fmt.Errorf("error waiting for command: %v", err)
	}

	output := out.String()
	// fmt.Println(output)

	returnStrings = strings.Split(strings.ToLower(output), embed.SeparatorExpectedString)

	for _, str := range returnStrings {
		str = fmt.Sprintf("---%s---\n", str)
	}

	return returnStrings, nil
}

func clearLine() {
	const clear = "\033[2K"
	// clear line
	fmt.Printf(clear)
	// the line is cleared but the cursor is in the wrong place. the carriage
	// return moves the cursor to the beginning of the line.
	fmt.Printf("\r")
}
