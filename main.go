package main

import (
	"log"

	"github.com/prattlOrg/prattl/cmd"
	ffmpeg "github.com/prattlOrg/prattl/internal"
)

func main() {
	err := ffmpeg.CheckInstall()
	if err != nil {
		log.Fatalf("error creating prattl env: %v\n", err)
	}
	cmd.Execute()
}
