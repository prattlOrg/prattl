package transcribe

import (
	"bytes"
	"fmt"
	"log"
	"os/exec"
	"time"
)

//// Logging

func trace(s string) (string, time.Time) {
	log.Println("START:", s)
	return s, time.Now()
}

func un(s string, startTime time.Time) {
	endTime := time.Now()
	log.Println("END:", s, "took", endTime.Sub(startTime))
}

func TranscribeLocal(base64 string) (transcription string, err error) {
	defer un(trace("transcribe.TranscribeLocal()"))
	var out bytes.Buffer
	var stderr bytes.Buffer

	cmd := exec.Command("python3", "transcribe/transcribe.py")
	stdin, err := cmd.StdinPipe()
	if err != nil {
		return "", err
	}

	cmd.Stdout = &out
	cmd.Stderr = &stderr

	if err = cmd.Start(); err != nil {
		return "", err
	}

	_, err = stdin.Write([]byte(base64))
	if err != nil {
		return "", err
	}
	stdin.Close()

	if err = cmd.Wait(); err != nil {
		fmt.Println(fmt.Sprint(err) + ": " + stderr.String())
		return "", err
	}

	output := out.String()
	return output, nil
}
