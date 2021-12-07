package app

import (
	"fmt"

	"github.com/fingcloud/cli/pkg/api"
	"github.com/fingcloud/cli/pkg/cli"
	"github.com/fingcloud/cli/pkg/ui"
	"github.com/fingcloud/cli/pkg/util"
	"github.com/spf13/cobra"
)

type RemoveOptions struct {
	Name  string
	Force bool
}

func NewCmdRemove(ctx *cli.Context) *cobra.Command {
	opts := new(RemoveOptions)

	cmd := &cobra.Command{
		Use:     "remove [app]",
		Short:   "remove an app",
		Aliases: []string{"rm"},
		Args:    cli.Exact(1),
		Run: func(cmd *cobra.Command, args []string) {
			opts.Name = args[0]

			util.CheckErr(runRemove(ctx, opts))
		},
	}

	cmd.Flags().BoolVarP(&opts.Force, "force", "f", false, "force the removal of app")

	return cmd
}

func runRemove(ctx *cli.Context, opts *RemoveOptions) error {
	if !opts.Force {
		var confirmation string
		ui.PromptInput(fmt.Sprintf("Enter app name to remove (%s):", opts.Name), &confirmation)
		if opts.Name != confirmation {
			return fmt.Errorf("invalid app name %s", confirmation)
		}
	}

	err := ctx.Client.AppsRemove(&api.RemoveAppOptions{
		Name: opts.Name,
	})
	util.CheckErr(err)

	fmt.Printf("app '%s' removed\n", opts.Name)
	return nil
}
