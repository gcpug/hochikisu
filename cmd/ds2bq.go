package cmd

import (
	"context"
	"fmt"
	"os"

	"github.com/gcpug/hochikisu/ds2bq"
	"github.com/gcpug/hochikisu/scheduler"
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

	jobs, err := ds2bq.ParseYaml(ctx, body)
	if err != nil {
		fmt.Println("failed to parse yaml file.")
		fmt.Println(err)
		os.Exit(1)
	}

	for _, job := range jobs {
		c, err := scheduler.NewClient(ctx, job.ProjectID)
		if err != nil {
			fmt.Println("failed scheduler new client.")
			fmt.Println(err)
			os.Exit(1)
		}

		req, err := job.CreateUpsertJobRequest()
		if err != nil {
			fmt.Println("failed create upsert job request.")
			fmt.Println(err)
			os.Exit(1)
		}
		rj, err := c.Upsert(ctx, req)
		if err != nil {
			fmt.Println("failed schedule upsert")
			fmt.Println(err)
			os.Exit(1)
		}
		fmt.Printf("upsert completed: %s\n", rj.Name)
	}

	fmt.Println("done!")
}
