package cmd

import (
	"fmt"
	"os"
	"prattl/pysrc"

	"github.com/spf13/cobra"
)

func init() {
	rootCmd.AddCommand(cleanCommand)
	rootCmd.AddCommand(prepareCommand)
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

		err = os.RemoveAll(env.ParentPath)
		if err != nil {
			return fmt.Errorf("Problem cleaning %s: %v", env.ParentPath, err)
		}
		fmt.Println("Successfully cleaned prattl directory")
		return nil
	},
}

var prepareCommand = &cobra.Command{
	Use:   "prepare",
	Short: "prepare the python distribution required by prattl",
	Long:  "This command prepares the distribution needed to use prattl",
	Args:  cobra.ExactArgs(0),
	RunE: func(cmd *cobra.Command, args []string) error {
		env, err := pysrc.PrattlEnv()
		if err != nil {
			return fmt.Errorf("Error getting prattl env: %v\n", err)
		}
		err = pysrc.PrepareDistribution(*env)
		if err != nil {
			if err.Error() == "dist exists" {
				fmt.Println("prattl distribution already prepared")
				return nil
			}
			return fmt.Errorf("Error preparing distribution env: %v\n", err)
		}
		fmt.Println("Successfully prepared prattl distribution")
		return nil
	},
}
