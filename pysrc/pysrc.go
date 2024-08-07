package pysrc

import (
	"embed"
	"fmt"
	"os"
	"path/filepath"

	"github.com/voidKandy/go-pyenv/pyenv"
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

func PrattlEnv() (*pyenv.PyEnv, error) {
	home, err := os.UserHomeDir()
	if err != nil {
		return nil, err
	}
	parentPath := filepath.Join(home, ".prattl/")
	env := pyenv.PyEnv{
		ParentPath: parentPath,
	}
	return &env, nil
}

func PrepareDistribution(env pyenv.PyEnv) error {
	exists, _ := env.DistExists()
	if !*exists {
		// mac install needs to return error
		env.MacInstall()
		err := downloadDeps(env)
		if err != nil {
			return fmt.Errorf("Error downloading prattl dependencies: %v\n", err)
		}
	} else {
		return fmt.Errorf("dist exists")
	}
	return nil
}

func downloadDeps(env pyenv.PyEnv) error {
	reqs, err := ReturnFile("requirements.txt")
	if err != nil {
		return err
	}
	path := filepath.Join(env.ParentPath, "requirements.txt")
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
