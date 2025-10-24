package version

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
	"path/filepath"

	"github.com/aide-cloud/gvm/pkg/log"
)

type OriginVersion struct {
	Version string `json:"version"`
	Stable  bool   `json:"stable"`
	Files   []File `json:"files"`
}

type File struct {
	Filename string `json:"filename"`
	OS       string `json:"os"`
	Arch     string `json:"arch"`
	Version  string `json:"version"`
	SHA256   string `json:"sha256"`
	Size     int    `json:"size"`
	Kind     string `json:"kind"`
}

func FetchOriginVersions(originURL, versionFilePath string, forceUpdate bool) ([]OriginVersion, error) {
	if _, err := os.Stat(versionFilePath); err == nil && !forceUpdate {
		var originVersions []OriginVersion
		content, err := os.ReadFile(versionFilePath)
		if err != nil {
			return nil, fmt.Errorf("failed to read the version file: %v", err)
		}
		if len(content) > 0 {
			if err = json.Unmarshal(content, &originVersions); err != nil {
				_ = os.Remove(versionFilePath)
			}
		}
		if len(originVersions) > 0 {
			return originVersions, nil
		}
	}

	log.Info("fetching origin versions", "url", originURL)

	resp, err := http.Get(originURL)
	if err != nil {
		return nil, fmt.Errorf("failed to fetch the webpage: %v", err)
	}
	defer resp.Body.Close()
	content, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("failed to read the response body: %v", err)
	}

	var originVersions []OriginVersion
	err = json.Unmarshal(content, &originVersions)
	if err != nil {
		return nil, fmt.Errorf("failed to decode the response: %v", err)
	}

	if err := os.MkdirAll(filepath.Dir(versionFilePath), 0755); err != nil {
		return nil, fmt.Errorf("failed to create the version file directory: %v", err)
	}

	if err := os.WriteFile(versionFilePath, append(content, '\n'), 0644); err != nil {
		return nil, fmt.Errorf("failed to write the version file: %v", err)
	}
	log.Info("fetched origin versions", "versionFilePath", versionFilePath)
	return originVersions, nil
}
