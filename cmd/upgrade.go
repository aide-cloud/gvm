/*
Copyright © 2025 NAME HERE <EMAIL ADDRESS>
*/
package cmd

import (
	"fmt"
	"os/exec"

	"github.com/spf13/cobra"
)

// upgradeCmd represents the upgrade command
var upgradeCmd = &cobra.Command{
	Use:   "upgrade",
	Short: "A brief description of your command",
	Long: `A longer description that spans multiple lines and likely contains examples
and usage of using your command. For example:

Cobra is a CLI library for Go that empowers applications.
This application is a tool to generate the needed files
to quickly create a Cobra application.`,
	Run: func(cmd *cobra.Command, args []string) {
		upgrade()
	},
}

func init() {
	rootCmd.AddCommand(upgradeCmd)
}

func upgrade() {
	// go install github.com/aide-cloud/gvm@latest
	if err := exec.Command("/bin/bash", "go", "install", "github.com/aide-cloud/gvm@latest").Run(); err != nil {
		fmt.Println(err)
		return
	}
	// 把当前go版本下的bin目录gvm移动到/usr/local/bin目录下
	if err := exec.Command("sudo", "mv", "$GOROOT/bin/gvm", "/usr/local/bin/").Run(); err != nil {
		fmt.Println(err)
	}
}
