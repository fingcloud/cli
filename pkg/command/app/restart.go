package app

import (
	"fmt"

	"github.com/fingcloud/cli/pkg/api"
	"github.com/fingcloud/cli/pkg/cli"
	"github.com/fingcloud/cli/pkg/command/util"
	"github.com/fingcloud/cli/pkg/ui"
	"github.com/spf13/cobra"
	"github.com/spf13/pflag"
)

type restartOptions struct {
	Name string
}

func newRestartCmd(ctx *cli.Context) *cobra.Command {
	opts := new(restartOptions)

	cmd := &cobra.Command{
		Use:   "restart APP",
		Short: "restart an app",
		Args:  cli.Exact(1),
		Run: func(cmd *cobra.Command, args []string) {
			ctx.SetupClient()
			opts.Name = args[0]

			if err := runRestart(ctx, cmd.Flags(), opts); err != nil {
				util.CheckErr(err)
			}
		},
	}

	return cmd
}

func runRestart(ctx *cli.Context, flags *pflag.FlagSet, opts *restartOptions) error {
	err := ctx.Client.AppsRestart(&api.RestartAppOptions{
		Name: opts.Name,
	})
	if err != nil {
		return err
	}

	fmt.Println(ui.Gray(fmt.Sprintf("%s restarted successfully", opts.Name)))

	return nil
}
