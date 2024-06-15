package pysrc

import (
	"embed"
)

//go:embed transcribe.py
var f embed.FS

func ReturnSrc() (string, error) {
	data, err := f.ReadFile("transcribe.py")
	if err != nil {
		return "", err
	}
	return (string(data)), nil
}
