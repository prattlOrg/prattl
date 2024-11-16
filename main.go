package main

import (
	"log"

	"github.com/prattlOrg/prattl/cmd"
	"github.com/prattlOrg/prattl/internal/ffmpeg"
)

var version = "dev"

func main() {
	err := ffmpeg.CheckInstall()
	if err != nil {
		log.Fatalf("error creating prattl env: %v\n", err)
	}
	cmd.Execute(version)
}
