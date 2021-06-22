package api

import (
	"net/http"
	"net/url"
)

type Option func(*Client)

func WithBaseURL(baseURL string) Option {
	return func(c *Client) {
		c.BaseURL, _ = url.Parse(baseURL)
	}
}

func WithHttpClient(client *http.Client) Option {
	return func(c *Client) {
		c.client = client
	}
}
