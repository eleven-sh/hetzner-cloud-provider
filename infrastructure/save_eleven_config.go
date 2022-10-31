package infrastructure

import (
	"os"
)

func SaveElevenConfig(
	elevenConfigDir string,
	apiToken string,
	region string,
	configJSON []byte,
) error {

	configFilePath, err := getElevenConfigFilePath(
		elevenConfigDir,
		apiToken,
		region,
	)

	if err != nil {
		return err
	}

	return os.WriteFile(
		configFilePath,
		configJSON,
		os.FileMode(0600),
	)
}
