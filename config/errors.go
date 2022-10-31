package config

// ErrInvalidRegion represents the error
// returned when the region in user config is invalid.
type ErrInvalidRegion struct {
	Region string
}

func (ErrInvalidRegion) Error() string {
	return "ErrInvalidRegion"
}

// ErrInvalidAPIToken represents the error
// returned when the API token in user config is invalid.
type ErrInvalidAPIToken struct {
	APIToken string
}

func (ErrInvalidAPIToken) Error() string {
	return "ErrInvalidAccessKeyID"
}

// ErrContextNotFound represents the error
// returned when the context passed in option was not found.
type ErrContextNotFound struct {
	Context        string
	ConfigFilePath string
}

func (ErrContextNotFound) Error() string {
	return "ErrContextNotFound"
}
