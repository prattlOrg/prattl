package transcribe

import (
	"bytes"
	"fmt"
	"os/exec"
)

func TranscribeLocal() {
	// cmd := exec.Command("ffmpeg")
	cmd := exec.Command("python3", "transcribe/transcribe.py")
	var out bytes.Buffer
	var stderr bytes.Buffer
	cmd.Stdout = &out
	cmd.Stderr = &stderr
	err := cmd.Run()
	if err != nil {
		fmt.Println(fmt.Sprint(err) + ": " + stderr.String())
		return
	}
	fmt.Println("Result: " + out.String())
}
