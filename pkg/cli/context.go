package cli

import (
	"github.com/fingcloud/cli/pkg/api"
	"github.com/fingcloud/cli/pkg/config"
	"github.com/spf13/pflag"
)

var (
	Version   = "dev"
	BuildDate = "now"
	Commit    = "commit"
)

type Context struct {
	Client *api.Client

	APIServer   *string
	AccessToken *string

	Path *string
}

func NewContext() *Context {
	return &Context{
		APIServer:   stringptr(""),
		AccessToken: stringptr(""),
		Path:        stringptr(""),
	}
}

func (c *Context) AddFlags(flags *pflag.FlagSet) {
	if c.Path != nil {
		flags.StringVar(c.Path, "path", ".", "your application path")
	}
	if c.AccessToken != nil {
		flags.StringVar(c.AccessToken, "access-token", *c.AccessToken, "access token for the API server authentication")
	}
	if c.APIServer != nil {
		flags.StringVar(c.APIServer, "server", *c.APIServer, "the address of the API server")
	}
}

func stringptr(val string) *string {
	return &val
}

func (c *Context) SetupClient() {
	var accessToken string

	cfg, _ := config.ReadAuthConfig()
	accessToken = cfg.Token

	if *c.AccessToken != "" {
		accessToken = *c.AccessToken
	}

	opts := make([]api.Option, 0)
	if *c.APIServer != "" {
		opts = append(opts, api.WithBaseURL(*c.APIServer))
	}

	c.Client = api.NewClient(accessToken, opts...)
}
