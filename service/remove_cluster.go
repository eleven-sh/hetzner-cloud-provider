package service

import (
	"encoding/json"

	"github.com/eleven-sh/eleven/entities"
	"github.com/eleven-sh/eleven/queues"
	"github.com/eleven-sh/eleven/stepper"
	"github.com/eleven-sh/hetzner-cloud-provider/infrastructure"
	"github.com/hetznercloud/hcloud-go/hcloud"
)

func (h *Hetzner) RemoveCluster(
	stepper stepper.Stepper,
	config *entities.Config,
	cluster *entities.Cluster,
) error {

	var clusterInfra *ClusterInfrastructure
	err := json.Unmarshal([]byte(cluster.InfrastructureJSON), &clusterInfra)

	if err != nil {
		return err
	}

	hetznerClient := hcloud.NewClient(hcloud.WithToken(h.config.Credentials.APIToken))
	clusterInfraQueue := queues.InfrastructureQueue[*ClusterInfrastructure]{}

	removeNetwork := func(infra *ClusterInfrastructure) error {
		if infra.Network == nil {
			return nil
		}

		err := infrastructure.RemoveNetwork(
			hetznerClient,
			infra.Network.ID,
		)

		if err != nil {
			return err
		}

		infra.Network = nil
		return nil
	}

	clusterInfraQueue = append(
		clusterInfraQueue,
		queues.InfrastructureQueueSteps[*ClusterInfrastructure]{
			func(*ClusterInfrastructure) error {
				stepper.StartTemporaryStep("Removing the network")
				return nil
			},
			removeNetwork,
		},
	)

	err = clusterInfraQueue.Run(
		clusterInfra,
	)

	// Cluster infra could be updated in the queue even
	// in case of error (partial infrastructure)
	cluster.SetInfrastructureJSON(clusterInfra)

	return err
}
