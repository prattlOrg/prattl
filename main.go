package main

import (
	"fmt"
	"os"
	"path/filepath"
	"prattl/cmd"
	ffmpeg "prattl/internal"
	"prattl/pysrc"

	"github.com/voidKandy/go-pyenv/pyenv"
)

func main() {
	err := ffmpeg.CheckInstall()
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	home, err := os.UserHomeDir()
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
	env := pyenv.PyEnv{
		ParentPath: home + "/.prattl/",
	}

	exists, _ := env.DistExists()

	if !*exists {
		env.MacInstall()
		downloadDeps(env)
	}

	cmd.Execute()
}

func downloadDeps(env pyenv.PyEnv) error {
	reqs, err := pysrc.ReturnFile("requirements.txt")
	if err != nil {
		return err
	}

	path := filepath.Join(env.ParentPath + "requirements.txt")

	err = os.WriteFile(path, []byte(reqs), 0o640)

	env.AddDependencies(path)
	return nil
}
