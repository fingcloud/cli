package api

import (
	"bytes"
	"fmt"
	"io"
	"mime/multipart"
	"net/http"
	"time"
)

type AppStatus string

const (
	AppStatusPending   AppStatus = "pending"
	AppStatusRunning             = "running"
	AppStatusPreparing           = "preparing"
	AppStatusAssigned            = "assigned"
	AppStatusRejected            = "rejected"
	AppStatusAccepted            = "accepted"
	AppStatusReady               = "ready"
	AppStatusStarting            = "starting"
	AppStatusComplete            = "complete"
	AppStatusShutdown            = "shutdown"
	AppStatusRemoved             = "removed"
)

type ListAppsOptions struct {
}

type CreateAppOptions struct {
	Label    string `json:"label"`
	Platform string `json:"platform"`
	PlanID   int64  `json:"plan_id"`
	Region   string `json:"region"`
}

type GetAppOptions struct {
	Name string `json:"name"`
}
type StartAppOptions struct {
	Name string `json:"name"`
}

type RestartAppOptions struct {
	Name string `json:"name"`
}

type ShutdownAppOptions struct {
	Name string `json:"name"`
}

type RemoveAppOptions struct {
	Name string `json:"name"`
}

type EnvItem struct {
	Key   string `json:"key"`
	Value string `json:"value"`
}

type App struct {
	ID            int64        `json:"id"`
	Name          string       `json:"name"`
	Label         string       `json:"label"`
	Platform      string       `json:"platform"`
	Image         string       `json:"image"`
	Status        AppStatus    `json:"status"`
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

type AppLog struct {
	Stream    string `json:"stream"`
	Message   string `json:"message"`
	Timestamp int64  `json:"timestamp"`
}

type AppConfig struct {
	App      string `mapstructure:"app" json:"app"`
	Port     int    `mapstructure:"port" json:"port"`
	Platform string `mapstructure:"platform" json:"platform"`
}

type AppLogsOptions struct {
	Since int64 `json:"since"`
}

func (c *Client) AppsList(opts *ListAppsOptions) ([]*App, error) {
	url := fmt.Sprintf("apps")

	req, err := c.NewRequest(http.MethodGet, url, opts)
	if err != nil {
		return nil, err
	}

	v := make([]*App, 0)
	_, err = c.Do(req, &v)
	if err != nil {
		return nil, err
	}

	return v, err
}

func (c *Client) AppsGet(opts *GetAppOptions) (*App, error) {
	url := fmt.Sprintf("apps/%s", opts.Name)

	req, err := c.NewRequest(http.MethodGet, url, opts)
	if err != nil {
		return nil, err
	}

	v := new(App)
	_, err = c.Do(req, v)
	if err != nil {
		return nil, err
	}

	return v, err
}

func (c *Client) AppsCreate(opts *CreateAppOptions) (*App, error) {
	url := fmt.Sprintf("apps")

	req, err := c.NewRequest(http.MethodPost, url, opts)
	if err != nil {
		return nil, err
	}

	v := new(App)
	_, err = c.Do(req, v)
	if err != nil {
		return nil, err
	}

	return v, err
}

func (c *Client) AppsStart(opts *StartAppOptions) error {
	url := fmt.Sprintf("apps/%s/start", opts.Name)

	req, err := c.NewRequest(http.MethodPost, url, nil)
	if err != nil {
		return err
	}

	_, err = c.Do(req, nil)
	if err != nil {
		return err
	}

	return err
}

func (c *Client) AppsRestart(opts *RestartAppOptions) error {
	url := fmt.Sprintf("apps/%s/restart", opts.Name)

	req, err := c.NewRequest(http.MethodPost, url, nil)
	if err != nil {
		return err
	}

	_, err = c.Do(req, nil)
	if err != nil {
		return err
	}

	return err
}

func (c *Client) AppsShutdown(opts *ShutdownAppOptions) error {
	url := fmt.Sprintf("apps/%s/shutdown", opts.Name)

	req, err := c.NewRequest(http.MethodPost, url, nil)
	if err != nil {
		return err
	}

	_, err = c.Do(req, nil)
	if err != nil {
		return err
	}

	return err
}

func (c *Client) AppsRemove(opts *RemoveAppOptions) error {
	url := fmt.Sprintf("apps/%s", opts.Name)

	req, err := c.NewRequest(http.MethodDelete, url, nil)
	if err != nil {
		return err
	}

	_, err = c.Do(req, nil)
	if err != nil {
		return err
	}

	return err
}

type ProgressReader struct {
	io.Reader
	Max    int64
	Add    func(int64)
	SetMax func(int64)
}

func (r *ProgressReader) Read(p []byte) (n int, err error) {
	n, err = r.Reader.Read(p)
	r.Add(int64(n))
	return
}

func (c *Client) AppsUploadFiles(app string, tarfile io.Reader, reporter *ProgressReader) error {
	url := fmt.Sprintf("apps/%s/files", app)

	body := new(bytes.Buffer)
	m := multipart.NewWriter(body)

	file, err := m.CreateFormFile("file", "file")
	if err != nil {
		return err
	}

	_, err = io.Copy(file, tarfile)
	if err != nil {
		return err
	}

	m.Close()

	reporter.Reader = body
	reporter.SetMax(int64(body.Len()))

	req, err := c.NewRequest(http.MethodPost, url, reporter)
	if err != nil {
		return err
	}

	req.Header.Set("Content-Type", m.FormDataContentType())

	_, err = c.Do(req, nil)
	if err != nil {
		return err
	}

	return err
}

func (c *Client) AppLogs(app string, opts *AppLogsOptions) ([]*AppLog, error) {
	url := fmt.Sprintf("apps/%s/logs?since=%d", app, opts.Since)

	req, err := c.NewRequest(http.MethodGet, url, opts)
	if err != nil {
		return nil, err
	}

	v := make([]*AppLog, 0)
	_, err = c.Do(req, &v)
	if err != nil {
		return nil, err
	}

	return v, err
}
