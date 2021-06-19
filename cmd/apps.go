package cmd

import (
	"github.com/fingcloud/cli/internal/cli"
	"github.com/spf13/cobra"
)

func NewAppsCommand() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "apps",
		Short: "manage your apps",
		Run: func(cmd *cobra.Command, args []string) {
			cli := cli.New(cmd, args, token, devMode)
			runListApp(cli)
		},
	}

	cmd.AddCommand(
		NewAppsListCommand(),
		NewAppsCreateCommand(),
	)

	return cmd
}
