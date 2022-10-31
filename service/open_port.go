package service

import (
	"encoding/json"

	"github.com/eleven-sh/eleven/entities"
	"github.com/eleven-sh/eleven/stepper"
	"github.com/eleven-sh/hetzner-cloud-provider/infrastructure"
	"github.com/hetznercloud/hcloud-go/hcloud"
)

func (h *Hetzner) OpenPort(
	stepper stepper.Stepper,
	config *entities.Config,
	cluster *entities.Cluster,
	env *entities.Env,
	portToOpen string,
) error {

	var envInfra *EnvInfrastructure
	err := json.Unmarshal([]byte(env.InfrastructureJSON), &envInfra)

	if err != nil {
		return err
	}

	hetznerClient := hcloud.NewClient(hcloud.WithToken(h.config.Credentials.APIToken))

	return infrastructure.OpenServerPort(
		hetznerClient,
		envInfra.Firewall.ID,
		portToOpen,
	)
}
