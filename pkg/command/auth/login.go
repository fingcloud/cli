package auth

import (
	"errors"
	"fmt"

	"github.com/fingcloud/cli/pkg/api"
	"github.com/fingcloud/cli/pkg/cli"
	"github.com/fingcloud/cli/pkg/command/util"
	"github.com/fingcloud/cli/pkg/config/session"
	"github.com/fingcloud/cli/pkg/ui"
	"github.com/spf13/cobra"
)

type LoginOptions struct {
	User          string
	Password      string
	PasswordStdin bool
}

func NewCmdLogin(ctx *cli.Context) *cobra.Command {
	o := new(LoginOptions)

	cmd := &cobra.Command{
		Use:   "login",
		Short: "login to fing service",
		Long:  "login to fing service",
		Run: func(cmd *cobra.Command, args []string) {
			ctx.SetupClient()

			util.CheckErr(o.Init(ctx, args))
			util.CheckErr(o.Validate())
			util.CheckErr(o.Run(ctx))
		},
	}

	cmd.Flags().StringVarP(&o.User, "user", "u", o.User, "your account username/email")
	cmd.Flags().StringVarP(&o.Password, "password", "p", o.Password, "your account password")
	cmd.Flags().BoolVar(&o.PasswordStdin, "password-stdin", false, "take the password from stdin")

	return cmd
}

func (o *LoginOptions) Init(ctx *cli.Context, args []string) error {
	if o.User == "" {
		util.CheckErr(ui.PromptEmail(&o.User))
	}

	if o.Password == "" {
		util.CheckErr(ui.PromptPassword(&o.Password))
	}

	if o.PasswordStdin {

	}
	return nil
}

func (o *LoginOptions) Validate() error {
	if o.User == "" {
		return errors.New("user/email not specified")
	}

	if o.Password == "" {
		return errors.New("password is empty")
	}

	return nil
}

func (o *LoginOptions) Run(ctx *cli.Context) error {
	auth, err := ctx.Client.AccountLogin(&api.AccountLoginOptions{
		Email:    o.User,
		Password: o.Password,
	})
	if err != nil {
		return err
	}

	fmt.Println(ui.Green("Successfully logged in."))

	sess := session.Session{
		Token: auth.Token,
		Email: auth.User.Email,
	}
	err = session.AddSession(sess)
	if err != nil {
		return err
	}

	fmt.Println(ui.Gray("Session saved"))
	return nil
}
