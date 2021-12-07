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

type UseOptions struct {
	Email string
}

func NewCmdUse(ctx *cli.Context) *cobra.Command {
	opts := new(UseOptions)

	cmd := &cobra.Command{
		Use:     "use [flags]",
		Short:   "use another account",
		Aliases: []string{"switch", "change"},
		Run: func(cmd *cobra.Command, args []string) {
			ctx.SetupClient()

			util.CheckErr(RunUse(ctx, opts))
		},
	}

	cmd.Flags().StringVarP(&opts.Email, "user", "u", opts.Email, "your account username/email")

	return cmd
}

func RunUse(ctx *cli.Context, opts *UseOptions) error {
	if opts.Email == "" {
		sessions, err := session.Read()
		util.CheckErr(err)

		options := funk.Map(sessions, func(s *session.Session) string { return s.Email }).([]string)
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
