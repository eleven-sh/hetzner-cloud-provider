package service

import (
	"encoding/json"
	"errors"

	"github.com/eleven-sh/eleven/entities"
	"github.com/eleven-sh/eleven/stepper"
	"github.com/eleven-sh/hetzner-cloud-provider/infrastructure"
)

func (h *Hetzner) CreateElevenConfigStorage(
	stepper stepper.Stepper,
) error {

	stepper.StartTemporaryStep("Creating Eleven's config storage")

	return infrastructure.CreateElevenConfigStorage(
		h.config.ElevenConfigDir,
		h.config.Credentials.APIToken,
		h.config.Region,
	)
}

func (h *Hetzner) LookupElevenConfig(
	stepper stepper.Stepper,
) (*entities.Config, error) {

	configJSON, err := infrastructure.LookupElevenConfig(
		h.config.ElevenConfigDir,
		h.config.Credentials.APIToken,
		h.config.Region,
	)

	if err != nil {

		if errors.Is(err, infrastructure.ErrElevenConfigNotFound) {
			return nil, entities.ErrElevenNotInstalled
		}

		return nil, err
	}

	var elevenConfig *entities.Config
	err = json.Unmarshal([]byte(configJSON), &elevenConfig)

	if err != nil {
		return nil, err
	}

	return elevenConfig, nil
}

func (h *Hetzner) SaveElevenConfig(
	stepper stepper.Stepper,
	config *entities.Config,
) error {

	configJSON, err := json.Marshal(config)

	if err != nil {
		return err
	}

	return infrastructure.SaveElevenConfig(
		h.config.ElevenConfigDir,
		h.config.Credentials.APIToken,
		h.config.Region,
		configJSON,
	)
}

func (h *Hetzner) RemoveElevenConfigStorage(
	stepper stepper.Stepper,
) error {

	stepper.StartTemporaryStep("Removing Eleven's config storage")

	return infrastructure.RemoveElevenConfig(
		h.config.ElevenConfigDir,
		h.config.Credentials.APIToken,
		h.config.Region,
	)
}
