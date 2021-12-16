package auth

import (
	"context"
	"errors"
	"log"
	"net/http"
	"time"

	"github.com/fingcloud/cli/pkg/api"
	"github.com/fingcloud/cli/pkg/cli"
	"github.com/pkg/browser"
	"github.com/r6m/spinner"
)

func browserLogin(ctx *cli.Context, opts *LoginOptions) (*api.Auth, error) {
	s := spinner.New().WithOptions(spinner.WithExitOnAbort(false), spinner.WithNotifySignals(false))
	s.Start("Getting Session...")

	cliSession, err := ctx.Client.AccountCreateCliSession()
	if err != nil {
		return nil, err
	}

	err = browser.OpenURL(cliSession.Link)
	if err != nil {
		log.Println("Could not open the browser, visit the following link in your browser")
		log.Printf("\t%s\n", cliSession.Link)
	}

	s.Success("Session OK")

	timeoutCtx, cancel := context.WithTimeout(context.Background(), 5*time.Minute)
	defer cancel()
	s.Start("Getting Authentication ...")

	for timeoutCtx.Err() == nil {
		auth, err := ctx.Client.AccountCheckCliSession(cliSession.Code)
		if err != nil {
			if resp, ok := err.(*api.Response); ok && resp.StatusCode == http.StatusBadRequest {
				time.Sleep(time.Second)
				continue
			}
			return nil, err
		}

		s.Success("Authentication OK")

		return auth, nil
	}

	return nil, errors.New("login failed")
}
