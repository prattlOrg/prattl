package cmd

import (
	"fmt"

	"github.com/prattlOrg/prattl/pysrc"
	"github.com/spf13/cobra"
)

func init() {
	rootCmd.AddCommand(prepareCommand)
}

var prepareCommand = &cobra.Command{
	Use:   "prepare",
	Short: "Prepare the python distribution required by prattl",
	Long:  "This command prepares the distribution needed to use prattl",
	Args:  cobra.ExactArgs(0),
	RunE: func(cmd *cobra.Command, args []string) error {
		env, err := pysrc.PrattlEnv()
		if err != nil {
			return err
		}
		// s := spinner.New(spinner.CharSets[35], 100*time.Millisecond, spinner.WithWriter(os.Stderr))
		// s.Suffix = "\n"
		// s.Start()
		err = pysrc.PrepareDistribution(*env)
		if err != nil {
			if err.Error() == "dist exists" {
				return fmt.Errorf("prattl distribution already prepared")
			}
			return err
		}
		// s.Stop()
		fmt.Println("Successfully prepared prattl distribution")
		return nil
	},
}
