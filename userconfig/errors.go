package userconfig

import "errors"

var (
	// ErrMissingConfig represents the error
	// returned when a resolver cannot resolve config.
	ErrMissingConfig = errors.New("ErrMissingConfig")

	// ErrMissingRegion represents the error
	// returned when a resolver cannot resolve region.
	ErrMissingRegion = errors.New("ErrMissingRegion")
)
