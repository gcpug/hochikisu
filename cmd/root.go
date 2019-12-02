package cmd

import (
	"github.com/spf13/cobra"
)

var (
	// Used for flags.
	cfgFile string

	rootCmd = &cobra.Command{
		Use:   "hochikisu",
		Short: "Manage tool settings under gcpug",
		Long:  `Manage tools that need to be configured with yaml file with tools under gcpug organization`,
	}
)

// Execute executes the root command.
func Execute() error {
	if err := rootCmd.MarkPersistentFlagRequired("config"); err != nil {
		return err
	}

	cmdDS2BQ.AddCommand(cmdDS2BQDeploy)
	rootCmd.AddCommand(cmdDS2BQ)
	return rootCmd.Execute()
}

func init() {
	rootCmd.PersistentFlags().StringVar(&cfgFile, "config", "", "config file path")
}
