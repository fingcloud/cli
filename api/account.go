package api

import (
	"fmt"
	"net/http"
)

type AccountLoginOptions struct {
	Email    string `json:"email"`
	Passowrd string `json:"password"`
}

type AccountRegisterOptions struct {
	Name     string `json:"name"`
	Email    string `json:"email"`
	Passowrd string `json:"password"`
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
	url := fmt.Sprintf("users/login")

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
	url := fmt.Sprintf("users/register")

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
