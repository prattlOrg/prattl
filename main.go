package main

import (
	"fmt"
	"os"
	"prattl/cmd"
	ffmpeg "prattl/internal"

	"github.com/voidKandy/go-pyenv/pyenv"
)

func main() {
	err := ffmpeg.CheckInstall()
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	env := pyenv.DefaultPyEnv()
	exists, _ := env.DistExists()
	if !*exists {
		env.MacInstall()
		env.AddDependencies("pysrc/requirements.txt")
	}

	cmd.Execute()
}
