package app

import (
	"github.com/fingcloud/cli/pkg/cli"
	"github.com/spf13/cobra"
)

func NewAppsCmd(ctx *cli.Context) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "app COMMAND",
		Short: "manage apps",
		Args:  cli.NoArgs,
	}

	cmd.AddCommand(
		newListCmd(ctx),
		newCreateCmd(ctx),
		newRestartCmd(ctx),
	)

	return cmd
}
