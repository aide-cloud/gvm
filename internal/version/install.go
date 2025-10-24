package version

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/aide-cloud/gvm/pkg/dir"
	"github.com/aide-cloud/gvm/pkg/download"
	"github.com/aide-cloud/gvm/pkg/log"
)

func Install(cacheFilePath, sdkFilePath, downloadFileUrl string) error {
	log.Info("checking cache file exists", "cacheFilePath", cacheFilePath)
	exist, err := dir.CheckFileExists(cacheFilePath)
	if err != nil {
		return fmt.Errorf("failed to check cache file exists: %v", err)
	}
	if !exist {
		log.Info("cache file not exists, downloading file", "cacheFilePath", cacheFilePath)

		log.Info("downloading file", "url", downloadFileUrl, "cacheFilePath", cacheFilePath)
		if err := download.FetchFile(downloadFileUrl, cacheFilePath); err != nil {
			return fmt.Errorf("failed to download file: %v", err)
		}
		log.Info("downloaded file", "url", downloadFileUrl, "cacheFilePath", cacheFilePath)
	}

	log.Info("creating sdk directory", "path", sdkFilePath)
	if err := os.MkdirAll(sdkFilePath, 0755); err != nil {
		return fmt.Errorf("failed to create sdk directory: %v", err)
	}

	log.Info("extracting file", "cacheFilePath", cacheFilePath, "sdkFilePath", sdkFilePath)
	if err := download.ExtractGoSdkTarGzFile(cacheFilePath, sdkFilePath); err != nil {
		return fmt.Errorf("failed to extract tar.gz file: %v", err)
	}
	log.Info("installed version", "sdkFilePath", sdkFilePath)

	// 递归设置权限
	binFilePath := filepath.Join(sdkFilePath, "bin")
	if err := setPermissionsRecursively(binFilePath, 0755); err != nil {
		return fmt.Errorf("failed to set permissions to bin directory: %v", err)
	}
	log.Info("set permissions to bin directory", "binFilePath", binFilePath)

	// 设置 pkg/tool 目录的权限（Go 工具需要执行权限）
	toolDirPath := filepath.Join(sdkFilePath, "pkg", "tool")
	if err := setPermissionsRecursively(toolDirPath, 0755); err != nil {
		return fmt.Errorf("failed to set permissions to tool directory: %v", err)
	}
	log.Info("set permissions to tool directory", "toolDirPath", toolDirPath)
	return nil
}

// setPermissionsRecursively 递归设置目录和文件的权限
func setPermissionsRecursively(rootPath string, mode os.FileMode) error {
	return filepath.Walk(rootPath, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}

		// 设置权限
		if err := os.Chmod(path, mode); err != nil {
			return fmt.Errorf("failed to set permissions for %s: %v", path, err)
		}

		return nil
	})
}
