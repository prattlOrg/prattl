package ffmpeg

import (
	"os/exec"
)

func CheckInstall() error {
	cmd := exec.Command("ffmpeg", "-h")
	if err := cmd.Run(); err != nil {
		return err
	}
	return nil
}
