package app

import (
	"fmt"

	"github.com/fingcloud/cli/pkg/api"
	"github.com/fingcloud/cli/pkg/cli"
	"github.com/fingcloud/cli/pkg/util"
	"github.com/spf13/cobra"
	"github.com/spf13/pflag"
)

type RestartOptions struct {
	Name string
}

func NewCmdRestart(ctx *cli.Context) *cobra.Command {
	opts := new(RestartOptions)

	cmd := &cobra.Command{
		Use:   "restart [app]",
		Short: "restart app",
		Args:  cli.Exact(1),
		Run: func(cmd *cobra.Command, args []string) {
			ctx.SetupClient()
			opts.Name = args[0]

			if err := RunRestart(ctx, cmd.Flags(), opts); err != nil {
				util.CheckErr(err)
			}
		},
	}

	return cmd
}

func RunRestart(ctx *cli.Context, flags *pflag.FlagSet, opts *RestartOptions) error {
	err := ctx.Client.AppsRestart(&api.RestartAppOptions{
		Name: opts.Name,
	})
	if err != nil {
		return err
	}

	fmt.Printf("app '%s' restarted\n", opts.Name)

	return nil
}