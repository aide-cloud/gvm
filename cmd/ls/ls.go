package ls

import (
	"github.com/spf13/cobra"

	"github.com/aide-cloud/gvm/cmd"
	"github.com/aide-cloud/gvm/internal/version"
)

func NewLsCmd() *cobra.Command {
	lsCmd := &cobra.Command{
		Use:   "ls",
		Short: "List out the versions of Go that have been installed locally.",
		Long: `List out the versions of Go that have been installed locally.
Example:
  gvm ls
`,
		Annotations: map[string]string{
			"group": cmd.VersionCommands,
		},
		Run: func(cmd *cobra.Command, args []string) {
			lsFlags.versions()
		},
	}
	lsFlags.initFlags(lsCmd)
	return lsCmd
}

var lsFlags = lsCmdFlags{}

type lsCmdFlags struct {
	cmd.GlobalFlags
}

func (l *lsCmdFlags) initFlags(c *cobra.Command) {
	cmd.InitFlags(c)
}

func (l *lsCmdFlags) versions() {
	l.GlobalFlags = cmd.GetGlobalFlags()
	v := version.NewVersion(
		version.WithSdkDir(l.SdkDir),
		version.WithCacheDir(l.CacheDir),
		version.WithOriginURL(l.OriginURL),
		version.WithDownloadURL(l.DownloadURL),
	)
	v.Ls()
}
