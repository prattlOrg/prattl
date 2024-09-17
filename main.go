package main

import (
	"fmt"
	"os"

	"github.com/benleem/prattl/cmd"
	ffmpeg "github.com/benleem/prattl/internal"
)

func main() {
	err := ffmpeg.CheckInstall()
	if err != nil {
		fmt.Printf("error creating prattl env: %v\n", err)
		os.Exit(1)
	}
	cmd.Execute()
}
