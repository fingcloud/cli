package app

import (
	"fmt"

	"github.com/fingcloud/cli/pkg/api"
	"github.com/fingcloud/cli/pkg/cli"
	"github.com/fingcloud/cli/pkg/util"
	"github.com/spf13/cobra"
)

type StartOptions struct {
	Name string
}

func NewCmdStart(ctx *cli.Context) *cobra.Command {
	opts := new(StartOptions)

	cmd := &cobra.Command{
		Use:   "start [app]",
		Short: "start app",
		Args:  cli.Exact(1),
		Run: func(cmd *cobra.Command, args []string) {
			ctx.SetupClient()
			opts.Name = args[0]

			util.CheckErr(RunStart(ctx, opts))
		},
	}

	return cmd
}

func RunStart(ctx *cli.Context, opts *StartOptions) error {
	err := ctx.Client.AppsStart(&api.StartAppOptions{
		Name: opts.Name,
	})
	util.CheckErr(err)
	fmt.Printf("app '%s' started\n", opts.Name)

	return nil
}
