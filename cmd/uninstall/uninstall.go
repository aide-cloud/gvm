package uninstall

import (
	"fmt"

	"github.com/aide-cloud/gvm/cmd"
	"github.com/spf13/cobra"
)

func NewUninstallCmd() *cobra.Command {
	uninstallCmd := &cobra.Command{
		Use:   "uninstall",
		Short: "Uninstall a specific Go version",
		Long: `Uninstall a specific Go version.
Example:
  gvm uninstall go1.25.3
  gvm uninstall latest
  gvm uninstall -l
`,
		Annotations: map[string]string{
			"group": cmd.VersionCommands,
		},
		Run: func(cmd *cobra.Command, args []string) {
			if len(args) > 0 {
				uninstallFlags.version = args[0]
			}
			if uninstallFlags.latest {
				uninstallFlags.version = "latest"
			}
			if uninstallFlags.version == "" {
				fmt.Println("Please specify the version to uninstall")
				return
			}
			uninstallFlags.uninstall()
		},
	}
	uninstallFlags.initFlags(uninstallCmd)
	return uninstallCmd
}

var uninstallFlags = uninstallCmdFlags{}

type uninstallCmdFlags struct {
	cmd.GlobalFlags
	version string
	latest  bool
}

func (u *uninstallCmdFlags) initFlags(c *cobra.Command) {
	cmd.InitFlags(c)
	c.Flags().BoolVarP(&u.latest, "latest", "l", false, "Uninstall the latest version")
}

func (u *uninstallCmdFlags) uninstall() {
	u.GlobalFlags = cmd.GetGlobalFlags()
	v := cmd.NewVersionManager()
	v.Uninstall(u.version)
}
