package config_test

import (
	"errors"
	"strings"
	"testing"

	"github.com/eleven-sh/hetzner-cloud-provider/config"
	"github.com/eleven-sh/hetzner-cloud-provider/userconfig"
)

func TestUserConfigValidator(t *testing.T) {
	testCases := []struct {
		test          string
		userconfig    *userconfig.Config
		expectedError error
	}{
		{
			test: "with valid config",
			userconfig: &userconfig.Config{
				Credentials: userconfig.Credentials{
					APIToken: strings.Repeat("b", 64),
				},
				Region: "fsn1",
			},
			expectedError: nil,
		},

		{
			test: "with invalid region",
			userconfig: &userconfig.Config{
				Credentials: userconfig.Credentials{
					APIToken: strings.Repeat("b", 64),
				},
				Region: "invalid_region",
			},
			expectedError: config.ErrInvalidRegion{},
		},

		{
			test: "with invalid API token",
			userconfig: &userconfig.Config{
				Credentials: userconfig.Credentials{
					APIToken: "invalid_api_token",
				},
				Region: "nbg1",
			},
			expectedError: config.ErrInvalidAPIToken{},
		},
	}

	for _, tc := range testCases {
		t.Run(tc.test, func(t *testing.T) {
			configvalidator := config.NewUserConfigValidator()
			err := configvalidator.Validate(tc.userconfig)

			if tc.expectedError == nil && err != nil {
				t.Fatalf("expected no error, got '%+v'", err)
			}

			if tc.expectedError == nil {
				return
			}

			if _, ok := tc.expectedError.(config.ErrInvalidRegion); ok {
				if !errors.As(err, &config.ErrInvalidRegion{}) {
					t.Fatalf(
						"expected error to equal '%+v', got '%+v'",
						tc.expectedError,
						err,
					)
				}
			}

			if _, ok := tc.expectedError.(config.ErrInvalidAPIToken); ok {
				if !errors.As(err, &config.ErrInvalidAPIToken{}) {
					t.Fatalf(
						"expected error to equal '%+v', got '%+v'",
						tc.expectedError,
						err,
					)
				}
			}
		})
	}
}
