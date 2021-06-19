package cmd

import (
	"fmt"

	"github.com/fingcloud/cli/api"
	"github.com/fingcloud/cli/internal/cli"
	"github.com/fingcloud/cli/internal/config"
	"github.com/fingcloud/cli/internal/ui"
	"github.com/spf13/cobra"
)

type loginOptions struct {
	user          string
	password      string
	passwordStdin bool
}

func NewLoginCommand() *cobra.Command {
	opts := new(loginOptions)

	cmd := &cobra.Command{
		Use:   "login",
		Short: "Login to Fing",
		Run: func(cmd *cobra.Command, args []string) {
			cli := cli.New(cmd, args, token, devMode)
			runLogin(cli, opts)
		},
	}

	cmd.Flags().StringVarP(&opts.user, "user", "u", "", "Your account username/email")
	cmd.Flags().StringVarP(&opts.password, "password", "p", "", "Your account password")
	cmd.Flags().BoolVar(&opts.passwordStdin, "password-stdin", false, "Take the password from stdin")

	return cmd
}

func runLogin(cli *cli.FingCli, opts *loginOptions) {
	if opts.user == "" {
		err := ui.PromptEmail(&opts.user)
		checkError(err)
	}

	if opts.password == "" {
		err := ui.PromptPassword(&opts.password)
		checkError(err)
	}

	auth, err := cli.Client.AccountLogin(&api.AccountLoginOptions{
		Email:    opts.user,
		Passowrd: opts.password,
	})
	checkError(err)

	fmt.Println(ui.Green("Successfully logged in."))

	cfg := &config.AuthConfig{
		Token: auth.Token,
		Email: auth.User.Email,
	}
	err = config.WriteAuthConfig(cfg)
	checkError(err)

	fmt.Println(ui.Gray("Session saved"))
}
