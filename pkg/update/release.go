package update

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"runtime"
)

type Release struct {
	Version     string `json:"version"`
	DownloadURL string `json:"download"`
}

type githubReleaseInfo struct {
	TagName string `json:"tag_name"`
}

func getLatestVersion(ctx context.Context) (*Release, error) {
	latestReleaseURL := fmt.Sprintf("https://api.github.com/repos/fingcloud/cli/releases/latest")
	req, err := http.NewRequestWithContext(ctx, "GET", latestReleaseURL, nil)
	if err != nil {
		return nil, err
	}

	req.Header.Add("Accept", "application/json")

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return nil, err
	}

	defer resp.Body.Close()

	var releaseInfo githubReleaseInfo
	if err := json.NewDecoder(resp.Body).Decode(&releaseInfo); err != nil {
		return nil, err
	}

	release := &Release{
		Version:     releaseInfo.TagName,
		DownloadURL: fmt.Sprintf("https://github.com/fingcloud/cli/releases/download/%s/fing-%s-%s.tar.gz", releaseInfo.TagName, runtime.GOOS, runtime.GOARCH),
	}

	return release, nil
}
