package ffmpeg

import (
	"bytes"
	"os/exec"
)

func CheckInstall() error {
	args := [1]string{"-h"}
	cmd := exec.Command("ffmpeg", args[:]...)
	var out bytes.Buffer
	var stderr bytes.Buffer
	cmd.Stdout = &out
	cmd.Stderr = &stderr
	if err := cmd.Start(); err != nil {
		return err
	}
	if err := cmd.Wait(); err != nil {
		return err
	}
	// output := out.String()
	return nil
}
