package pysrc

import (
	"embed"
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"runtime"
	"strings"

	pyenv "github.com/prattlOrg/go-pyenv"
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

type PrattlEnv struct {
	pyenv      pyenv.PyEnv
	compressed bool
}

func GetPrattlEnv() (*pyenv.PyEnv, error) {
	home, err := os.UserHomeDir()
	if err != nil {
		return nil, err
	}
	parentPath := filepath.Join(home, ".prattl")
	osArch := fmt.Sprintf("%s/%s", runtime.GOOS, runtime.GOARCH)
	env, err := pyenv.NewPyEnv(parentPath, osArch)
	if err != nil {
		return nil, err
	}
	_, err = os.Stat(pyenv.DistZipPath(&env.EnvOptions))
	if err != nil {
		if errors.Is(err, os.ErrNotExist) {
			return env, nil
		}
		return nil, err
	}
	env.EnvOptions.Compressed = true
	return env, nil
}

func PrepareDistribution(env pyenv.PyEnv) error {
	exists, _ := env.EnvOptions.DistExists()
	if !*exists {
		// s.Prefix = "installing python distribution: "
		// install needs to return error
		env.Install()
		// s.Prefix = "downloading dependencies:"
		err := downloadDeps(env)
		if err != nil {
			return fmt.Errorf("error downloading prattl dependencies: %v\n", err)
		}
	} else {
		return fmt.Errorf("dist exists")
	}
	return nil
}

func downloadDeps(env pyenv.PyEnv) error {
	var requirementsFp string
	switch {
	case strings.Contains(env.EnvOptions.Distribution, "darwin"):
		requirementsFp = "requirements-darwin.txt"
	case strings.Contains(env.EnvOptions.Distribution, "linux"):
		requirementsFp = "requirements-linux.txt"
	case strings.Contains(env.EnvOptions.Distribution, "windows"):
		requirementsFp = "requirements-windows.txt"
	}

	reqs, err := ReturnFile(requirementsFp)
	if err != nil {
		return err
	}
	path := filepath.Join(env.EnvOptions.ParentPath, requirementsFp)
	err = os.WriteFile(path, []byte(reqs), 0o640)
	if err != nil {
		return err
	}
	err = env.AddDependencies(path)
	if err != nil {
		return err
	}
	return nil
}
