package list

import (
	"github.com/aide-cloud/gvm/cmd"
	"github.com/aide-cloud/gvm/internal/version"
	"github.com/spf13/cobra"
)

func NewListCmd() *cobra.Command {
	listCmd := &cobra.Command{
		Use:   "list",
		Short: "List out the available Go versions that can be installed",
		Long: `List out the available Go versions that can be installed.
Example:
  gvm list
  gvm list -n 10
  gvm list -l
`,
		Annotations: map[string]string{
			"group": cmd.VersionCommands,
		},
		Run: func(cmd *cobra.Command, args []string) {
			listFlags.versions()
		},
	}

	listFlags.initFlags(listCmd)

	return listCmd
}

var listFlags = listCmdFlags{}

type listCmdFlags struct {
	cmd.GlobalFlags
	number      int
	forceUpdate bool
	latest      bool
}

func (l *listCmdFlags) initFlags(c *cobra.Command) {
	cmd.InitFlags(c)
	c.Flags().IntVarP(&l.number, "number", "n", 10, "The number of versions to list")
	c.Flags().BoolVar(&l.forceUpdate, "force-update", false, "Force update the origin versions cache")
	c.Flags().BoolVarP(&l.latest, "latest", "l", false, "Show the latest version only")

}

func (l *listCmdFlags) versions() {
	l.GlobalFlags = cmd.GetGlobalFlags()
	v := version.NewVersion(
		version.WithSdkDir(l.SdkDir),
		version.WithCacheDir(l.CacheDir),
		version.WithOriginURL(l.OriginURL),
		version.WithDownloadURL(l.DownloadURL),
	)
	v.List(l.latest, l.number, l.forceUpdate)
}
