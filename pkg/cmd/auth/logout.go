package auth

import (
	"fmt"

	"github.com/fingcloud/cli/pkg/cli"
	"github.com/fingcloud/cli/pkg/config/session"
	"github.com/fingcloud/cli/pkg/util"
	"github.com/spf13/cobra"
)

type LogoutOptions struct{}

func NewCmdLogout(ctx *cli.Context) *cobra.Command {
	opts := new(LogoutOptions)

	cmd := &cobra.Command{
		Use:     "logout",
		Short:   "logout your account",
		Aliases: []string{"rm"},
		Run: func(cmd *cobra.Command, args []string) {
			ctx.SetupClient()

			runLogout(ctx, opts)
		},
	}

	return cmd
}

func runLogout(ctx *cli.Context, opts *LogoutOptions) error {
	sess, err := session.RemoveSession()
	util.CheckErr(err)

	fmt.Printf("Logged out %s\n", sess.Email)
	return nil
}
