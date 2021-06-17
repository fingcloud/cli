package api

import (
	"fmt"
	"io/fs"
	"net/http"
	"time"

	"github.com/fingcloud/cli/config"
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

func (c *Client) CreateDeployment(app string, opts *CreateDeploymentOptions) (*Deployment, []*FileInfo, error) {
	url := fmt.Sprintf("apps/%s/deployments", app)

	v := new(Deployment)
	req, err := c.NewRequest(http.MethodPost, url, opts)
	if err != nil {
		return nil, nil, err
	}

	resp, err := c.Do(req, v)
	if err != nil {
		if resp != nil && resp.IsFilesError() {
			return nil, resp.Files, nil
		}
		return nil, nil, err
	}

	return v, nil, nil
}
