package service

import (
	"encoding/json"

	"github.com/eleven-sh/eleven/entities"
	"github.com/eleven-sh/eleven/queues"
	"github.com/eleven-sh/eleven/stepper"
	"github.com/eleven-sh/hetzner-cloud-provider/infrastructure"
	"github.com/hetznercloud/hcloud-go/hcloud"
)

func (h *Hetzner) RemoveEnv(
	stepper stepper.Stepper,
	config *entities.Config,
	cluster *entities.Cluster,
	env *entities.Env,
) error {

	var envInfra *EnvInfrastructure
	err := json.Unmarshal([]byte(env.InfrastructureJSON), &envInfra)

	if err != nil {
		return err
	}

	hetznerClient := hcloud.NewClient(hcloud.WithToken(h.config.Credentials.APIToken))
	envInfraQueue := queues.InfrastructureQueue[*EnvInfrastructure]{}

	detachFirewallFromServer := func(infra *EnvInfrastructure) error {
		if infra.Firewall == nil || !infra.Firewall.AttachedToServer {
			return nil
		}

		err := infrastructure.DetachFirewallFromServer(
			hetznerClient,
			infra.Firewall.ID,
			infra.Server.ID,
		)

		if err != nil {
			return err
		}

		infra.Firewall.AttachedToServer = false
		return nil
	}

	envInfraQueue = append(
		envInfraQueue,
		queues.InfrastructureQueueSteps[*EnvInfrastructure]{
			func(*EnvInfrastructure) error {
				stepper.StartTemporaryStep("Detaching firewall from server")
				return nil
			},
			detachFirewallFromServer,
		},
	)

	removeServer := func(infra *EnvInfrastructure) error {
		if infra.Server == nil {
			return nil
		}

		err := infrastructure.RemoveServer(
			hetznerClient,
			infra.Server.ID,
		)

		if err != nil {
			return err
		}

		infra.Server = nil
		return nil
	}

	envInfraQueue = append(
		envInfraQueue,
		queues.InfrastructureQueueSteps[*EnvInfrastructure]{
			func(*EnvInfrastructure) error {
				stepper.StartTemporaryStep("Removing the server")
				return nil
			},
			removeServer,
		},
	)

	removeSSHKey := func(infra *EnvInfrastructure) error {
		if infra.SSHKey == nil {
			return nil
		}

		err := infrastructure.RemoveSSHKey(
			hetznerClient,
			infra.SSHKey.ID,
		)

		if err != nil {
			return err
		}

		infra.SSHKey = nil
		return nil
	}

	removeFirewall := func(infra *EnvInfrastructure) error {
		if infra.Firewall == nil {
			return nil
		}

		err := infrastructure.RemoveFirewall(
			hetznerClient,
			infra.Firewall.ID,
		)

		if err != nil {
			return err
		}

		infra.Firewall = nil
		return nil
	}

	envInfraQueue = append(
		envInfraQueue,
		queues.InfrastructureQueueSteps[*EnvInfrastructure]{
			func(*EnvInfrastructure) error {
				stepper.StartTemporaryStep("Removing the SSH key and the firewall")
				return nil
			},
			removeSSHKey,
			removeFirewall,
		},
	)

	err = envInfraQueue.Run(
		envInfra,
	)

	// Env infra could be updated in the queue even
	// in case of error (partial infrastructure)
	env.SetInfrastructureJSON(envInfra)

	return err
}
