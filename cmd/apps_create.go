package cmd

import (
	"github.com/fingcloud/cli/internal/cli"
	"github.com/spf13/cobra"
)

func NewAppsCreateCommand() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "create",
		Short: "create an app",
		Run: func(cmd *cobra.Command, args []string) {
			cli := cli.New(cmd, args, token, devMode)
			runCreateApp(cli)
		},
	}

	return cmd
}

func runCreateApp(cli *cli.FingCli) {

}
