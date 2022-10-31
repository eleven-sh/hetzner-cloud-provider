package userconfig

//go:generate go run github.com/golang/mock/mockgen -destination ../mocks/user_config_context_loader.go -package mocks -mock_names ContextLoader=UserConfigContextLoader github.com/eleven-sh/hetzner-cloud-provider/userconfig ContextLoader
// ContextLoader represents the interface
// used to load configuration context from files.
type ContextLoader interface {
	Load(
		context string,
		configFilePath string,
	) (*Config, error)
}

// FilesResolverOpts represents the options
// used to configure the FilesResolver.
type FilesResolverOpts struct {
	// Context specifies which configuration context will be loaded.
	Context string

	// Region specifies which region will be used in the resulting config.
	// Default to the one found in config files if not set.
	Region string
}

// FilesResolver retrieves the Hetzner account
// configuration from config files.
type FilesResolver struct {
	opts          FilesResolverOpts
	contextLoader ContextLoader
	envVars       EnvVarsGetter
}

// NewFilesResolver constructs the FilesResolver struct.
func NewFilesResolver(
	contextLoader ContextLoader,
	opts FilesResolverOpts,
	envVars EnvVarsGetter,
) FilesResolver {

	return FilesResolver{
		contextLoader: contextLoader,
		opts:          opts,
		envVars:       envVars,
	}
}

// Resolve retrieves the Hetzner account configuration from config files.
//
// The ConfigFilePath option is used to locate the config file.
//
// The Context option specifies which configuration context
// will be loaded.
//
// The Region option takes precedence over the region found in config file.
//
// Config file is loaded via the ContextLoader interface
// passed as constructor argument.
func (f FilesResolver) Resolve() (*Config, error) {
	loadedContext, err := f.contextLoader.Load(
		f.resolveContext(),
		DefaultConfigFilePath,
	)

	if err != nil {
		return nil, err
	}

	if !loadedContext.Credentials.HasKeys() {
		return nil, ErrMissingConfig
	}

	resolvedRegion := f.resolveRegion()

	if len(resolvedRegion) == 0 {
		return nil, ErrMissingRegion
	}

	resolvedConfig := NewConfig(
		loadedContext.Credentials.APIToken,
		resolvedRegion,
	)

	return resolvedConfig, nil
}

func (f FilesResolver) resolveContext() string {
	return f.opts.Context
}

func (f FilesResolver) resolveRegion() string {
	if len(f.opts.Region) > 0 {
		return f.opts.Region
	}

	return f.envVars.Get(HetznerRegionEnvVar)
}
