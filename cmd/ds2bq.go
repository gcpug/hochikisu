package cmd

import (
	"context"
	"fmt"
	"os"

	"github.com/spf13/cobra"
)

var cmdDS2BQ = &cobra.Command{
	Use:   "ds2bq",
	Short: "Commands for gcpug/ds2bq",
	Long:  `ds2bq is a function to manage gcpug/ds2bq settings in yaml file`,
}

var cmdDS2BQDeploy = &cobra.Command{
	Use:   "deploy",
	Short: "Deploy the yaml file.",
	Long:  `Create or update settings according to yaml. Do not delete.`,
	Run: func(cmd *cobra.Command, args []string) {
		ctx := context.Background()
		ds2bqDeploy(ctx, cfgFile)
	},
}

func ds2bqDeploy(ctx context.Context, cfgFile string) {
	body, err := ReadFile(cfgFile)
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
	fmt.Println(string(body))
}
