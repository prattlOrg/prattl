package cmd

import (
	"bufio"
	"fmt"
	"os"
	"unicode"

	"github.com/benleem/prattl/pysrc"
	"github.com/spf13/cobra"
)

func init() {
	rootCmd.AddCommand(cleanCommand)
}

var cleanCommand = &cobra.Command{
	Use:   "clean",
	Short: "Remove the python distribution built by prattl",
	Long:  "This command removes everything prattl adds to your filesystem",
	Args:  cobra.ExactArgs(0),
	RunE: func(cmd *cobra.Command, args []string) error {
		env, err := pysrc.PrattlEnv()
		if err != nil {
			return fmt.Errorf("Error getting prattl env: %v\n", err)
		}

		fmt.Printf("Are you sure you want to delete %s? [Y/N]\n", env.ParentPath)
		reader := bufio.NewReader(os.Stdin)

		var proceed bool
		for {

			char, _, err := reader.ReadRune()
			if err != nil {
				return fmt.Errorf("Error reading from stdin: %v\n", err)
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
			err = os.RemoveAll(env.ParentPath)
			if err != nil {
				return fmt.Errorf("Problem cleaning %s: %v", env.ParentPath, err)
			}
			fmt.Println("Successfully cleaned prattl directory")
			return nil
		} else {
			fmt.Println("Clean Cancelled")
			return nil

		}
	},
}
