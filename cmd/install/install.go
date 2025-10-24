package install

import (
	"fmt"

	"github.com/spf13/cobra"

	"github.com/aide-cloud/gvm/cmd"
	"github.com/aide-cloud/gvm/internal/version"
)

func NewInstallCmd() *cobra.Command {
	installCmd := &cobra.Command{
		Use:   "install",
		Short: "Install a specific Go version",
		Long: `Install a specific Go version.
Example:
  gvm install 1.25.3
  gvm install latest
  gvm install -l
`,
		Annotations: map[string]string{
			"group": cmd.VersionCommands,
		},
		Run: func(cmd *cobra.Command, args []string) {
			if len(args) > 0 {
				installFlags.version = args[0]
			}
			if installFlags.latest {
				installFlags.version = "latest"
			}
			if installFlags.version == "" {
				fmt.Println("Please specify the version to install")
				return
			}
			installFlags.install()
		},
	}
	installFlags.initFlags(installCmd)
	return installCmd
}

var installFlags = installCmdFlags{}

type installCmdFlags struct {
	cmd.GlobalFlags
	version string
	latest  bool
	isForce bool
}

func (i *installCmdFlags) initFlags(c *cobra.Command) {
	cmd.InitFlags(c)
	c.Flags().BoolVarP(&i.latest, "latest", "l", false, "Install the latest version")
	c.Flags().BoolVarP(&i.isForce, "force", "f", false, "Force install the version")
}

func (i *installCmdFlags) install() {
	i.GlobalFlags = cmd.GetGlobalFlags()
	v := version.NewVersion(
		version.WithSdkDir(i.SdkDir),
		version.WithCacheDir(i.CacheDir),
		version.WithOriginURL(i.OriginURL),
		version.WithDownloadURL(i.DownloadURL),
	)
	v.Install(i.version, i.isForce)
}
