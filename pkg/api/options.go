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

func WithUserAgent(userAgent string) Option {
	return func(c *Client) {
		c.headers["User-Agent"] = userAgent
	}
}
