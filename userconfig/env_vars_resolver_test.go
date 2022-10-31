package userconfig_test

import (
	"errors"
	"testing"

	"github.com/eleven-sh/hetzner-cloud-provider/mocks"
	"github.com/eleven-sh/hetzner-cloud-provider/userconfig"
	"github.com/golang/mock/gomock"
)

func TestEnvVarsResolving(t *testing.T) {
	testCases := []struct {
		test           string
		apiTokenEnvVar string
		regionEnvVar   string
		regionOpts     string
		expectedError  error
		expectedConfig *userconfig.Config
	}{
		{
			test:           "valid",
			apiTokenEnvVar: "a",
			regionEnvVar:   "b",
			expectedConfig: userconfig.NewConfig("a", "b"),
			expectedError:  nil,
		},

		{
			test:           "valid with region opts",
			apiTokenEnvVar: "a",
			regionEnvVar:   "b",
			regionOpts:     "c",
			expectedConfig: userconfig.NewConfig("a", "c"),
			expectedError:  nil,
		},

		{
			test:           "missing region with region opts",
			apiTokenEnvVar: "a",
			regionOpts:     "b",
			expectedConfig: userconfig.NewConfig("a", "b"),
			expectedError:  nil,
		},

		{
			test:           "missing region",
			apiTokenEnvVar: "a",
			expectedError:  userconfig.ErrMissingRegionInEnv,
			expectedConfig: nil,
		},

		{
			test:           "no env vars",
			expectedError:  userconfig.ErrMissingConfig,
			expectedConfig: nil,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.test, func(t *testing.T) {
			mockCtrl := gomock.NewController(t)
			defer mockCtrl.Finish()

			envVarsGetterMock := mocks.NewUserConfigEnvVarsGetter(mockCtrl)
			envVarsGetterMock.EXPECT().Get(userconfig.HetznerAPITokenEnvVar).Return(tc.apiTokenEnvVar).AnyTimes()
			envVarsGetterMock.EXPECT().Get(userconfig.HetznerRegionEnvVar).Return(tc.regionEnvVar).AnyTimes()

			resolver := userconfig.NewEnvVarsResolver(
				envVarsGetterMock,
				userconfig.EnvVarsResolverOpts{
					Region: tc.regionOpts,
				},
			)

			resolvedConfig, err := resolver.Resolve()

			if tc.expectedError == nil && err != nil {
				t.Fatalf("expected no error, got '%+v'", err)
			}

			if tc.expectedError != nil && !errors.Is(err, tc.expectedError) {
				t.Fatalf("expected error to equal '%+v', got '%+v'", tc.expectedError, err)
			}

			if tc.expectedConfig != nil && *resolvedConfig != *tc.expectedConfig {
				t.Fatalf("expected config to equal '%+v', got '%+v'", *tc.expectedConfig, *resolvedConfig)
			}

			if tc.expectedConfig == nil && resolvedConfig != nil {
				t.Fatalf("expected no config, got '%+v'", *resolvedConfig)
			}
		})
	}
}
