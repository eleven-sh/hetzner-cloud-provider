package userconfig

import (
	"errors"
)

var (
	// ErrMissingRegionInEnv represents the error
	// returned when API token is set but a region was not
	// passed as an option nor set as a environment variable.
	ErrMissingRegionInEnv = errors.New("ErrMissingRegionInEnv")
)

const (
	// HetznerAPITokenEnvVar represents the environment variable name
	// that the resolver will look for when resolving the Hetzner API token.
	HetznerAPITokenEnvVar = "HCLOUD_TOKEN"

	// HetznerRegionEnvVar represents the environment variable name
	// that the resolver will look for when resolving the Hetzner region.
	HetznerRegionEnvVar = "HCLOUD_REGION"
)

//go:generate go run github.com/golang/mock/mockgen -destination ../mocks/user_config_env_vars_getter.go -package mocks -mock_names EnvVarsGetter=UserConfigEnvVarsGetter  github.com/eleven-sh/hetzner-cloud-provider/userconfig EnvVarsGetter
// EnvVarsGetter represents the interface
// used to access environment variables.
type EnvVarsGetter interface {
	Get(string) string
}

// EnvVarsResolverOpts represents the options
// used to configure the EnvVarsResolver.
type EnvVarsResolverOpts struct {
	// Region specifies which region will be used in the resulting config.
	// Default to the one found in sandbox if not set.
	Region string
}

// EnvVarsResolver retrieves the Hetzner account
// configuration from environment variables.
type EnvVarsResolver struct {
	opts    EnvVarsResolverOpts
	envVars EnvVarsGetter
}

// NewFilesResolver constructs the EnvVarsResolver struct.
func NewEnvVarsResolver(
	envVars EnvVarsGetter,
	opts EnvVarsResolverOpts,
) EnvVarsResolver {

	return EnvVarsResolver{
		opts:    opts,
		envVars: envVars,
	}
}

// Resolve retrieves the Hetzner account configuration
// from environment variables.
//
// The Region option takes precedence over the one
// found in sandbox.
//
// Partial configurations return an adequate errror.
//
// Env vars are retrieved via the EnvVarsGetter interface
// passed as constructor argument.
func (e EnvVarsResolver) Resolve() (*Config, error) {
	regionInEnv := e.envVars.Get(HetznerRegionEnvVar)

	resolvedConfig := NewConfig(
		e.envVars.Get(HetznerAPITokenEnvVar),
		e.resolveRegion(regionInEnv),
	)

	if resolvedConfig.Credentials.HasKeys() &&
		len(resolvedConfig.Region) > 0 {

		return resolvedConfig, nil
	}

	if resolvedConfig.Credentials.HasKeys() &&
		len(resolvedConfig.Region) == 0 {

		return nil, ErrMissingRegionInEnv
	}

	return nil, ErrMissingConfig
}

func (e EnvVarsResolver) resolveRegion(regionInEnv string) string {
	if len(e.opts.Region) > 0 {
		return e.opts.Region
	}

	return regionInEnv
}
