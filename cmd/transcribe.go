package cmd

import (
	"bytes"
	"fmt"
	"os"
	"prattl/internal/python-libs/data"
	pysrc "prattl/python-src"

	"github.com/kluctl/go-embed-python/embed_util"
	"github.com/kluctl/go-embed-python/python"
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
			return fmt.Errorf("no file path provided\n")
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
	if err != nil {
		return "", err
	}

	tmpDir := "prattl-embedded"
	ep, err := python.NewEmbeddedPythonWithTmpDir(tmpDir+"-python", true)
	if err != nil {
		return "", err
	}
	prattlLib, err := embed_util.NewEmbeddedFilesWithTmpDir(data.Data, tmpDir+"-lib", true)
	if err != nil {
		return "", err
	}
	ep.AddPythonPath(prattlLib.GetExtractedPath())

	py, err := pysrc.ReturnSrc()
	if err != nil {
		return "", err
	}
	cmd := ep.PythonCmd("-c", py)

	var out bytes.Buffer
	var stderr bytes.Buffer

	stdin, err := cmd.StdinPipe()
	if err != nil {
		return "", fmt.Errorf(stderr.String())
	}
	cmd.Stdout = &out
	cmd.Stderr = &stderr
	if err = cmd.Start(); err != nil {
		return "", fmt.Errorf(stderr.String())
	}
	_, err = stdin.Write(fileBytes)
	if err != nil {
		return "", fmt.Errorf(stderr.String())
	}
	stdin.Close()
	if err = cmd.Wait(); err != nil {
		return "", fmt.Errorf(stderr.String())
	}
	output := out.String()
	return output, nil
}
