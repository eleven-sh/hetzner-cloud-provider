package infrastructure

import (
	"os"
	"path/filepath"
)

func RemoveElevenConfig(
	elevenConfigDir string,
	apiToken string,
	region string,
) error {

	configFilePath, err := getElevenConfigFilePath(
		elevenConfigDir,
		apiToken,
		region,
	)

	if err != nil {
		return err
	}

	err = os.Remove(configFilePath)

	if err != nil {
		return err
	}

	configDirName := filepath.Dir(
		configFilePath,
	)

	return os.Remove(configDirName)
}
