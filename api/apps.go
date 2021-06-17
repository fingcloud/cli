package api

import (
	"bytes"
	"fmt"
	"io"
	"mime/multipart"
	"net/http"
	"time"
)

type ListAppsOptions struct {
}

type CreateAppOptions struct {
	Name     string `json:"name"`
	Platform string `json:"platform"`
	PlanID   int64  `json:"plan_id"`
	Region   string `json:"region"`
}

type EnvItem struct {
	Key   string `json:"key"`
	Value string `json:"value"`
}

type App struct {
	ID            int64        `json:"id"`
	Name          string       `json:"name"`
	Platform      string       `json:"platform"`
	Port          int          `json:"port"`
	Region        string       `json:"region"`
	DomainEnabled bool         `json:"domain_enabled"`
	Env           []EnvItem    `json:"env"`
	Resource      *AppResource `json:"resource"`
	CreatedAt     *time.Time   `json:"created_at"`
}

type AppResource struct {
	CPU     float32 `json:"cpu"`
	Memory  float32 `json:"memory"`
	Storage float32 `json:"storage"`
}

func (c *Client) AppsList(opts *ListAppsOptions) ([]*App, error) {
	url := fmt.Sprintf("apps")

	v := make([]*App, 0)
	req, err := c.NewRequest(http.MethodGet, url, opts)
	if err != nil {
		return nil, err
	}

	_, err = c.Do(req, v)
	if err != nil {
		return nil, err
	}

	return v, err
}

func (c *Client) AppsCreate(opts *CreateAppOptions) (*App, error) {
	url := fmt.Sprintf("apps")

	v := new(App)
	req, err := c.NewRequest(http.MethodPost, url, opts)
	if err != nil {
		return nil, err
	}

	_, err = c.Do(req, v)
	if err != nil {
		return nil, err
	}

	return v, err
}

func (c *Client) AppsUploadFiles(app string, tarfile io.Reader) error {
	url := fmt.Sprintf("apps/%s/files", app)

	body := new(bytes.Buffer)
	writer := multipart.NewWriter(body)

	file, err := writer.CreateFormFile("file", "file")
	if err != nil {
		return err
	}

	_, err = io.Copy(file, tarfile)
	if err != nil {
		return err
	}

	writer.Close()

	req, err := c.NewRequest(http.MethodPost, url, body)
	if err != nil {
		return err
	}

	req.Header.Set("Content-Type", writer.FormDataContentType())
	fmt.Println(req.Header)

	_, err = c.Do(req, nil)
	if err != nil {
		return err
	}

	return err
}
