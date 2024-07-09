package main

import (
	"prattl/cmd"

	pyenv "github.com/voidKandy/go-pyenv"
)

func main() {
	pyenv.Main()
	cmd.Execute()
}
