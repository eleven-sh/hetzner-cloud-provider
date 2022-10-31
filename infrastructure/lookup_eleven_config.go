package infrastructure

import (
	"errors"
	"io/fs"
	"os"
)

var (
	ErrElevenConfigNotFound = errors.New("ErrElevenConfigNotFound")
)

func LookupElevenConfig(
	elevenConfigDir string,
	apiToken string,
	region string,
) (string, error) {

	configFilePath, err := getElevenConfigFilePath(
		elevenConfigDir,
		apiToken,
		region,
	)

	if err != nil {
		return "", err
	}

	configFileContent, err := os.ReadFile(configFilePath)

	if err != nil && errors.Is(err, fs.ErrNotExist) {
		return "", ErrElevenConfigNotFound
	}

	if err != nil {
		return "", err
	}

	return string(configFileContent), nil
}
