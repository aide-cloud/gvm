/*
Copyright © 2025 NAME HERE <EMAIL ADDRESS>
*/
package cmd

import (
	"fmt"
	"os"
	"os/exec"
	"os/user"
	"path/filepath"
	"regexp"
	"strings"

	"github.com/spf13/cobra"
)

// useCmd represents the use command
var useCmd = &cobra.Command{
	Use:   "use",
	Short: "A brief description of your command",
	Long: `A longer description that spans multiple lines and likely contains examples
and usage of using your command. For example:

Cobra is a CLI library for Go that empowers applications.
This application is a tool to generate the needed files
to quickly create a Cobra application.`,
	Run: func(cmd *cobra.Command, args []string) {
		version := args[0]
		useGoVersion(version)
	},
	Args: cobra.ExactArgs(1),
}

func init() {
	rootCmd.AddCommand(useCmd)
}

func useGoVersion(version string) {
	version = buildGoVersion(version)
	version = "go" + version
	// 检查版本目录是否存在
	if _, err := os.Stat(filepath.Join(rootPath, "versions", version)); os.IsNotExist(err) {
		fmt.Println("版本不存在:", version)
		return
	}
	fmt.Println("use go version:", version)
	// 更新环境变量
	if err := updateZshrc(filepath.Join(rootPath, "versions", version, "go")); err != nil {
		fmt.Println("更新环境变量失败:", err)
	}
}

func getZshrcPath() string {
	filename := ".zshrc"
	usr, err := user.Current()
	if err != nil {
		panic(err)
	}
	return filepath.Join(usr.HomeDir, filename)
}

// updateZshrc 更新 ~/.zshrc 文件中的 GOROOT 设置
func updateZshrc(newGOROOT string) error {
	filename := ".zshrc"

	zshrcPath := getZshrcPath()

	// 如果文件不存在则创建文件
	if _, err := os.Stat(zshrcPath); os.IsNotExist(err) {
		file, err := os.Create(zshrcPath)
		if err != nil {
			return err
		}
		defer file.Close()
	}
	// 读取原始内容
	contentBytes, err := os.ReadFile(zshrcPath)
	if err != nil {
		// 如果文件不存在，则创建
		if os.IsNotExist(err) {
			contentBytes = []byte{}
		} else {
			return err
		}
	}
	content := string(contentBytes)

	// 定义 GOROOT 相关的环境变量
	newLines := fmt.Sprintf(
		"export GOROOT=%s\nexport PATH=$GOROOT/bin:$PATH\nexport GOBIN=$GOPATH/bin\n",
		newGOROOT,
	)

	// 先检查是否已有 GOROOT 变量
	re := regexp.MustCompile(`(?m)^export GOROOT=.*$`)
	if re.MatchString(content) {
		// 替换已存在的 GOROOT
		content = re.ReplaceAllString(content, fmt.Sprintf("export GOROOT=%s", newGOROOT))
	} else {
		// 追加新配置到文件末尾
		if !strings.HasSuffix(content, "\n") {
			content += "\n"
		}
		content += newLines
	}

	// 写回文件
	if err := os.WriteFile(zshrcPath, []byte(content), 0644); err != nil {
		return err
	}
	defer func() {
		fmt.Println("请执行以下命令，使 GOROOT 立即生效：")
		fmt.Println("source ~/.zshrc")
	}()

	// 执行 source ~/.zshrc 来使配置生效
	return exec.Command("zsh", "-c", fmt.Sprintf("source ~/%s", filename)).Run()
}
