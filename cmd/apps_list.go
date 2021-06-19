package cmd

import (
	"fmt"
	"log"

	"github.com/fingcloud/cli/api"
	"github.com/fingcloud/cli/internal/cli"
	"github.com/spf13/cobra"
)

func NewAppsListCommand() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "List",
		Short: "list all apps",
		Run: func(cmd *cobra.Command, args []string) {
			cli := cli.New(cmd, args, token, devMode)
			runListApp(cli)
		},
	}

	return cmd
}

func runListApp(cli *cli.FingCli) {
	apps, err := cli.Client.AppsList(&api.ListAppsOptions{})
	if err != nil {
		fmt.Println(err)
	}

	for _, app := range apps {
		log.Println(app)
	}
}
