package cmd

import (
	"bytes"
	"fmt"
	"os"

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
	fileBytes, err := os.ReadFile(fp) //read the content of file
	if err != nil {
		return "", err
	}

	ep, err := python.NewEmbeddedPython("transcribe")
	if err != nil {
		return "", err
	}

	cmd := ep.PythonCmd("-c", `import sys
import torch
from transformers import AutoModelForSpeechSeq2Seq, AutoProcessor, pipeline

def transcribe (file_bytes) :
    # device = "cuda:0" if torch.cuda.is_available() else "cpu"
    device = torch.device("mps" if torch.backends.mps.is_available() else "cpu")
    torch_dtype = torch.float16 if torch.cuda.is_available() else torch.float32

    model_id = "distil-whisper/distil-medium.en"
    model = AutoModelForSpeechSeq2Seq.from_pretrained(
        model_id, torch_dtype=torch_dtype, low_cpu_mem_usage=True, use_safetensors=True
    )
    model.to(device)
    processor = AutoProcessor.from_pretrained(model_id)
    
    pipe = pipeline(
        "automatic-speech-recognition",
        model=model,
        tokenizer=processor.tokenizer,
        feature_extractor=processor.feature_extractor,
        max_new_tokens=128,
        torch_dtype=torch_dtype,
        device=device,
    )
    result = pipe(file_bytes, return_timestamps=True)
    sys.stdout.write(result["text"])

def main ():
    file_bytes = sys.stdin.buffer.read()
    transcribe(file_bytes)

if __name__ == "__main__":
    main()`)

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
