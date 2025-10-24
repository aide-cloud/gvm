package cmd

import (
	"github.com/spf13/cobra"

	"github.com/aide-cloud/gvm/internal/version"
	"github.com/aide-cloud/gvm/pkg/log"
)

var globalFlags = GlobalFlags{}

type GlobalFlags struct {
	OriginURL        string
	DownloadURL      string
	CacheDir         string
	SdkDir           string
	VersionFilePath  string
	LocalVersionFile string

	Eval bool
}

func InitFlags(cmd *cobra.Command) {
	cmd.Flags().StringVar(&globalFlags.OriginURL, "origin-url", "https://go.dev/dl/?mode=json&include=all", "The URL to fetch the origin versions")
	cmd.Flags().StringVar(&globalFlags.DownloadURL, "download-url", "https://dl.google.com/go/", "The URL to download the sdk")
	cmd.Flags().StringVar(&globalFlags.CacheDir, "cache-dir", version.DefaultCacheDir, "The directory to cache the origin versions")
	cmd.Flags().StringVar(&globalFlags.SdkDir, "sdk-dir", version.DefaultSdkDir, "The directory to store the sdk")
	cmd.Flags().StringVar(&globalFlags.VersionFilePath, "version-file-path", version.DefaultVersionFilePath, "The file path to store the versions")
	cmd.Flags().StringVar(&globalFlags.LocalVersionFile, "local-version-file", version.DefaultLocalVersionFilePath, "The file path to store the local versions")
	cmd.Flags().BoolVar(&globalFlags.Eval, "eval", false, "Eval the command")
}

func GetGlobalFlags() GlobalFlags {
	log.SetPrintEnable(!globalFlags.Eval)
	return globalFlags
}
