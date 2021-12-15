package api

import (
	"fmt"
	"net/http"
)

type AccountLoginOptions struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

type AccountRegisterOptions struct {
	Name     string `json:"name"`
	Email    string `json:"email"`
	Password string `json:"password"`
	Phone    string `json:"phone"`
}

type Auth struct {
	Token string `json:"token"`
	User  User   `json:"user"`
}

type User struct {
	Email string `json:"email"`
}

func (c *Client) AccountLogin(opts *AccountLoginOptions) (*Auth, error) {
	url := fmt.Sprintf("user/login")

	req, err := c.NewRequest(http.MethodPost, url, opts)
	if err != nil {
		return nil, err
	}

	v := new(Auth)
	_, err = c.Do(req, v)
	if err != nil {
		return nil, err
	}

	return v, err
}

func (c *Client) AccountRegister(opts *AccountRegisterOptions) (*Auth, error) {
	url := fmt.Sprintf("user/register")

	req, err := c.NewRequest(http.MethodPost, url, opts)
	if err != nil {
		return nil, err
	}

	v := new(Auth)
	_, err = c.Do(req, v)
	if err != nil {
		return nil, err
	}

	return v, err
}

type CliSession struct {
	Code string `json:"code"`
	Link string `json:"link"`
}

func (c *Client) AccountCreateCliSession() (*CliSession, error) {
	url := fmt.Sprintf("user/cli-sessions")

	req, err := c.NewRequest(http.MethodPost, url, nil)
	if err != nil {
		return nil, err
	}

	v := new(CliSession)
	_, err = c.Do(req, v)
	if err != nil {
		return nil, err
	}

	return v, err
}

func (c *Client) AccountCheckCliSession(code string) (*Auth, error) {
	url := fmt.Sprintf("user/cli-sessions/%s", code)

	req, err := c.NewRequest(http.MethodGet, url, nil)
	if err != nil {
		return nil, err
	}

	v := new(Auth)
	_, err = c.Do(req, v)
	if err != nil {
		return nil, err
	}

	return v, err
}
