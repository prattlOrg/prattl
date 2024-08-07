package main

import (
	"fmt"
	"os"
	"prattl/cmd"
	ffmpeg "prattl/internal"
)

func main() {
	err := ffmpeg.CheckInstall()
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
	if err != nil {
		fmt.Printf("error creating prattl env: %v\n", err)
		os.Exit(1)
	}
	cmd.Execute()
}
