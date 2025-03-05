/*
Copyright © 2025 NAME HERE <EMAIL ADDRESS>
*/
package cmd

import (
	"fmt"
	"sort"
	"strconv"
	"strings"

	"github.com/aide-cloud/gvm/internal"
	"github.com/spf13/cobra"
)

// listCmd represents the list command
var listCmd = &cobra.Command{
	Use:   "list",
	Short: "A brief description of your command",
	Long: `A longer description that spans multiple lines and likely contains examples
and usage of using your command. For example:

Cobra is a CLI library for Go that empowers applications.
This application is a tool to generate the needed files
to quickly create a Cobra application.`,
	Run: func(cmd *cobra.Command, args []string) {
		if originFlag {
			OriginVersions()
			return
		}
		LocalVersions()
	},
}

var originFlag bool
var versionAll bool
var versionFlag string

func init() {
	rootCmd.AddCommand(listCmd)
	listCmd.Flags().BoolVarP(&originFlag, "origin", "o", false, "List available Go versions from official source")
	listCmd.Flags().BoolVarP(&versionAll, "all", "a", false, "List all available Go versions from official source")
	listCmd.Flags().StringVarP(&versionFlag, "version", "v", "", "List available Go versions for a specific version")
}

func OriginVersions() {
	fmt.Println("Go origin versions: ")
	versions := internal.FetchGoVersions()
	if len(versions) == 0 {
		return
	}
	list := make([]string, 0, len(versions))
	for version := range versions {
		list = append(list, version)
	}
	showVersions(list)
}

func LocalVersions() {
	fmt.Println("Go installed versions: ")
	versions := internal.FetchLocalVersions()
	if len(versions) == 0 {
		return
	}
	showVersions(versions)
}

func getLastVersion() string {
	versions := internal.FetchGoVersions()
	list := make([]string, 0, len(versions))
	for version := range versions {
		list = append(list, version)
	}
	// 排序
	sortVersions(list)
	for _, s := range list {
		if len(strings.Split(s, ".")) == 3 {
			return s
		}
	}
	return ""
}

func showVersions(list []string, currentVersion ...string) {
	// 排序
	sortVersions(list)

	versions := make([]string, 0, len(list))
	// 如果versionFlag不为空，则只列出versionFlag对应的版本号
	if versionFlag != "" {
		versionFlag = strings.TrimPrefix(versionFlag, "go")
		versionFlag = strings.TrimPrefix(versionFlag, "v")
		versionFlag = strings.TrimPrefix(versionFlag, "V")
		for _, s := range list {
			if strings.HasPrefix(s, versionFlag) {
				versions = append(versions, s)
			}
		}
		return
	}
	if !versionAll && len(list) > 10 {
		list = list[:10]
	}
	for _, s := range list {
		if len(currentVersion) > 0 && s == currentVersion[0] {
			versions = append(versions, "* "+s)
		} else {
			versions = append(versions, "  "+s)
		}
	}
	fmt.Println(strings.Join(versions, "\n"))
}

// 对 Go 版本号进行排序（从大到小）
func sortVersions(versions []string) {
	sort.Slice(versions, func(i, j int) bool {
		return compareVersions(versions[i], versions[j])
	})
}

// 比较两个 Go 版本号，决定是否需要交换位置
func compareVersions(v1, v2 string) bool {
	parts1 := strings.Split(v1, ".")
	parts2 := strings.Split(v2, ".")
	if len(parts1) < 3 {
		parts1 = append(parts1, "99999999999")
	}
	if len(parts2) < 3 {
		parts2 = append(parts2, "99999999999")
	}

	for i := 0; i < 3; i++ {
		num1, _ := strconv.Atoi(parts1[i])
		num2, _ := strconv.Atoi(parts2[i])
		if num1 != num2 {
			return num1 > num2 // 从大到小排序
		}
	}
	return false
}
