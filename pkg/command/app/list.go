package app

import (
	"os"

	"github.com/fingcloud/cli/pkg/api"
	"github.com/fingcloud/cli/pkg/cli"
	"github.com/fingcloud/cli/pkg/command/util"
	"github.com/spf13/cobra"
	"github.com/spf13/pflag"
)

type listOptions struct {
	kind   string
	format string
}

func newListCmd(ctx *cli.Context) *cobra.Command {
	opts := new(listOptions)

	cmd := &cobra.Command{
		Use:     "ls",
		Short:   "list apps",
		Aliases: []string{"list"},
		Args:    cli.NoArgs,
		Run: func(cmd *cobra.Command, args []string) {
			ctx.SetupClient()

			if err := runList(ctx, cmd.Flags(), opts); err != nil {
				util.CheckErr(err)
			}
		},
	}

	cmd.Flags().StringVar(&opts.kind, "type", "", "one of [web, service]")
	cmd.Flags().StringVar(&opts.format, "format", "", "format of result in go template")

	return cmd
}

func runList(ctx *cli.Context, flags *pflag.FlagSet, opts *listOptions) error {
	apps, err := ctx.Client.AppsList(&api.ListAppsOptions{})
	if err != nil {
		return err
	}

	if opts.format == "" {
		opts.format = defaultFormat
	}

	return PrintFormat(os.Stdout, opts.format, apps)
}
