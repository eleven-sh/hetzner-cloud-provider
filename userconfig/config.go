package userconfig

var DefaultConfigFilePath string

// Credentials represents the Hetzner credentials resolved from user config.
type Credentials struct {
	// APIToken represents the API token.
	APIToken string
}

// HasKeys is an helper method used to check
// that the credentials in the Credentials struct are not empty.
func (c Credentials) HasKeys() bool {
	return len(c.APIToken) > 0
}

// Config represents the resolved user config.
type Config struct {
	// Credentials represents the resolved credentials (API token).
	Credentials Credentials

	// ElevenConfigDir represents the path to the directory where the Eleven's configuration is stored.
	ElevenConfigDir string

	// Region represents the resolved region.
	Region string
}

// NewConfig constructs a new resolved user config.
func NewConfig(
	apiToken string,
	region string,
) *Config {

	return &Config{
		Credentials: Credentials{
			APIToken: apiToken,
		},
		Region: region,
	}
}
