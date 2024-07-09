package pysrc

import (
	"embed"
)

//go:embed py
var PythonSrc embed.FS

func ReturnFile(fp string) (string, error) {
	data, err := PythonSrc.ReadFile("py/" + fp)
	if err != nil {
		return "", err
	}
	return (string(data)), nil
}
