package auth

import (
	"github.com/fingcloud/cli/pkg/cli"
	"github.com/spf13/cobra"
)

func NewCmd(ctx *cli.Context) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "auth [command]",
		Short: "login, logout, switch your authentications",
		Args:  cobra.NoArgs,
	}

	cmd.AddCommand(
		NewCmdLogin(ctx),
		NewCmdLogout(ctx),
		NewCmdUse(ctx),
		NewCmdList(ctx),
	)

	return cmd
}
