package config

import (
	"github.com/eleven-sh/hetzner-cloud-provider/userconfig"
)

var validHetznerRegions = map[string]bool{
	"fsn1": true,
	"nbg1": true,
	"hel1": true,
	"ash":  true,
}

type UserConfigValidator struct{}

func NewUserConfigValidator() UserConfigValidator {
	return UserConfigValidator{}
}

func (u UserConfigValidator) Validate(userConfig *userconfig.Config) error {
	region := userConfig.Region

	if err := u.validateRegion(region); err != nil {
		return err
	}

	creds := userConfig.Credentials
	apiToken := creds.APIToken

	if err := u.validateAPIToken(apiToken); err != nil {
		return err
	}

	return nil
}

func (UserConfigValidator) validateRegion(region string) error {
	if _, ok := validHetznerRegions[region]; !ok {
		return ErrInvalidRegion{
			Region: region,
		}
	}

	return nil
}

func (UserConfigValidator) validateAPIToken(apiToken string) error {
	if len(apiToken) != 64 {
		return ErrInvalidAPIToken{
			APIToken: apiToken,
		}
	}

	return nil
}
