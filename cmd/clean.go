package cmd

import (
	"fmt"
	"os"

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

		err = os.RemoveAll(env.ParentPath)
		if err != nil {
			return fmt.Errorf("Problem cleaning %s: %v", env.ParentPath, err)
		}
		fmt.Println("Successfully cleaned prattl directory")
		return nil
	},
}
