package version

import (
	"fmt"
	"os"
	"path/filepath"
	"regexp"

	"github.com/aide-cloud/gvm/pkg/dir"
	"github.com/aide-cloud/gvm/pkg/log"
)

func Use(version, sdkDir string, isEval bool) error {
	sdkDir = dir.ExpandHomeDir(sdkDir)
	sdkFilePath := filepath.Join(sdkDir, version)
	exist, err := dir.CheckFileExists(sdkFilePath)
	if err != nil {
		return err
	}
	if !exist {
		return fmt.Errorf("version %s not found", version)
	}
	// ~/.zshrc set the go root
	shell := os.Getenv("SHELL")
	var (
		shellConfig     []byte
		shellConfigPath string
	)
	switch shell {
	case "/bin/zsh", "/bin/sh":
		shellConfigPath = filepath.Join(os.Getenv("HOME"), ".zshrc")
		shellConfig, err = os.ReadFile(shellConfigPath)
		if err != nil {
			return err
		}
	case "/bin/bash":
		shellConfigPath = filepath.Join(os.Getenv("HOME"), ".bashrc")
		shellConfig, err = os.ReadFile(shellConfigPath)
		if err != nil {
			return err
		}
	default:
		return fmt.Errorf("unsupported shell: %s", shell)
	}

	exportGOROOT := fmt.Sprintf(`export GOROOT="%s"`, sdkFilePath)

	gorootRegex := regexp.MustCompile(`(?m)^export GOROOT=.*$`)
	if gorootRegex.Match(shellConfig) {
		shellConfig = gorootRegex.ReplaceAll(shellConfig, []byte(exportGOROOT))
	} else {
		shellConfig = append(shellConfig, []byte(exportGOROOT)...)
	}
	if err := os.WriteFile(shellConfigPath, []byte(shellConfig), 0644); err != nil {
		return err
	}
	log.Info("set", "GOROOT", sdkFilePath)
	if isEval {
		fmt.Printf("source %s\n", shellConfigPath)
	} else {
		fmt.Printf("execute command:\n\t eval \"source %s\"\n", shellConfigPath)
	}
	return nil
}
