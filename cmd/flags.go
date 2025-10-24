package cmd

import (
	"github.com/aide-cloud/gvm/pkg/log"
	"github.com/spf13/cobra"
)

var globalFlags = GlobalFlags{}

type GlobalFlags struct {
	OriginURL       string
	DownloadURL     string
	CacheDir        string
	SdkDir          string
	VersionFilePath string

	Eval bool
}

func InitFlags(cmd *cobra.Command) {
	cmd.Flags().StringVar(&globalFlags.OriginURL, "origin-url", "https://go.dev/dl/?mode=json&include=all", "The URL to fetch the origin versions")
	cmd.Flags().StringVar(&globalFlags.DownloadURL, "download-url", "https://dl.google.com/go/", "The URL to download the sdk")
	cmd.Flags().StringVar(&globalFlags.CacheDir, "cache-dir", "~/.gvm/cache", "The directory to cache the origin versions")
	cmd.Flags().StringVar(&globalFlags.SdkDir, "sdk-dir", "~/.gvm/sdk", "The directory to store the sdk")
	cmd.Flags().StringVar(&globalFlags.VersionFilePath, "version-file-path", "~/.gvm/versions.json", "The file path to store the versions")
	cmd.Flags().BoolVar(&globalFlags.Eval, "eval", false, "Eval the command")
}

func GetGlobalFlags() GlobalFlags {
	log.SetPrintEnable(!globalFlags.Eval)
	return globalFlags
}
