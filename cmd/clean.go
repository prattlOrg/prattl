package cmd

import (
	"bufio"
	"fmt"
	"os"
	"unicode"

	"github.com/prattlOrg/prattl/internal/pysrc"
	"github.com/spf13/cobra"
)

var Confirm bool

func init() {
	cleanCmd.Flags().BoolVarP(&Confirm, "confirm", "y", false, "skips confirmation prompt")
	RootCmd.AddCommand(cleanCmd)
}

var cleanCmd = &cobra.Command{
	Use:   "clean",
	Short: "Remove the python distribution built by prattl",
	Long:  "This command removes everything prattl adds to your filesystem",
	Args:  cobra.ExactArgs(0),
	RunE: func(cmd *cobra.Command, args []string) error {
		env, err := pysrc.GetPrattlEnv()
		if err != nil {
			return err
		}

		if Confirm {
			fmt.Printf("Removing %s\n", env.EnvOptions.ParentPath)
		} else {
			fmt.Printf("Are you sure you want to delete %s? [Y/N]\n", env.EnvOptions.ParentPath)
		}
		reader := bufio.NewReader(os.Stdin)

		var proceed bool
		for {
			if Confirm {
				proceed = true
				break
			}

			char, _, err := reader.ReadRune()
			if err != nil {
				return fmt.Errorf("error reading from stdin: %v", err)
			}
			switch unicode.ToLower(char) {
			case 'y':
				proceed = true
			case 'n':
				proceed = false
			case '\n':
				continue
			default:
				fmt.Printf("Unexpected input: %c\nExpected [Y/N]\n", char)
				continue

			}
			break
		}

		if proceed {
			err = os.RemoveAll(env.EnvOptions.ParentPath)
			if err != nil {
				return fmt.Errorf("Problem cleaning %s: %v", env.EnvOptions.ParentPath, err)
			}
			fmt.Println("Successfully cleaned prattl directory")
			return nil
		} else {
			fmt.Println("Clean Cancelled")
			return nil

		}
	},
}
