package config_test

import (
	"errors"
	"testing"

	"github.com/eleven-sh/hetzner-cloud-provider/config"
	"github.com/eleven-sh/hetzner-cloud-provider/userconfig"
)

func TestLoadWithExistingContexts(t *testing.T) {
	testCases := []struct {
		test           string
		context        string
		expectedConfig *userconfig.Config
	}{
		{
			test:    "with default context",
			context: "",
			expectedConfig: userconfig.NewConfig(
				"abcdef",
				"",
			),
		},

		{
			test:    "with non-default context",
			context: "production",
			expectedConfig: userconfig.NewConfig(
				"123456",
				"",
			),
		},
	}

	for _, tc := range testCases {
		t.Run(tc.test, func(t *testing.T) {
			contextLoader := config.NewContextLoader()
			loadedContext, err := contextLoader.Load(
				tc.context,
				"./testdata/user_config",
			)

			if err != nil {
				t.Fatalf("expected no error, got '%+v'", err)
			}

			if loadedContext.Credentials.APIToken != tc.expectedConfig.Credentials.APIToken ||
				loadedContext.Region != tc.expectedConfig.Region {

				t.Fatalf("expected config to equal '%+v', got '%+v'", *tc.expectedConfig, loadedContext)
			}
		})
	}
}

func TestLoadWithInvalidContext(t *testing.T) {
	contextLoader := config.NewContextLoader()
	_, err := contextLoader.Load(
		"non_existing_context",
		"./testdata/user_config",
	)

	if !errors.As(err, &config.ErrContextNotFound{}) {
		t.Fatalf(
			"expected error to equal '%+v', got '%+v'",
			&config.ErrContextNotFound{},
			err,
		)
	}
}

func TestLoadWithNonExistingConfigFileAndContext(t *testing.T) {
	contextLoader := config.NewContextLoader()
	_, err := contextLoader.Load(
		"context",
		"./testdata/non_existing_user_config",
	)

	if !errors.As(err, &config.ErrContextNotFound{}) {
		t.Fatalf(
			"expected error to equal '%+v', got '%+v'",
			&config.ErrContextNotFound{},
			err,
		)
	}
}

func TestLoadWithNonExistingConfigFileAndNoContext(t *testing.T) {
	contextLoader := config.NewContextLoader()
	userConfig, err := contextLoader.Load(
		"",
		"./testdata/non_existing_user_config",
	)

	if err != nil {
		t.Fatalf("expected no error, got '%+v'", err)
	}

	if userConfig.Credentials.APIToken != "" {
		t.Fatalf("expected no API token in config")
	}

	if userConfig.Region != "" {
		t.Fatalf("expected no region in config")
	}
}

func TestLoadWithInvalidConfigFile(t *testing.T) {
	contextLoader := config.NewContextLoader()
	_, err := contextLoader.Load(
		"",
		"./testdata/invalid_user_config",
	)

	if err == nil {
		t.Fatalf("expected error, got nothing")
	}
}
