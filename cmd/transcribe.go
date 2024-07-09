package cmd

import (
	"fmt"

	"github.com/spf13/cobra"
	pyenv "github.com/voidKandy/go-pyenv/pyenv"
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
			return fmt.Errorf("no file path provided")
		}
		err := transcribe(args[0])
		if err != nil {
			return err
		}
		return nil
	},
}

func transcribe(fp string) error {
	env := pyenv.DefaultPyEnv()

	args := []string{
		"-c",
		`import base64
		import sys
		# import torch
		# from transformers import AutoModelForSpeechSeq2Seq, AutoProcessor, pipeline

		def transcribe (fp) :
			sys.stdout.write(fp)
		#     device = "cuda:0" if torch.cuda.is_available() else "cpu"
		#     # device = torch.device("mps" if torch.backends.mps.is_available() else "cpu")
		#     torch_dtype = torch.float16 if torch.cuda.is_available() else torch.float32

		#     model_id = "distil-whisper/distil-medium.en"
		#     model = AutoModelForSpeechSeq2Seq.from_pretrained(
		#         model_id, torch_dtype=torch_dtype, low_cpu_mem_usage=True, use_safetensors=True
		#     )
		#     model.to(device)
		#     processor = AutoProcessor.from_pretrained(model_id)
			
		#     pipe = pipeline(
		#         "automatic-speech-recognition",
		#         model=model,
		#         tokenizer=processor.tokenizer,
		#         feature_extractor=processor.feature_extractor,
		#         max_new_tokens=128,
		#         torch_dtype=torch_dtype,
		#         device=device,
		#     )
		#     result = pipe(base64_bytes, return_timestamps=True)
		#     print(result["text"])

		def main ():
			fp = sys.stdin.read().strip()
			transcribe(fp)

		if __name__ == "__main__":
			main()`,
	}
	stdout, err := env.ExecutePython(args)
	// stdin, err := cmd.StdinPipe()
	if err != nil {
		return err
	}
	fmt.Println(stdout)
	return nil
	// cmd.Stdout = os.Stdout
	// cmd.Stderr = os.Stderr
	// err = cmd.Run()
	// if err != nil {
	// 	return err
	// }
}
