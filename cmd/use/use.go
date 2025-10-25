package use

import (
	"fmt"

	"github.com/aide-cloud/gvm/cmd"
	"github.com/spf13/cobra"
)

func NewUseCmd() *cobra.Command {
	useCmd := &cobra.Command{
		Use:   "use",
		Short: "Use a specific Go version",
		Long: `Use a specific Go version.
Example:
  gvm use go1.25.3
  gvm use latest
  gvm use -l
`,
		Annotations: map[string]string{
			"group": cmd.VersionCommands,
		},
		Run: func(cmd *cobra.Command, args []string) {
			if len(args) > 0 {
				useFlags.version = args[0]
			}
			if useFlags.latest {
				useFlags.version = "latest"
			}
			if useFlags.version == "" {
				fmt.Println("Please specify the version to use")
				return
			}

			useFlags.use()
		},
	}
	useFlags.initFlags(useCmd)
	return useCmd
}

var useFlags = useCmdFlags{}

type useCmdFlags struct {
	cmd.GlobalFlags
	version string
	latest  bool
	isForce bool
}

func (u *useCmdFlags) initFlags(c *cobra.Command) {
	cmd.InitFlags(c)
	c.Flags().BoolVarP(&u.latest, "latest", "l", false, "Use the latest version")
	c.Flags().BoolVarP(&u.isForce, "force", "f", false, "Force use the version")
}

func (u *useCmdFlags) use() {
	u.GlobalFlags = cmd.GetGlobalFlags()
	v := cmd.NewVersionManager()
	v.Use(u.version, u.isForce, u.Eval)
}
