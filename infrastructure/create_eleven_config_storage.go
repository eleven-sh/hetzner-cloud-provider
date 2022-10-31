package infrastructure

import (
	"crypto/sha1"
	"encoding/hex"
	"io/fs"
	"os"
	"path"
	"path/filepath"
)

const ElevenConfigFileName = "config.json"

func CreateElevenConfigStorage(
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

	return os.MkdirAll(
		filepath.Dir(configFilePath),
		fs.FileMode(0700),
	)
}

func getElevenConfigFilePath(
	elevenConfigDir string,
	apiToken string,
	region string,
) (string, error) {

	storageDirName, err := computeElevenConfigStorageDirName(apiToken, region)

	if err != nil {
		return "", err
	}

	return path.Join(
		elevenConfigDir,
		storageDirName,
		ElevenConfigFileName,
	), nil
}

func computeElevenConfigStorageDirName(
	apiToken string,
	region string,
) (string, error) {

	sha1Computer := sha1.New()
	_, err := sha1Computer.Write([]byte("hetzner." + apiToken + "." + region))

	if err != nil {
		return "", err
	}

	return hex.EncodeToString(sha1Computer.Sum(nil)), nil
}
