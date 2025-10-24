/*
Copyright © 2025 aide

Permission is hereby granted, free of charge, to any person obtaining a copy
of this software and associated documentation files (the "Software"), to deal
in the Software without restriction, including without limitation the rights
to use, copy, modify, merge, publish, distribute, sublicense, and/or sell
copies of the Software, and to permit persons to whom the Software is
furnished to do so, subject to the following conditions:

The above copyright notice and this permission notice shall be included in
all copies or substantial portions of the Software.

THE SOFTWARE IS PROVIDED "AS IS", WITHOUT WARRANTY OF ANY KIND, EXPRESS OR
IMPLIED, INCLUDING BUT NOT LIMITED TO THE WARRANTIES OF MERCHANTABILITY,
FITNESS FOR A PARTICULAR PURPOSE AND NONINFRINGEMENT. IN NO EVENT SHALL THE
AUTHORS OR COPYRIGHT HOLDERS BE LIABLE FOR ANY CLAIM, DAMAGES OR OTHER
LIABILITY, WHETHER IN AN ACTION OF CONTRACT, TORT OR OTHERWISE, ARISING FROM,
OUT OF OR IN CONNECTION WITH THE SOFTWARE OR THE USE OR OTHER DEALINGS IN
THE SOFTWARE.
*/
package main

import (
	"os"

	"github.com/spf13/cobra"

	"github.com/aide-cloud/gvm/cmd"
	"github.com/aide-cloud/gvm/cmd/install"
	"github.com/aide-cloud/gvm/cmd/list"
	"github.com/aide-cloud/gvm/cmd/ls"
	"github.com/aide-cloud/gvm/cmd/uninstall"
	"github.com/aide-cloud/gvm/cmd/use"
	"github.com/aide-cloud/gvm/pkg/log"
)

func main() {
	rootCmd := cmd.NewCmd()

	commands := []*cobra.Command{
		list.NewListCmd(),
		ls.NewLsCmd(),
		install.NewInstallCmd(),
		uninstall.NewUninstallCmd(),
		use.NewUseCmd(),
	}

	rootCmd.AddCommand(commands...)

	if err := rootCmd.Execute(); err != nil {
		log.Error("Failed to execute root command", "error", err)
		os.Exit(1)
	}
}
