package version

import (
	"fmt"
	"net/url"
	"os"
	"path/filepath"
	"runtime"
	"slices"
	"strings"

	"github.com/aide-cloud/gvm/pkg/dir"
	"github.com/aide-cloud/gvm/pkg/log"
)

const (
	DefaultSdkDir               = "~/go/sdk"
	DefaultCacheDir             = "~/.gvm/cache"
	DefaultOriginURL            = "https://go.dev/dl/?mode=json&include=all"
	DefaultDownloadURL          = "https://dl.google.com/go/"
	DefaultVersionFilePath      = "~/.gvm/versions.json"
	DefaultLocalVersionFilePath = "~/.gvm/version"
)

type Version struct {
	sdkDir               string
	cacheDir             string
	versionFilePath      string
	localVersionFilePath string
	originURL            string
	downloadURL          string
}

type VersionOption func(*Version)

func NewVersion(opts ...VersionOption) *Version {
	v := &Version{
		sdkDir:               dir.ExpandHomeDir(DefaultSdkDir),
		cacheDir:             dir.ExpandHomeDir(DefaultCacheDir),
		versionFilePath:      dir.ExpandHomeDir(DefaultVersionFilePath),
		localVersionFilePath: dir.ExpandHomeDir(DefaultLocalVersionFilePath),
		originURL:            DefaultOriginURL,
		downloadURL:          DefaultDownloadURL,
	}
	for _, opt := range opts {
		opt(v)
	}
	// 检查目录是否存在，如果不存在，则创建
	if _, err := os.Stat(v.sdkDir); os.IsNotExist(err) {
		if err := os.MkdirAll(v.sdkDir, 0755); err != nil {
			log.Error("Failed to create sdk directory:", "error", err)
		}
	}
	if _, err := os.Stat(v.cacheDir); os.IsNotExist(err) {
		if err := os.MkdirAll(v.cacheDir, 0755); err != nil {
			log.Error("Failed to create cache directory:", "error", err)
		}
	}
	versionFileDir := filepath.Dir(v.versionFilePath)
	if _, err := os.Stat(versionFileDir); os.IsNotExist(err) {
		if err := os.MkdirAll(versionFileDir, 0755); err != nil {
			log.Error("Failed to create version file directory:", "error", err)
		}
	}
	if _, err := os.Stat(v.localVersionFilePath); os.IsNotExist(err) {
		if err := os.WriteFile(v.localVersionFilePath, []byte(""), 0644); err != nil {
			log.Error("Failed to create local version file:", "error", err)
		}
	}
	// 检查sdk目前权限，如果权限不正确，则设置为755
	if err := setPermissionsRecursively(v.sdkDir, 0755); err != nil {
		log.Error("Failed to set permissions to sdk directory:", "error", err)
	}
	return v
}

func (v *Version) Use(targetVersion string, isForce, isEval bool) {
	version, err := v.getOriginVersion(targetVersion, false)
	if err != nil {
		log.Error("Failed to get version:", "error", err)
		return
	}
	exist, err := v.checkLocalVersion(version)
	if err != nil {
		log.Error("Failed to check local version:", "error", err)
		return
	}

	tarGzFilename := v.tarGzFilename(version)
	cacheFilePath := v.cacheFilePath(tarGzFilename)
	sdkFilePath := v.sdkFilePath(version)
	downloadFileUrl, err := url.JoinPath(v.downloadURL, tarGzFilename)
	if err != nil {
		log.Error("Failed to join download file url", "error", err)
		return
	}

	if !exist || isForce {
		if isForce && exist {
			_ = os.RemoveAll(cacheFilePath)
			_ = os.RemoveAll(sdkFilePath)
		}
		if err := Install(cacheFilePath, sdkFilePath, downloadFileUrl); err != nil {
			log.Error("Failed to install version:", "error", err)
			return
		}
	}
	if err := Use(version, v.sdkDir, v.localVersionFilePath, isEval); err != nil {
		log.Error("Failed to use version:", "error", err)
		return
	}
}

func (v *Version) Install(targetVersion string, isForce bool) {
	version, err := v.getOriginVersion(targetVersion, false)
	if err != nil {
		log.Error("Failed to get version:", "error", err)
		return
	}
	exist, err := v.checkLocalVersion(version)
	if err != nil {
		log.Error("Failed to check local version:", "error", err)
		return
	}
	tarGzFilename := v.tarGzFilename(version)
	cacheFilePath := v.cacheFilePath(tarGzFilename)
	sdkFilePath := v.sdkFilePath(version)
	downloadFileUrl, err := url.JoinPath(v.downloadURL, tarGzFilename)
	if err != nil {
		log.Error("Failed to join download file url", "error", err)
		return
	}

	if exist && !isForce {
		log.Info("Version already installed", "version", version)
		return
	}

	if isForce && exist {
		_ = os.RemoveAll(cacheFilePath)
		_ = os.RemoveAll(sdkFilePath)
	}

	if err := Install(cacheFilePath, sdkFilePath, downloadFileUrl); err != nil {
		log.Error("Failed to install version:", "error", err)
		return
	}
}

func (v *Version) Uninstall(targetVersion string) {
	version, err := v.getOriginVersion(targetVersion, false)
	if err != nil {
		log.Error("Failed to get version:", "error", err)
		return
	}
	exist, err := v.checkLocalVersion(version)
	if err != nil {
		log.Error("Failed to check local version:", "error", err)
		return
	}
	if !exist {
		log.Error("Version not installed", "version", version)
		return
	}
	if err := Uninstall(version, v.sdkFilePath(version)); err != nil {
		log.Error("Failed to uninstall version:", "error", err)
		return
	}
}

func (v *Version) Ls() {
	vs, err := FetchLocalVersions(v.sdkDir)
	if err != nil {
		log.Error("Failed to fetch local versions:", "error", err)
		return
	}
	if len(vs) == 0 {
		log.Info("No local versions found")
		return
	}
	// 读取本地版本文件
	content, _ := os.ReadFile(v.localVersionFilePath)
	localVersion := string(content)
	for _, v := range vs {
		if v == localVersion {
			fmt.Println("*", v)
		} else {
			fmt.Println(" ", v)
		}
	}
}

func (v *Version) List(isLatest bool, showNumber int, forceUpdate bool) {
	originVersions, err := FetchOriginVersions(v.originURL, v.versionFilePath, forceUpdate)
	if err != nil {
		log.Error("Failed to fetch origin versions:", "error", err)
		return
	}
	if len(originVersions) == 0 {
		log.Info("No origin versions found")
		return
	}
	if isLatest {
		fmt.Println(originVersions[0].Version)
		return
	}
	versions := make([]string, 0, len(originVersions))
	for index, o := range originVersions {
		if index >= showNumber {
			break
		}
		versions = append(versions, o.Version)
	}
	fmt.Println(strings.Join(versions, "\n"))
}

func (v *Version) getOriginVersion(targetVersion string, forceUpdate bool) (string, error) {
	vs, err := FetchOriginVersions(v.originURL, v.versionFilePath, forceUpdate)
	if err != nil {
		return "", fmt.Errorf("failed to fetch origin versions: %v", err)
	}
	if len(vs) == 0 {
		return "", fmt.Errorf("no origin versions found")
	}
	if targetVersion == "latest" {
		return vs[0].Version, nil
	}
	suspiciousVersion := make([]string, 0, len(vs))
	for _, o := range vs {
		if o.Version == targetVersion {
			return o.Version, nil
		}
		if strings.Contains(o.Version, targetVersion) {
			suspiciousVersion = append(suspiciousVersion, o.Version)
		}
	}
	if len(suspiciousVersion) == 0 {
		return "", fmt.Errorf("version %s not found", targetVersion)
	}
	return suspiciousVersion[0], nil
}

func (v *Version) checkLocalVersion(targetVersion string) (bool, error) {
	vs, err := FetchLocalVersions(v.sdkDir)
	if err != nil {
		return false, fmt.Errorf("failed to fetch local versions: %v", err)
	}
	if slices.Contains(vs, targetVersion) {
		return true, nil
	}
	return false, nil
}

func (v *Version) tarGzFilename(version string) string {
	return fmt.Sprintf("%s.%s-%s.tar.gz", version, runtime.GOOS, runtime.GOARCH)
}

func (v *Version) cacheFilePath(tarGzFilename string) string {
	return filepath.Join(v.cacheDir, tarGzFilename)
}

func (v *Version) sdkFilePath(version string) string {
	return filepath.Join(v.sdkDir, version)
}

func WithSdkDir(sdkDir string) VersionOption {
	return func(v *Version) {
		v.sdkDir = dir.ExpandHomeDir(sdkDir)
	}
}

func WithCacheDir(cacheDir string) VersionOption {
	return func(v *Version) {
		v.cacheDir = dir.ExpandHomeDir(cacheDir)
	}
}

func WithOriginURL(originURL string) VersionOption {
	return func(v *Version) {
		v.originURL = originURL
	}
}

func WithDownloadURL(downloadURL string) VersionOption {
	return func(v *Version) {
		v.downloadURL = downloadURL
	}
}
