package pysrc

import (
	"embed"
	"fmt"
	"os"
	"path/filepath"
	"time"

	"github.com/briandowns/spinner"
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
	s := spinner.New(spinner.CharSets[35], 100*time.Millisecond, spinner.WithWriter(os.Stderr))
	if !*exists {
		s.Suffix = "\n"
		s.Start()
		s.Prefix = "installing python distribution: "
		// mac install needs to return error
		env.MacInstall()
		s.Prefix = "downloading dependencies: "
		s.Restart()
		err := downloadDeps(env)
		if err != nil {
			s.Stop()
			return fmt.Errorf("Error downloading prattl dependencies: %v\n", err)
		}
	} else {
		s.Stop()
		return fmt.Errorf("dist exists")
	}
	s.Stop()
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
