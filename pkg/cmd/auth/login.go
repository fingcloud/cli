package auth

import (
	"errors"
	"fmt"

	"github.com/fingcloud/cli/pkg/api"
	"github.com/fingcloud/cli/pkg/cli"
	"github.com/fingcloud/cli/pkg/config/session"
	"github.com/fingcloud/cli/pkg/ui"
	"github.com/fingcloud/cli/pkg/util"
	"github.com/spf13/cobra"
)

type LoginOptions struct {
	User          string
	Password      string
	PasswordStdin bool
}

func NewCmdLogin(ctx *cli.Context) *cobra.Command {
	opts := new(LoginOptions)

	cmd := &cobra.Command{
		Use:     "login",
		Short:   "login your account",
		Aliases: []string{"add"},
		Run: func(cmd *cobra.Command, args []string) {
			ctx.SetupClient()

			util.CheckErr(runLogin(ctx, opts))
		},
	}

	cmd.Flags().StringVarP(&opts.User, "user", "u", opts.User, "your account username/email")
	cmd.Flags().StringVarP(&opts.Password, "password", "p", opts.Password, "your account password")
	cmd.Flags().BoolVar(&opts.PasswordStdin, "password-stdin", false, "take the password from stdin")

	return cmd
}

func (opts *LoginOptions) validate() error {
	if opts.User == "" {
		return errors.New("user/email not specified")
	}

	if opts.Password == "" {
		return errors.New("password is empty")
	}

	return nil
}

func runLogin(ctx *cli.Context, opts *LoginOptions) error {
	if opts.User == "" {
		util.CheckErr(ui.PromptEmail(&opts.User))
	}

	if opts.Password == "" {
		util.CheckErr(ui.PromptPassword(&opts.Password))
	}

	if opts.PasswordStdin {

	}

	util.CheckErr(opts.validate())

	auth, err := ctx.Client.AccountLogin(&api.AccountLoginOptions{
		Email:    opts.User,
		Password: opts.Password,
	})
	util.CheckErr(err)

	fmt.Println(ui.Green("Successfully logged in."))

	sess := &session.Session{
		Token: auth.Token,
		Email: auth.User.Email,
	}
	err = session.AddSession(sess)
	util.CheckErr(err)

	fmt.Println("Session saved")
	return nil
}
