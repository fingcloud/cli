package app

import (
	"os"

	"github.com/fingcloud/cli/pkg/api"
	"github.com/fingcloud/cli/pkg/cli"
	"github.com/fingcloud/cli/pkg/util"
	"github.com/spf13/cobra"
)

type ListOptions struct {
	kind   string
	format string
}

func NewCmdList(ctx *cli.Context) *cobra.Command {
	opts := new(ListOptions)

	cmd := &cobra.Command{
		Use:     "list",
		Short:   "list apps",
		Aliases: []string{"ls"},
		Args:    cli.NoArgs,
		Run: func(cmd *cobra.Command, args []string) {

			util.CheckErr(RunList(ctx, opts))
		},
	}

	cmd.Flags().StringVar(&opts.format, "format", "", "format of result in go template")

	return cmd
}

func RunList(ctx *cli.Context, opts *ListOptions) error {
	apps, err := ctx.Client.AppsList(&api.ListAppsOptions{})
	util.CheckErr(err)

	if opts.format == "" {
		opts.format = defaultFormat
	}

	return PrintFormat(os.Stdout, opts.format, apps)
}
