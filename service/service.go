package service

import (
	"github.com/eleven-sh/hetzner-cloud-provider/userconfig"
)

type Hetzner struct {
	config *userconfig.Config
}

func NewHetzner(config *userconfig.Config) *Hetzner {
	return &Hetzner{
		config,
	}
}
