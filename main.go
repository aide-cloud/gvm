package main

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
	"path/filepath"
	"runtime"
	"strings"

	"github.com/spf13/cobra"
)

const (
	GVM_DIR          = ".gvm"
	VERSIONS_DIR     = "versions"
	CURRENT_VERSION  = "current"
	GO_VERSIONS_API  = "https://golang.org/dl/?mode=json"
	DEFAULT_OS       = runtime.GOOS
	DEFAULT_ARCH     = runtime.GOARCH
)

var rootCmd = &cobra.Command{
	Use:   "gvm",
	Short: "Go Version Manager",
	Long:  "A simple Go Version Manager inspired by nvm",
}

var installCmd = &cobra.Command{
	Use:   "install [version]",
	Short: "Install a Go version",
	Args:  cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		version := args[0]
		installVersion(version)
	},
}

var useCmd = &cobra.Command{
	Use:   "use [version]",
	Short: "Switch to use specified Go version",
	Args:  cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		version := args[0]
		setCurrentVersion(version)
	},
}

var listCmd = &cobra.Command{
	Use:   "list",
	Short: "List installed Go versions",
	Run: func(cmd *cobra.Command, args []string) {
		listInstalledVersions()
	},
}

var listRemoteCmd = &cobra.Command{
	Use:   "list-remote",
	Short: "List available remote Go versions",
	Run: func(cmd *cobra.Command, args []string) {
		listRemoteVersions()
	},
}

var uninstallCmd = &cobra.Command{
	Use:   "uninstall [version]",
	Short: "Uninstall a Go version",
	Args:  cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		version := args[0]
		uninstallVersion(version)
	},
}

func main() {
	rootCmd.AddCommand(installCmd, useCmd, listCmd, listRemoteCmd, uninstallCmd)
	rootCmd.Execute()
}

// 安装指定版本
func installVersion(version string) {
	versions := getRemoteVersions()
	found := false
	var targetVersion *GoVersion

	for _, v := range versions {
		if v.Version == "go"+version {
			targetVersion = &v
			found = true
			break
		}
	}

	if !found {
		fmt.Printf("Version %s not found\n", version)
		return
	}

	versionDir := getVersionDir(version)
	if _, err := os.Stat(versionDir); !os.IsNotExist(err) {
		fmt.Printf("Version %s already installed\n", version)
		return
	}

	downloadURL := ""
	for _, file := range targetVersion.Files {
		if file.OS == DEFAULT_OS && file.Arch == DEFAULT_ARCH && file.Kind == "archive" {
			downloadURL = file.URL
			break
		}
	}

	if downloadURL == "" {
		fmt.Println("No suitable package found for your OS/Arch")
		return
	}

	fmt.Printf("Downloading %s...\n", downloadURL)
	if err := downloadAndExtract(downloadURL, versionDir); err != nil {
		fmt.Printf("Installation failed: %v\n", err)
		return
	}

	fmt.Printf("Successfully installed go %s\n", version)
}

// 设置当前使用版本
func setCurrentVersion(version string) {
	versionDir := getVersionDir(version)
	if _, err := os.Stat(versionDir); os.IsNotExist(err) {
		fmt.Printf("Version %s is not installed\n", version)
		return
	}

	currentPath := filepath.Join(getGvmDir(), CURRENT_VERSION)
	os.Remove(currentPath)
	
	if err := os.Symlink(versionDir, currentPath); err != nil {
		fmt.Printf("Failed to switch version: %v\n", err)
		return
	}

	fmt.Printf("Now using go %s\n", version)
	fmt.Println("Please run:")
	fmt.Printf("export PATH=%s/bin:$PATH\n", currentPath)
}

// 获取远程版本列表
func listRemoteVersions() {
	versions := getRemoteVersions()
	for _, v := range versions {
		fmt.Println(strings.TrimPrefix(v.Version, "go"))
	}
}

// 列出已安装版本
func listInstalledVersions() {
	versionsDir := getVersionsDir()
	files, err := os.ReadDir(versionsDir)
	if err != nil {
		return
	}

	for _, f := range files {
		if f.IsDir() {
			fmt.Println(f.Name())
		}
	}
}

// 卸载指定版本
func uninstallVersion(version string) {
	versionDir := getVersionDir(version)
	if err := os.RemoveAll(versionDir); err != nil {
		fmt.Printf("Uninstall failed: %v\n", err)
		return
	}
	fmt.Printf("Successfully uninstalled go %s\n", version)
}

// 辅助函数
func getGvmDir() string {
	home, _ := os.UserHomeDir()
	return filepath.Join(home, GVM_DIR)
}

func getVersionsDir() string {
	return filepath.Join(getGvmDir(), VERSIONS_DIR)
}

func getVersionDir(version string) string {
	return filepath.Join(getVersionsDir(), version)
}

type GoVersion struct {
	Version string     `json:"version"`
	Files   []FileInfo `json:"files"`
}

type FileInfo struct {
	Filename string `json:"filename"`
	OS       string `json:"os"`
	Arch     string `json:"arch"`
	Version  string `json:"version"`
	Checksum string `json:"sha256"`
	Size     int64  `json:"size"`
	Kind     string `json:"kind"`
	URL      string `json:"url"`
}

// 获取远程版本信息
func getRemoteVersions() []GoVersion {
	resp, err := http.Get(GO_VERSIONS_API)
	if err != nil {
		return nil
	}
	defer resp.Body.Close()

	var versions []GoVersion
	json.NewDecoder(resp.Body).Decode(&versions)
	return versions
}

// 下载并解压
func downloadAndExtract(url, targetDir string) error {
	resp, err := http.Get(url)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	tmpFile := filepath.Join(os.TempDir(), filepath.Base(url))
	out, err := os.Create(tmpFile)
	if err != nil {
		return err
	}
	defer out.Close()

	_, err = io.Copy(out, resp.Body)
	if err != nil {
		return err
	}

	// 这里需要添加解压逻辑（根据文件类型使用tar或unzip）
	// 伪代码示例：
	// if strings.HasSuffix(tmpFile, ".tar.gz") {
	//  执行tar解压到targetDir
	// }

	return os.Rename(tmpFile, targetDir)
}