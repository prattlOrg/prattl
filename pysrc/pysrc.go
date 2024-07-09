package pysrc

import (
	"embed"
)

//go:embed transcribe.py
var pythonSrc embed.FS

func ReturnSrc() (string, error) {
	data, err := pythonSrc.ReadFile("transcribe.py")
	if err != nil {
		return "", err
	}
	return (string(data)), nil
}
