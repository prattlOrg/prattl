package pysrc

import (
	"embed"
	"fmt"
	"os"
	"path/filepath"
	"runtime"
	"strings"

	"github.com/prattlOrg/go-pyenv/pyenv"
)

//go:embed py
var PythonSrc embed.FS

func ReturnFile(fp string) (string, error) {
	data, err := PythonSrc.ReadFile("py/" + fp)
	if err != nil {
		return "", fmt.Errorf("error returning file: %v", fp)
	}
	return (string(data)), nil
}

func PrattlEnv() (*pyenv.PyEnv, error) {
	home, err := os.UserHomeDir()
	if err != nil {
		return nil, fmt.Errorf("error getting $HOME directory: %v", err)
	}
	parentPath := filepath.Join(home, ".prattl")
	osArch := fmt.Sprintf("%s/%s", runtime.GOOS, runtime.GOARCH)
	env, err := pyenv.NewPyEnv(parentPath)
	if err != nil {
		return nil, err
	}
	env.Distribution = osArch
	return env, nil
}

func PrepareDistribution(env pyenv.PyEnv) error {
	exists, _ := env.DistExists()
	if !*exists {
		// s.Prefix = "installing python distribution: "
		// install needs to return error
		err := env.Install()
		if err != nil {
			return err
		}
		// s.Prefix = "downloading dependencies:"
		err = downloadDeps(env)
		if err != nil {
			return err
		}
	} else {
		return fmt.Errorf("dist exists")
	}
	return nil
}

func downloadDeps(env pyenv.PyEnv) error {
	var requirementsFp string
	switch {
	case strings.Contains(env.Distribution, "darwin"):
		requirementsFp = "requirements-darwin.txt"
	case strings.Contains(env.Distribution, "linux"):
		requirementsFp = "requirements-linux.txt"
	case strings.Contains(env.Distribution, "windows"):
		requirementsFp = "requirements-windows.txt"
	}

	reqs, err := ReturnFile(requirementsFp)
	if err != nil {
		return err
	}
	path := filepath.Join(env.ParentPath, requirementsFp)
	err = os.WriteFile(path, []byte(reqs), 0o640)
	if err != nil {
		return fmt.Errorf("error copying python requirements to %v: %v", path, err)
	}
	err = env.AddDependencies(path)
	if err != nil {
		return err
	}
	return nil
}
