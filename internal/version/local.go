package version

import (
	"fmt"
	"os"
	"strings"

	"github.com/aide-cloud/gvm/pkg/dir"
)

func FetchLocalVersions(sdkDir string) ([]string, error) {
	sdkDir = dir.ExpandHomeDir(sdkDir)
	if _, err := os.Stat(sdkDir); os.IsNotExist(err) {
		if err := os.MkdirAll(sdkDir, 0755); err != nil {
			return nil, fmt.Errorf("failed to create the sdk directory: %v", err)
		}
	}
	dis, err := os.ReadDir(sdkDir)
	if err != nil {
		return nil, fmt.Errorf("failed to read the sdk directory: %w", err)
	}
	var versions []string
	for _, di := range dis {
		if di.IsDir() {
			if dirName := di.Name(); strings.HasPrefix(dirName, "go") {
				versions = append(versions, dirName)
			}
		}
	}
	return versions, nil
}
