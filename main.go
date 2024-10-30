package main

import (
	"fmt"
	"os"

	"github.com/prattlOrg/prattl/cmd"
	ffmpeg "github.com/prattlOrg/prattl/internal"
)

func main() {
	err := ffmpeg.CheckInstall()
	if err != nil {
		fmt.Printf("error creating prattl env: %v\n", err)
		os.Exit(1)
	}
	cmd.Execute()
}
