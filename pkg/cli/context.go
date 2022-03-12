package cli

import (
	"fmt"
	"io"

	"github.com/fingcloud/cli/pkg/api"
	"github.com/fingcloud/cli/pkg/config/session"
	"github.com/fingcloud/cli/pkg/ui"
	"github.com/spf13/pflag"
	"github.com/thoas/go-funk"
)

var (
	Version   = "v0.0.0-dev"
	BuildDate = "now"
	Commit    = "commit"
)

type Context struct {
	Stdout io.Writer
	Stderr io.Writer

	Client *api.Client

	APIServer   string
	AccessToken string
	Auth        string
	AskAuth     bool
	Path        string
}

func NewContext(out io.Writer, err io.Writer) *Context {
	return &Context{
		Stdout: out,
		Stderr: err,
	}
}

func (c *Context) AddFlags(flags *pflag.FlagSet) {
	flags.StringVar(&c.AccessToken, "access-token", c.AccessToken, "access token for the API server authentication")
	flags.StringVar(&c.APIServer, "server", c.APIServer, "the address of the API server")
	flags.StringVar(&c.Auth, "auth", c.Auth, "specify auth account")
	flags.BoolVar(&c.AskAuth, "use-auth", c.AskAuth, "show logged in auth to use")
}

func stringptr(val string) *string {
	return &val
}

func (c *Context) SetupClient() error {
	var accessToken string

	if c.AskAuth {
		sessions, err := session.Read()
		if err != nil {
			return err
		}

		options := funk.Map(sessions, func(s *session.Session) string { return s.Email }).([]string)
		if err := ui.PromptSelect("select account", options, &c.Auth); err != nil {
			return err
		}
	}

	if c.Auth != "" {
		if err := session.UseSession(c.Auth); err != nil {
			return err
		}
	}

	sess, _ := session.CurrentSession()
	if sess != nil {
		accessToken = sess.Token
	}

	if c.AccessToken != "" {
		accessToken = c.AccessToken
	}

	opts := make([]api.Option, 0)
	if c.APIServer != "" {
		opts = append(opts, api.WithBaseURL(c.APIServer))
		opts = append(opts, api.WithUserAgent(fmt.Sprintf("fingcli/%s", Version)))
	}

	c.Client = api.NewClient(accessToken, opts...)
	return nil
}
