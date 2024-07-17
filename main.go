package main

import (
	"fmt"
	"os"
	"prattl/cmd"
	ffmpeg "prattl/internal"
	"prattl/pysrc"
)

func main() {
	err := ffmpeg.CheckInstall()
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	env, err := pysrc.PrattlEnv()
	if err != nil {
		fmt.Printf("error creating prattl env: %v\n", err)
		os.Exit(1)
	}
	// pysrc.PrepareDistribution(*env)
	cmd.Execute(*env)
}
