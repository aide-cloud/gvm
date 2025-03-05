/*
Copyright © 2025 NAME HERE <EMAIL ADDRESS>
*/
package cmd

import (
	"fmt"
	"os"
	"os/exec"
	"regexp"

	"github.com/aide-cloud/gvm/internal"
	"github.com/spf13/cobra"
)

// lsCmd represents the ls command
var lsCmd = &cobra.Command{
	Use:   "ls",
	Short: "A brief description of your command",
	Long: `A longer description that spans multiple lines and likely contains examples
and usage of using your command. For example:

Cobra is a CLI library for Go that empowers applications.
This application is a tool to generate the needed files
to quickly create a Cobra application.`,
	Run: func(cmd *cobra.Command, args []string) {
		ls()
	},
}

func init() {
	rootCmd.AddCommand(lsCmd)
}

func ls() {
	// 从 .zshrc 文件中提取 GOROOT 路径
	goroot, err := extractGORootFromZshrc()
	if err != nil {
		// 如果无法从 .zshrc 提取版本号，则执行 go version 获取版本
		version, err := getGoVersion()
		if err != nil {
			fmt.Println("获取 Go 版本失败:", err)
			return
		}
		fmt.Printf("go%s\n", version)
		return
	}

	// 尝试从 GOROOT 路径解析版本号
	version, err := extractVersionFromPath(goroot)
	if err != nil {
		// 如果路径中没有版本号，则执行 go version 获取版本
		version, err = getGoVersion()
		if err != nil {
			fmt.Println("获取 Go 版本失败:", err)
			return
		}
	}
	// 列举当前本地安装的 Go 版本， 当前版本显示* gox.x.x
	versions := internal.FetchLocalVersions()

	showVersions(versions, version)
}

// 从 ~/.zshrc 文件中提取 GOROOT 配置
func extractGORootFromZshrc() (string, error) {
	// 打开 ~/.zshrc 文件
	zshrcPath := getZshrcPath()
	file, err := os.Open(zshrcPath)
	if err != nil {
		return "", fmt.Errorf("无法打开文件 %s: %v", zshrcPath, err)
	}
	defer file.Close()

	// 读取原始内容
	contentBytes, err := os.ReadFile(zshrcPath)
	if err != nil {
		// 如果文件不存在，则创建
		if os.IsNotExist(err) {
			contentBytes = []byte{}
		} else {
			return "", err
		}
	}
	content := string(contentBytes)
	// 先检查是否已有 GOROOT 变量
	re := regexp.MustCompile(`(?m)^export GOROOT=.*$`)
	if re.MatchString(content) {
		// 提取 GOROOT 路径
		gorootRe := regexp.MustCompile(`(?m)^export GOROOT=(.*)$`)
		gorootMatches := gorootRe.FindStringSubmatch(content)
		if len(gorootMatches) > 1 {
			return gorootMatches[1], nil
		}
	}

	return "", nil
}

// 从 GOROOT 路径中解析版本号
func extractVersionFromPath(goroot string) (string, error) {
	// 使用正则表达式从路径中提取 Go 版本
	re := regexp.MustCompile(`go(\d+\.\d+\.\d+)`)
	matches := re.FindStringSubmatch(goroot)
	if len(matches) > 1 {
		return matches[1], nil
	}
	return "", fmt.Errorf("无法从路径中解析版本")
}

// 使用 go version 获取当前安装的 Go 版本
func getGoVersion() (string, error) {
	// 执行 go version 命令
	cmd := exec.Command("go", "version")
	output, err := cmd.CombinedOutput()
	if err != nil {
		return "", fmt.Errorf("执行 go version 失败: %v", err)
	}
	// 提取版本号
	re := regexp.MustCompile(`go(\d+\.\d+\.\d+)`)
	matches := re.FindStringSubmatch(string(output))
	if len(matches) > 1 {
		return matches[1], nil
	}
	return "", fmt.Errorf("无法解析 go version 输出")
}
