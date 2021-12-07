package auth

import (
	"github.com/fingcloud/cli/pkg/cli"
	"github.com/fingcloud/cli/pkg/config/session"
	"github.com/fingcloud/cli/pkg/util"
	"github.com/spf13/cobra"
)

type ListOptions struct {
	format string
}

func NewCmdList(ctx *cli.Context) *cobra.Command {
	opts := new(ListOptions)

	cmd := &cobra.Command{
		Use:     "list",
		Short:   "list your accounts",
		Aliases: []string{"ls"},
		Run: func(cmd *cobra.Command, args []string) {
			ctx.SetupClient()

			util.CheckErr(runList(ctx, opts))
		},
	}

	cmd.Flags().StringVar(&opts.format, "format", "", "format of result in go template")

	return cmd
}

func runList(ctx *cli.Context, opts *ListOptions) error {
	sessions, err := session.Read()
	util.CheckErr(err)

	if opts.format == "" {
		opts.format = defaultFormat
	}

	err = PrintFormat(ctx.Stdout, opts.format, sessions)
	util.CheckErr(err)

	return nil
}
