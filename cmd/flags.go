package cmd

import (
	"github.com/spf13/cobra"

	"github.com/aide-cloud/gvm/pkg/env"
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
	cmd.Flags().StringVar(&globalFlags.OriginURL, "origin-url", env.GetEnv("GVM_ORIGIN_URL", "https://go.dev/dl/?mode=json&include=all"), "The URL to fetch the origin versions, env: GVM_ORIGIN_URL")
	cmd.Flags().StringVar(&globalFlags.DownloadURL, "download-url", env.GetEnv("GVM_DOWNLOAD_URL", "https://dl.google.com/go/"), "The URL to download the sdk, env: GVM_DOWNLOAD_URL")
	cmd.Flags().StringVar(&globalFlags.CacheDir, "cache-dir", env.GetEnv("GVM_CACHE_DIR", "~/.gvm/cache"), "The directory to cache the origin versions, env: GVM_CACHE_DIR")
	cmd.Flags().StringVar(&globalFlags.SdkDir, "sdk-dir", env.GetEnv("GVM_SDK_DIR", "~/go/sdk"), "The directory to store the sdk, env: GVM_SDK_DIR")
	cmd.Flags().StringVar(&globalFlags.VersionFilePath, "version-file-path", env.GetEnv("GVM_VERSION_FILE_PATH", "~/.gvm/versions.json"), "The file path to store the versions, env: GVM_VERSION_FILE_PATH")
	cmd.Flags().StringVar(&globalFlags.LocalVersionFile, "local-version-file", env.GetEnv("GVM_LOCAL_VERSION_FILE_PATH", "~/.gvm/version"), "The file path to store the local versions, env: GVM_LOCAL_VERSION_FILE_PATH")
	cmd.Flags().BoolVar(&globalFlags.Eval, "eval", false, "Eval the command")
}

func GetGlobalFlags() GlobalFlags {
	log.SetPrintEnable(!globalFlags.Eval)
	return globalFlags
}
