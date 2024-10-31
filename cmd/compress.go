package cmd

import (
	"fmt"

	"github.com/prattlOrg/prattl/internal/pysrc"
	"github.com/spf13/cobra"
)

func init() {
	rootCmd.AddCommand(compressCommand)
	rootCmd.AddCommand(decompressCommand)
}

var compressCommand = &cobra.Command{
	Use:   "compress",
	Short: "Compress python distribution to save space",
	Long:  "This command will compress prattl's internal python distribution in order to save space on your computer",
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

		if env.EnvOptions.Compressed {
			fmt.Println("Distribution is already compressed")
			return nil
		}

		err = env.EnvOptions.CompressDist()
		if err != nil {
			fmt.Printf("An error occurred when compressing distribution: %v\n", err)
			return err
		}
		fmt.Println("Successfully compressed distribution")
		return nil
	},
}

var decompressCommand = &cobra.Command{
	Use:   "decompress",
	Short: "Decompress python distribution",
	Long:  "This command will decompress prattl's internal python distribution in order to use",
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

		if !env.EnvOptions.Compressed {
			fmt.Println("Distribution is already decompressed")
			return nil
		}

		err = env.EnvOptions.DecompressDist()
		if err != nil {
			fmt.Printf("An error occurred when decompressing distribution: %v\n", err)
			return err
		}
		fmt.Println("Successfully decompressed distribution")
		return nil
	},
}
