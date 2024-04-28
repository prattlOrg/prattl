#!/usr/bin/env python3

import torch
from transformers import WhisperProcessor, WhisperForConditionalGeneration
from datasets import Audio, load_dataset

def main ():
    # load model and processor
    processor = WhisperProcessor.from_pretrained("openai/whisper-tiny")
    model = WhisperForConditionalGeneration.from_pretrained("openai/whisper-tiny")
    forced_decoder_ids = processor.get_decoder_prompt_ids(language="english", task="transcribe")

    # load dummy dataset and read audio files
    ds = load_dataset("mozilla-foundation/common_voice_11_0", "en", streaming=True, trust_remote_code=True)
    ds = ds.cast_column("audio", Audio(sampling_rate=16_000))
    input_speech = next(iter(ds))["audio"]
    input_features = processor(input_speech["array"], sampling_rate=input_speech["sampling_rate"], return_tensors="pt").input_features 

    # generate token ids
    predicted_ids = model.generate(input_features, forced_decoder_ids=forced_decoder_ids)
    # decode token ids to text
    transcription = processor.batch_decode(predicted_ids)
    
    print(transcription)
if __name__ == "__main__":
    main()