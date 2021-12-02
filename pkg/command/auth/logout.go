package auth

import (
	"fmt"

	"github.com/fingcloud/cli/pkg/cli"
	"github.com/fingcloud/cli/pkg/command/util"
	"github.com/fingcloud/cli/pkg/config/session"
	"github.com/fingcloud/cli/pkg/ui"
	"github.com/spf13/cobra"
)

type LogoutOptions struct{}

func NewCmdLogout(ctx *cli.Context) *cobra.Command {
	o := new(LogoutOptions)

	cmd := &cobra.Command{
		Use:   "logout",
		Short: "logout from fing service",
		Long:  "logout from fing service",
		Run: func(cmd *cobra.Command, args []string) {
			ctx.SetupClient()

			util.CheckErr(o.Init(ctx, args))
			util.CheckErr(o.Validate())
			util.CheckErr(o.Run(ctx))
		},
	}

	return cmd
}

func (o *LogoutOptions) Init(ctx *cli.Context, args []string) error {
	return nil
}

func (o *LogoutOptions) Validate() error {
	return nil
}

func (o *LogoutOptions) Run(ctx *cli.Context) error {
	sess, err := session.RemoveSession()
	if err != nil {
		return err
	}

	fmt.Println(ui.Gray(fmt.Sprintf("Logged out %s", sess.Email)))
	return nil
}
