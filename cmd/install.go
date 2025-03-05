/*
Copyright © 2025 NAME HERE <EMAIL ADDRESS>
*/
package cmd

import (
	"fmt"
	"os"
	"path/filepath"
	"runtime"
	"strings"

	"github.com/aide-cloud/gvm/internal"
	"github.com/spf13/cobra"
)

// installCmd represents the install command
var installCmd = &cobra.Command{
	Use:   "install",
	Short: "A brief description of your command",
	Long: `A longer description that spans multiple lines and likely contains examples
and usage of using your command. For example:

Cobra is a CLI library for Go that empowers applications.
This application is a tool to generate the needed files
to quickly create a Cobra application.`,
	Run: func(cmd *cobra.Command, args []string) {
		version := buildGoVersion(args[0])
		fmt.Printf("install go version is %s\n", version)
		downloadGoVersion(version)
	},
	Args: cobra.ExactArgs(1),
	PreRun: func(cmd *cobra.Command, args []string) {
		fmt.Printf("os is %s, arch is %s\n", systemType, systemArch)
	},
}

// 系统类型
var systemType = runtime.GOOS

// 系统架构
var systemArch = runtime.GOARCH

func init() {
	rootCmd.AddCommand(installCmd)
}

// 构建下载链接
func buildDownloadUrl(version string) string {
	return fmt.Sprintf("https://mirrors.aliyun.com/golang/go%s.%s-%s.tar.gz", version, systemType, systemArch)
}

// 构建go版本
func buildGoVersion(version string) string {
	if version == "latest" {
		return getLastVersion()
	}
	version = strings.TrimPrefix(version, "v")
	version = strings.TrimPrefix(version, "V")
	version = strings.TrimPrefix(version, "go")

	return version
}

// download go version pkg
func downloadGoVersion(version string) {
	downloadUrl := buildDownloadUrl(version)
	fmt.Printf("download url is %s\n", downloadUrl)

	// 1. 下载文件
	tarGzPath := filepath.Join(rootPath, "versions", "go"+version+".tar.gz")
	err := internal.DownloadFile(downloadUrl, tarGzPath)
	if err != nil {
		fmt.Println("下载失败:", err)
		return
	}
	fmt.Println("文件下载成功:", tarGzPath)

	defer func() {
		// 3. 删除下载的 tar.gz 文件
		if err = os.Remove(tarGzPath); err != nil {
			fmt.Println("删除 tar.gz 失败:", err)
		}
	}()

	destDir := filepath.Join(rootPath, "versions")
	fileName := "go" + version

	// 2. 解压到目标目录
	targetPath := filepath.Join(destDir, fileName)
	err = internal.ExtractTarGz(tarGzPath, targetPath)
	if err != nil {
		fmt.Println("解压失败:", err)
		return
	}
	fmt.Println("解压完成:", targetPath)
	// 授予执行权限
	goBin := filepath.Join(targetPath, "go", "bin")
	if err := setPermissions(goBin); err != nil {
		fmt.Println("授予执行权限失败:", err)
		return
	}
}

// setPermissions 设置文件或目录的执行权限
func setPermissions(path string) error {
	// 获取文件的文件信息
	info, err := os.Stat(path)
	if err != nil {
		return err
	}

	// 设置权限：为文件和目录添加执行权限
	if info.IsDir() {
		// 如果是目录，递归地设置其内部文件权限
		return filepath.Walk(path, func(path string, info os.FileInfo, err error) error {
			if err != nil {
				return err
			}
			// 为文件添加执行权限
			if !info.IsDir() {
				err := os.Chmod(path, 0755) // 为文件设置可执行权限
				if err != nil {
					return fmt.Errorf("设置文件权限失败: %v", err)
				}
			}
			// 为目录添加执行权限
			if info.IsDir() {
				err := os.Chmod(path, 0755) // 为目录设置可执行权限
				if err != nil {
					return fmt.Errorf("设置目录权限失败: %v", err)
				}
			}
			return nil
		})
	} else {
		// 如果是文件，直接设置执行权限
		return os.Chmod(path, 0755)
	}
}
