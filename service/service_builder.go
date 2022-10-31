package service

import (
	"github.com/eleven-sh/eleven/entities"
	"github.com/eleven-sh/hetzner-cloud-provider/userconfig"
)

//go:generate go run github.com/golang/mock/mockgen -destination ../mocks/user_config_resolver.go -package mocks -mock_names UserConfigResolver=UserConfigResolver github.com/eleven-sh/hetzner-cloud-provider/service UserConfigResolver
type UserConfigResolver interface {
	Resolve() (*userconfig.Config, error)
}

//go:generate go run github.com/golang/mock/mockgen -destination ../mocks/user_config_validator.go -package mocks -mock_names UserConfigValidator=UserConfigValidator github.com/eleven-sh/hetzner-cloud-provider/service UserConfigValidator
type UserConfigValidator interface {
	Validate(userConfig *userconfig.Config) error
}

type BuilderOpts struct {
	ElevenConfigDir string
}

type Builder struct {
	opts                BuilderOpts
	userConfigResolver  UserConfigResolver
	userConfigValidator UserConfigValidator
}

func NewBuilder(
	opts BuilderOpts,
	userConfigResolver UserConfigResolver,
	userConfigValidator UserConfigValidator,
) Builder {

	return Builder{
		opts:                opts,
		userConfigResolver:  userConfigResolver,
		userConfigValidator: userConfigValidator,
	}
}

func (b Builder) Build() (entities.CloudService, error) {
	userConfig, err := b.userConfigResolver.Resolve()

	if err != nil {
		return nil, err
	}

	if err := b.userConfigValidator.Validate(userConfig); err != nil {
		return nil, err
	}

	userConfig.ElevenConfigDir = b.opts.ElevenConfigDir

	return NewHetzner(userConfig), nil
}
