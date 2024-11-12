package cmd

import (
	"fmt"
	"os"
	"path/filepath"
	"sync"

	"github.com/prattlOrg/go-pyenv"
	"github.com/prattlOrg/prattl/internal/pysrc"
	"github.com/spf13/cobra"
)

func init() {
	RootCmd.AddCommand(reportCommand)
}

// ripped straight from stack overflow: https://stackoverflow.com/questions/32482673/how-to-get-directory-total-size
func DirSize(path string) (int64, error) {
	var size int64
	var mu sync.Mutex

	// Function to calculate size for a given path
	var calculateSize func(string) error
	calculateSize = func(p string) error {
		fileInfo, err := os.Lstat(p)
		if err != nil {
			return err
		}

		// Skip symbolic links to avoid counting them multiple times
		if fileInfo.Mode()&os.ModeSymlink != 0 {
			return nil
		}

		if fileInfo.IsDir() {
			entries, err := os.ReadDir(p)
			if err != nil {
				return err
			}
			for _, entry := range entries {
				if err := calculateSize(filepath.Join(p, entry.Name())); err != nil {
					return err
				}
			}
		} else {
			mu.Lock()
			size += fileInfo.Size()
			mu.Unlock()
		}
		return nil
	}

	// Start calculation from the root path
	if err := calculateSize(path); err != nil {
		return 0, err
	}

	return size, nil
}

var reportCommand = &cobra.Command{
	Use:   "report",
	Short: "Gives a report of prattl's python distribution",
	Long:  "Gives information on whether prattl's python distribution is compressed, and how much space it is taking up",
	Args:  cobra.ExactArgs(0),
	RunE: func(cmd *cobra.Command, args []string) error {
		env, err := pysrc.GetPrattlEnv()
		if err != nil {
			return fmt.Errorf("Error getting prattl env: %v\n", err)
		}

		exists, err := env.EnvOptions.DistExists()
		if err != nil {
			return fmt.Errorf("Error checking prattl's python distribution exists: %v\n", err)
		}
		if !*exists {
			fmt.Println("Distribution doesn't exist, you need to run prattl prepare before running this command")
			return nil
		}

		var distPath string
		if env.EnvOptions.Compressed {
			distPath = pyenv.DistZipPath(&env.EnvOptions)
		} else {
			distPath = pyenv.DistDirPath(&env.EnvOptions)
		}

		info, err := os.Stat(distPath)
		if err != nil {
			return fmt.Errorf("Error getting info of %s: %v\n", distPath, err)
		}

		var size int64
		if info.IsDir() {
			size, err = DirSize(distPath)
			if err != nil {
				return fmt.Errorf("Error getting size of dist directory: %v\n", err)
			}
		} else {
			size = info.Size()
		}

		fmt.Println("------prattl distribution report------")
		fmt.Printf("Location: %v\n", distPath)
		fmt.Printf("Compressed: %v\n", env.EnvOptions.Compressed)
		fmt.Printf("Size: %d bytes\n", size)
		fmt.Printf("Last Modified: %v\n", info.ModTime())
		fmt.Println("--------------------------------------")
		return nil
	},
}
