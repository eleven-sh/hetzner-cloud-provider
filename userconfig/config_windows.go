//go:build windows
// +build windows

package userconfig

import (
	"os"
	"path/filepath"
)

func init() {
	dir := os.Getenv("APPDATA")

	if dir != "" {
		DefaultConfigFilePath = filepath.Join(dir, "hcloud", "cli.toml")
	}
}
