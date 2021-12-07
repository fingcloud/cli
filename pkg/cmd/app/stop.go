package app

import (
	"fmt"

	"github.com/fingcloud/cli/pkg/api"
	"github.com/fingcloud/cli/pkg/cli"
	"github.com/fingcloud/cli/pkg/util"
	"github.com/spf13/cobra"
)

type StopOptions struct {
	Name string
}

func NewCmdStop(ctx *cli.Context) *cobra.Command {
	opts := new(StopOptions)

	cmd := &cobra.Command{
		Use:     "stop [app]",
		Short:   "stop app",
		Aliases: []string{"shutdown"},
		Args:    cli.Exact(1),
		Run: func(cmd *cobra.Command, args []string) {
			ctx.SetupClient()
			opts.Name = args[0]

			util.CheckErr(RunStop(ctx, opts))
		},
	}

	return cmd
}

func RunStop(ctx *cli.Context, opts *StopOptions) error {
	err := ctx.Client.AppsShutdown(&api.ShutdownAppOptions{
		Name: opts.Name,
	})
	util.CheckErr(err)
	fmt.Printf("app '%s' stopped\n", opts.Name)

	return nil
}
