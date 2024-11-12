package cmd

import (
	"bufio"
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"os"
	"strings"

	"github.com/prattlOrg/prattl/internal/embed"
	"github.com/prattlOrg/prattl/internal/pysrc"

	"github.com/spf13/cobra"
)

func init() {
	RootCmd.AddCommand(transcribeCmd)
}

var transcribeCmd = &cobra.Command{
	Use:   "transcribe <filepath/to/audio.mp3>",
	Short: "Transcribe the provided audio file (file path)",
	Long:  "This command transcribes the provided audiofile and prints the resulting string",
	// Args:  cobra.MinimumNArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		isStdin, err := checkStdin()
		if err != nil {
			return err
		}
		if isStdin {
			fileBytes, err := readStdin()
			if err != nil {
				return err
			}
			err = transcribeStdin(fileBytes)
			if err != nil {
				return err
			}
		} else {
			if len(args) == 0 {
				return fmt.Errorf("requires at least 1 arg(s), only received 0")
			}
			err = transcribeFp(args...)
			if err != nil {
				return err
			}
		}

		return nil
	},
}

func checkStdin() (bool, error) {
	// fileStat, err := os.Stdin.Stat()
	// if err != nil {
	// 	return false, fmt.Errorf("getting stdin stat failed: %v", err)
	// }
	// // check if stdin is pipe
	// return fileStat.Mode()&os.ModeNamedPipe != 0, nil
	return true, nil
}

func readStdin() ([]byte, error) {
	reader := bufio.NewReader(os.Stdin)
	var fileBytes []byte
	for {
		b, err := reader.ReadByte()
		if err != nil && !errors.Is(err, io.EOF) {
			return nil, fmt.Errorf("failed to read file byte: %v", err)
		}
		// process the one byte b
		if err != nil {
			// end of file
			break
		}
		fileBytes = append(fileBytes, b)
	}
	return fileBytes, nil
}

func transcribeStdin(fileBytes []byte) error {
	fmt.Fprintln(os.Stderr, "Transcribing..")
	transcription, err := transcribe(fileBytes)
	if err != nil {
		return err
	}
	clearLine()
	fmt.Println(transcription[0])
	return nil
}

func transcribeFp(fps ...string) error {
	fmt.Fprintln(os.Stderr, "Transcribing..")
	var allBytes []byte
	for i, fp := range fps {
		fileBytes, err := os.ReadFile(fp)
		if err != nil {
			return fmt.Errorf("error reading file: %v", err)
		}
		allBytes = append(allBytes, fileBytes...)

		if i < len(fps)-1 {
			allBytes = append(allBytes, embed.CodeBytes...)
		}

	}
	transcriptionMap := make(map[string]string)
	transcriptions, err := transcribe(allBytes)
	if err != nil {
		return err
	}
	for i, trans := range transcriptions {
		transcriptionMap[fps[i]] = trans
	}
	jsonOutput, err := json.Marshal(transcriptionMap)
	if err != nil {
		return fmt.Errorf("marshaling to JSON failed: %v", err)
	}
	clearLine()
	_, err = io.WriteString(os.Stdout, string(jsonOutput)+"\n")
	if err != nil {
		return fmt.Errorf("writing to stdout failed: %v", err)
	}
	return nil
}

func transcribe(file []byte) ([]string, error) {
	fmt.Println(file)
	program, err := pysrc.ReturnFile("transcribe.py")
	if err != nil {
		return nil, err
	}

	env, err := pysrc.GetPrattlEnv()
	if err != nil {
		return nil, err
	}

	cmd, err := env.ExecutePython("-c", program)
	if err != nil {
		return nil, err
	}

	var out bytes.Buffer
	var stderr bytes.Buffer
	cmd.Stdout = &out
	cmd.Stderr = &stderr
	stdin, err := cmd.StdinPipe()
	if err != nil {
		return nil, fmt.Errorf("error instantiating pipe: %v", stderr.String())
	}
	if err = cmd.Start(); err != nil {
		return nil, fmt.Errorf("error starting command: %v", stderr.String())
	}
	_, err = stdin.Write(file)
	if err != nil {
		return nil, fmt.Errorf("error writing to stdin: %v", stderr.String())
	}
	stdin.Close()
	if err = cmd.Wait(); err != nil {
		return nil, fmt.Errorf("error waiting for command: %v", stderr.String())
	}

	output := out.String()
	returnStrings := strings.Split(strings.ToLower(output), embed.SeparatorExpectedString)

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
