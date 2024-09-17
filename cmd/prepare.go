package cmd

import (
	"fmt"

	"github.com/benleem/prattl/pysrc"
	"github.com/spf13/cobra"
)

func init() {
	rootCmd.AddCommand(prepareCommand)
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

		// s := spinner.New(spinner.CharSets[35], 100*time.Millisecond, spinner.WithWriter(os.Stderr))
		// s.Suffix = "\n"
		// s.Start()
		err = pysrc.PrepareDistribution(*env)
		if err != nil {
			if err.Error() == "dist exists" {
				fmt.Println("prattl distribution already prepared")
				return nil
			}
			return fmt.Errorf("Error preparing distribution env: %v\n", err)
		}
		// s.Stop()
		fmt.Println("successfully prepared prattl distribution")
		return nil
	},
}
