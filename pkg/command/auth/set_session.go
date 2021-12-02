package auth

import (
	"errors"
	"fmt"

	"github.com/AlecAivazis/survey/v2/terminal"
	"github.com/fingcloud/cli/pkg/cli"
	"github.com/fingcloud/cli/pkg/command/util"
	"github.com/fingcloud/cli/pkg/config/session"
	"github.com/fingcloud/cli/pkg/ui"
	"github.com/spf13/cobra"
	"github.com/thoas/go-funk"
)

type SetSessionOptions struct {
	Email string
}

func NewCmdSetSession(ctx *cli.Context) *cobra.Command {
	o := new(SetSessionOptions)

	cmd := &cobra.Command{
		Use:   "set-session",
		Short: "set default session",
		Long:  "use set-session to switch between your logged in accounts",
		Run: func(cmd *cobra.Command, args []string) {
			ctx.SetupClient()

			util.CheckErr(o.Init(ctx, args))
			util.CheckErr(o.Validate())
			util.CheckErr(o.Run(ctx))
		},
	}

	cmd.Flags().StringVarP(&o.Email, "user", "u", o.Email, "your account username/email")

	return cmd
}

func (o *SetSessionOptions) Init(ctx *cli.Context, args []string) error {
	if o.Email == "" {
		sessions, err := session.Sessions()
		if err != nil {
			return err
		}

		options := funk.Map(sessions, func(s session.Session) string { return s.Email }).([]string)
		err = ui.PromptSelect("choose session", options, &o.Email)
		if err == terminal.InterruptErr {
			return terminal.InterruptErr
		}
	}

	return nil
}

func (o *SetSessionOptions) Validate() error {
	if o.Email == "" {
		return errors.New("you need to login")
	}
	return nil
}

func (o *SetSessionOptions) Run(ctx *cli.Context) error {

	err := session.UseSession(o.Email)
	if err != nil {
		return err
	}

	fmt.Println(ui.Gray("Session changed"))
	return nil
}
