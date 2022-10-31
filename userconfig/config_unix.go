//go:build !windows
// +build !windows

package userconfig

import (
	"os/user"
	"path/filepath"
)

func init() {
	usr, err := user.Current()

	if err != nil {
		return
	}

	if usr.HomeDir != "" {
		DefaultConfigFilePath = filepath.Join(usr.HomeDir, ".config", "hcloud", "cli.toml")
	}
}
