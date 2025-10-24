package version

import (
	"fmt"
	"os"

	"github.com/aide-cloud/gvm/pkg/dir"
	"github.com/aide-cloud/gvm/pkg/log"
)

func Uninstall(version, sdkFilePath string) error {
	exist, err := dir.CheckFileExists(sdkFilePath)
	if err != nil {
		return fmt.Errorf("failed to check the sdk file exists: %v", err)
	}
	if !exist {
		return fmt.Errorf("sdk file %s not found", sdkFilePath)
	}
	log.Info("uninstalling version", "version", version, "sdkFilePath", sdkFilePath)
	if err := os.RemoveAll(sdkFilePath); err != nil {
		return fmt.Errorf("failed to remove the sdk file: %v", err)
	}
	log.Info("uninstalled version", "version", version, "sdkFilePath", sdkFilePath)
	return nil
}
