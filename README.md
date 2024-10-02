# prattl

CLI tool for transcription


## Build Locally

1. Clone repository
2. Either run `go build` and then manually move the resulting binary into your `$PATH`, or just run `go install` to move the executable file to $GOPATH/bin
3. You're good to go!

## Prattl manages its own python distribution
When you run `prattl prepare` prattl installs a python distribution specific to your OS to your system, This means you don't need to manage Python dependencies or risk corrupting your existing environments.
`prattl clean` will completely remove this distribution. 

## Transcribing
Under the hood, prattl is using [distil-whisper](https://huggingface.co/distil-whisper/distil-large-v3), which runs locally on your system. The better your GPU, the faster the transcription, if you do not have a GPU, it will use your CPU. To create a transcription, use the command:

`prattl transcribe <filepath>`


You can provide multiple file paths, and prattl will transcribe all of them as a single batch. This means the efficiency increases with more files!

Upon completion, the output will be a JSON object. For instance, if you run:

`prattl transcribe test1.mp3 test2.mp3 test3.mp3`

the output will be:
```json
{
    "test1.mp3": "test1.mp3's transcription",
    "test2.mp3": "test2.mp3's transcription",
    "test3.mp3": "test3.mp3's transcription",
}
```

### OS/CPU Architecture Support
- windows/386
- windows/amd64
- darwin/arm64
- darwin/amd64
- linux-gnu/arm64
- linux-gnu/amd64

> **_NOTE:_** CUDA architecture GPUs can take advantage of GPU acceleration for transcription

### Prerequisites
[ffmpeg](https://www.ffmpeg.org/) installed and included in `$PATH`
