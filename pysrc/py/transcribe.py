import sys
import torch
from transformers import AutoModelForSpeechSeq2Seq, AutoProcessor, pipeline

def transcribe (file_bytes) :
    device = 'cuda:0' if torch.cuda.is_available() else ('mps' if torch.backends.mps.is_available() else 'cpu')
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
    main()
