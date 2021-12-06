package app

import (
	"github.com/fingcloud/cli/pkg/cli"
	"github.com/spf13/cobra"
)

func NewAppsCmd(ctx *cli.Context) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "app [command]",
		Short: "manage apps",
		Args:  cli.NoArgs,
	}

	cmd.AddCommand(
		NewCmdList(ctx),
		NewCmdCreate(ctx),
		NewCmdRemove(ctx),
		NewCmdStart(ctx),
		NewCmdRestart(ctx),
		NewCmdStop(ctx),
		NewCmdLogs(ctx),
	)

	return cmd
}
