package auth

import (
	"errors"
	"fmt"

	"github.com/fingcloud/cli/pkg/cli"
	"github.com/fingcloud/cli/pkg/config/session"
	"github.com/fingcloud/cli/pkg/ui"
	"github.com/fingcloud/cli/pkg/util"
	"github.com/spf13/cobra"
	"github.com/thoas/go-funk"
)

type SwitchOptions struct {
	Email string
}

func NewCmdSwitch(ctx *cli.Context) *cobra.Command {
	opts := new(SwitchOptions)

	cmd := &cobra.Command{
		Use:     "switch [flags]",
		Short:   "switch to another account",
		Aliases: []string{"use"},
		RunE: func(cmd *cobra.Command, args []string) error {
			ctx.SetupClient()

			return runSwitch(ctx, opts)
		},
	}

	cmd.Flags().StringVarP(&opts.Email, "user", "u", opts.Email, "your account username/email")

	return cmd
}

func runSwitch(ctx *cli.Context, opts *SwitchOptions) error {
	if opts.Email == "" {
		sessions, err := session.AllSessions()
		util.CheckErr(err)

		options := funk.Map(sessions, func(s session.Session) string { return s.Email }).([]string)
		err = ui.PromptSelect("select account", options, &opts.Email)
		util.CheckErr(err)
	}

	if opts.Email == "" {
		return errors.New("you need to login")
	}

	err := session.UseSession(opts.Email)
	util.CheckErr(err)

	fmt.Println("Session changed")
	return nil
}
