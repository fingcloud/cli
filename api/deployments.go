package api

import (
	"fmt"
	"io/fs"
	"net/http"
	"time"

	"github.com/fingcloud/cli/internal/config"
)

type FileInfo struct {
	Path     string      `json:"path"`
	Checksum string      `json:"checksum"`
	Size     int64       `json:"size"`
	Dir      bool        `json:"dir"`
	Mode     fs.FileMode `json:"mode"`
}

type CreateDeploymentOptions struct {
	Files  []*FileInfo    `json:"files"`
	Config *config.Config `json:"config"`
}

type Deployment struct {
	ID        int64            `json:"id"`
	Platform  string           `json:"platform"`
	Image     string           `json:"image"`
	Port      int              `json:"port"`
	Status    DeploymentStatus `json:"status"`
	CreatedAt *time.Time       `json:"created_at"`
}

type ListLogsOptions struct {
	From int64 `json:"from"`
}

type BuildLog struct {
	ID        int64     `json:"id"`
	Level     string    `json:"level"`
	Message   string    `json:"message"`
	CreatedAt time.Time `json:"created_at"`
}

type BuildLogs struct {
	Deployment *Deployment `json:"deployment"`
	Logs       []*BuildLog `json:"logs"`
}

type DeploymentLog struct {
	Stream    string    `json:"stream"`
	Message   string    `json:"message"`
	Timestamp time.Time `json:"timestamp"`
}

type DeploymentStatus string

const (
	DeploymentStatusReady    = "ready"
	DeploymentStatusPending  = "pending"
	DeploymentStatusBuilding = "building"
	DeploymentStatusStarting = "starting"
	DeploymentStatusRunning  = "running"
	DeploymentStatusShutdown = "shutdown"
	DeploymentStatusCancel   = "cancel"
	DeploymentStatusFailed   = "failed"
)

func (c *Client) DeployemntCreate(app string, opts *CreateDeploymentOptions) (*Deployment, []*FileInfo, error) {
	url := fmt.Sprintf("apps/%s/deployments", app)

	req, err := c.NewRequest(http.MethodPost, url, opts)
	if err != nil {
		return nil, nil, err
	}

	v := new(Deployment)
	resp, err := c.Do(req, v)
	if err != nil {
		if resp != nil && resp.IsFilesError() {
			return nil, resp.Files, nil
		}
		return nil, nil, err
	}

	return v, nil, nil
}

func (c *Client) DeploymentListBuildLogs(app string, deploymentId int64, opts *ListLogsOptions) (*BuildLogs, error) {
	url := fmt.Sprintf("apps/%s/deployments/%d/build-logs?from=%d", app, deploymentId, opts.From)

	req, err := c.NewRequest(http.MethodGet, url, opts)
	if err != nil {
		return nil, err
	}

	v := new(BuildLogs)
	_, err = c.Do(req, v)
	if err != nil {
		return nil, err
	}

	return v, err
}

func (c *Client) DeploymentListLogs(app string, deploymentId int64, opts *ListLogsOptions) ([]*DeploymentLog, error) {
	url := fmt.Sprintf("apps/%s/deployments/%d/logs?from=%d", app, deploymentId, opts.From)

	req, err := c.NewRequest(http.MethodGet, url, opts)
	if err != nil {
		return nil, err
	}

	v := make([]*DeploymentLog, 0)
	_, err = c.Do(req, v)
	if err != nil {
		return nil, err
	}

	return v, err
}

func (c *Client) DeploymentCancel(app string, deploymentId int64) error {
	url := fmt.Sprintf("apps/%s/deployments/%d/cancel", app, deploymentId)

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
