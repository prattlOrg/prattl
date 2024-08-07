package ffmpeg

import (
	"os/exec"
)

func CheckInstall() error {
	cmd := exec.Command("ffmpeg", "-h")
	if err := cmd.Run(); err != nil {
		return err
	}
	// var out bytes.Buffer
	// var stderr bytes.Buffer
	// cmd.Stdout = &out
	// cmd.Stderr = &stderr
	// if err := cmd.Start(); err != nil {
	// 	return err
	// }
	// if err := cmd.Wait(); err != nil {
	// 	return err
	// }
	// output := out.String()
	return nil
}
