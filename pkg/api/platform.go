package api

import (
	"fmt"
	"net/http"
)

type Platform struct {
	ID   string `json:"id"`
	Name string `json:"name"`
}

func (c *Client) PlatformList() ([]*Platform, error) {
	url := fmt.Sprintf("platforms")

	req, err := c.NewRequest(http.MethodGet, url, nil)
	if err != nil {
		return nil, err
	}

	v := make([]*Platform, 0)
	_, err = c.Do(req, &v)
	if err != nil {
		return nil, err
	}

	return v, err
}
