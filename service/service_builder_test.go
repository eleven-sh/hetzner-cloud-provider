package service_test

import (
	"errors"
	"testing"

	"github.com/eleven-sh/hetzner-cloud-provider/config"
	"github.com/eleven-sh/hetzner-cloud-provider/mocks"
	"github.com/eleven-sh/hetzner-cloud-provider/service"
	"github.com/eleven-sh/hetzner-cloud-provider/userconfig"
	"github.com/golang/mock/gomock"
)

func TestServiceBuilderBuildWithResolvedUserConfig(t *testing.T) {
	mockCtrl := gomock.NewController(t)
	defer mockCtrl.Finish()

	resolvedUserConfig := userconfig.NewConfig("a", "b")

	userConfigResolver := mocks.NewUserConfigResolver(mockCtrl)
	userConfigResolver.EXPECT().Resolve().Return(resolvedUserConfig, nil).Times(1)

	userConfigValidator := mocks.NewUserConfigValidator(mockCtrl)
	userConfigValidator.EXPECT().Validate(resolvedUserConfig).Return(nil).Times(1)

	builder := service.NewBuilder(
		service.BuilderOpts{
			ElevenConfigDir: "",
		},
		userConfigResolver,
		userConfigValidator,
	)
	_, err := builder.Build()

	if err != nil {
		t.Fatalf("expected no error, got '%+v'", err)
	}
}

func TestServiceBuilderBuildWithUserConfigResolverError(t *testing.T) {
	mockCtrl := gomock.NewController(t)
	defer mockCtrl.Finish()

	resolvedUserConfig := userconfig.NewConfig("a", "b")

	userConfigResolverErr := userconfig.ErrMissingConfig
	userConfigResolver := mocks.NewUserConfigResolver(mockCtrl)
	userConfigResolver.EXPECT().Resolve().Return(nil, userConfigResolverErr).Times(1)

	userConfigValidator := mocks.NewUserConfigValidator(mockCtrl)
	userConfigValidator.EXPECT().Validate(resolvedUserConfig).Return(nil).Times(0)

	builder := service.NewBuilder(
		service.BuilderOpts{
			ElevenConfigDir: "",
		},
		userConfigResolver,
		userConfigValidator,
	)
	_, err := builder.Build()

	if err == nil {
		t.Fatalf("expected error, got nothing")
	}

	if !errors.Is(err, userConfigResolverErr) {
		t.Fatalf(
			"expected error to equal '%+v', got '%+v'",
			userConfigResolverErr,
			err,
		)
	}
}

func TestServiceBuilderBuildWithUserConfigValidatorError(t *testing.T) {
	mockCtrl := gomock.NewController(t)
	defer mockCtrl.Finish()

	resolvedUserConfig := userconfig.NewConfig("a", "b")

	userConfigResolver := mocks.NewUserConfigResolver(mockCtrl)
	userConfigResolver.EXPECT().Resolve().Return(resolvedUserConfig, nil).Times(1)

	userConfigValidatorErr := config.ErrInvalidRegion{}
	userConfigValidator := mocks.NewUserConfigValidator(mockCtrl)
	userConfigValidator.EXPECT().Validate(resolvedUserConfig).Return(userConfigValidatorErr).Times(1)

	builder := service.NewBuilder(
		service.BuilderOpts{
			ElevenConfigDir: "",
		},
		userConfigResolver,
		userConfigValidator,
	)
	_, err := builder.Build()

	if err == nil {
		t.Fatalf("expected error, got nothing")
	}

	if !errors.Is(err, userConfigValidatorErr) {
		t.Fatalf(
			"expected error to equal '%+v', got '%+v'",
			userConfigValidatorErr,
			err,
		)
	}
}
