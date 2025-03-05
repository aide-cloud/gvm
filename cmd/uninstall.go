/*
Copyright © 2025 NAME HERE <EMAIL ADDRESS>
*/
package cmd

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/spf13/cobra"
)

// uninstallCmd represents the uninstall command
var uninstallCmd = &cobra.Command{
	Use:   "uninstall",
	Short: "A brief description of your command",
	Long: `A longer description that spans multiple lines and likely contains examples
and usage of using your command. For example:

Cobra is a CLI library for Go that empowers applications.
This application is a tool to generate the needed files
to quickly create a Cobra application.`,
	Run: func(cmd *cobra.Command, args []string) {
		version := args[0]
		removeVersion(version)
	},
	Args: cobra.ExactArgs(1),
}

func init() {
	rootCmd.AddCommand(uninstallCmd)
}

func removeVersion(uninstallVersionFlag string) {
	version := buildGoVersion(uninstallVersionFlag)
	version = "go" + version
	// 检查版本目录是否存在
	if _, err := os.Stat(filepath.Join(rootPath, "versions", version)); os.IsNotExist(err) {
		fmt.Println("版本不存在:", uninstallVersionFlag)
		return
	}

	versionDir := filepath.Join(rootPath, "versions")
	// 删除此目录下的对应版本文件
	if err := os.RemoveAll(filepath.Join(versionDir, version)); err != nil {
		fmt.Println("删除版本失败:", err)
		return
	}
	fmt.Println("删除版本成功:", version)
}
