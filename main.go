package main

import (
	"prattl/cmd"
	"prattl/internal/pyenv"
)

func main() {
	pyenv.Main()
	cmd.Execute()
}
