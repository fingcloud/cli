package auth

import (
	"errors"
	"fmt"
	"io/ioutil"
	"os"
	"strings"

	"github.com/fingcloud/cli/pkg/api"
	"github.com/fingcloud/cli/pkg/cli"
	"github.com/fingcloud/cli/pkg/config/session"
	"github.com/fingcloud/cli/pkg/ui"
	"github.com/fingcloud/cli/pkg/util"
	"github.com/spf13/cobra"
)

type LoginOptions struct {
	Username      string
	Password      string
	PasswordStdin bool
	Browser       bool
}

func NewCmdLogin(ctx *cli.Context) *cobra.Command {
	opts := new(LoginOptions)

	cmd := &cobra.Command{
		Use:     "login",
		Short:   "login your account",
		Aliases: []string{"add"},
		Run: func(cmd *cobra.Command, args []string) {
			loginHandler := cliLogin
			if opts.Browser {
				loginHandler = browserLogin
			}
			util.CheckErr(runLogin(ctx, opts, loginHandler))
		},
	}

	cmd.Flags().StringVarP(&opts.Username, "username", "u", opts.Username, "your account username/email")
	cmd.Flags().StringVarP(&opts.Password, "password", "p", opts.Password, "your account password")
	cmd.Flags().BoolVar(&opts.PasswordStdin, "password-stdin", false, "take the password from stdin")
	cmd.Flags().BoolVar(&opts.Browser, "browser", false, "login easily using your browser")

	return cmd
}

func (opts *LoginOptions) validate() error {
	if opts.Username == "" {
		return errors.New("username not specified")
	}

	if opts.Password == "" {
		return errors.New("password is empty")
	}

	return nil
}

func getCredentials(opts *LoginOptions) error {
	// warn user if uses --password in none development environment
	if opts.Password != "" {
		fmt.Println(ui.Yellow("WARNING! Using --password via the CLI is insecure. Use --password-stdin instead."))
		if opts.PasswordStdin {
			return errors.New("--password and --password-stdin can't be used together")
		}
	}

	// read password from stdin --password-stdin
	if opts.PasswordStdin {
		if opts.Username == "" {
			return errors.New("Must provider --username with --password-stdin")
		}

		input, err := ioutil.ReadAll(os.Stdin)
		util.CheckErr(err)

		opts.Password = strings.TrimRight(string(input), "\n\r")
		return nil
	}

	if opts.Username == "" {
		util.CheckErr(ui.PromptEmail(&opts.Username))
	}

	if opts.Password == "" {
		util.CheckErr(ui.PromptPassword(&opts.Password))
	}

	return nil
}

func cliLogin(ctx *cli.Context, opts *LoginOptions) (*api.Auth, error) {
	if err := getCredentials(opts); err != nil {
		return nil, err
	}

	if err := opts.validate(); err != nil {
		return nil, err
	}

	return ctx.Client.AccountLogin(&api.AccountLoginOptions{
		Email:    opts.Username,
		Password: opts.Password,
	})
}

type LoginHandler func(*cli.Context, *LoginOptions) (*api.Auth, error)

func runLogin(ctx *cli.Context, opts *LoginOptions, loginHandler LoginHandler) (err error) {
	auth, err := loginHandler(ctx, opts)
	if err != nil {
		return err
	}

	fmt.Println(ui.Green("Successfully logged in."))

	sess := &session.Session{
		Token: auth.Token,
		Email: auth.User.Email,
	}
	err = session.AddSession(sess)
	if err != nil {
		return err
	}

	fmt.Println("Session saved")
	return nil
}
