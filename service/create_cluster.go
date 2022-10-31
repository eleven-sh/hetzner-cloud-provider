package service

import (
	"encoding/json"

	"github.com/eleven-sh/eleven/entities"
	"github.com/eleven-sh/eleven/queues"
	"github.com/eleven-sh/eleven/stepper"
	"github.com/eleven-sh/hetzner-cloud-provider/infrastructure"
	"github.com/hetznercloud/hcloud-go/hcloud"
)

type ClusterInfrastructure struct {
	Network  *infrastructure.Network  `json:"network"`
	Location *infrastructure.Location `json:"location"`
}

func (h *Hetzner) CreateCluster(
	stepper stepper.Stepper,
	config *entities.Config,
	cluster *entities.Cluster,
) error {

	clusterInfra := &ClusterInfrastructure{}
	if len(cluster.InfrastructureJSON) > 0 {
		err := json.Unmarshal([]byte(cluster.InfrastructureJSON), clusterInfra)

		if err != nil {
			return err
		}
	}

	prefixResource := prefixClusterResource(cluster.GetNameSlug())
	hetznerClient := hcloud.NewClient(hcloud.WithToken(h.config.Credentials.APIToken))

	clusterInfraQueue := queues.InfrastructureQueue[*ClusterInfrastructure]{}

	lookupLocation := func(infra *ClusterInfrastructure) error {
		if infra.Location != nil {
			return nil
		}

		location, err := infrastructure.LookupLocation(
			hetznerClient,
			h.config.Region,
		)

		if err != nil {
			return err
		}

		infra.Location = location
		return nil
	}

	clusterInfraQueue = append(
		clusterInfraQueue,
		queues.InfrastructureQueueSteps[*ClusterInfrastructure]{
			func(*ClusterInfrastructure) error {
				stepper.StartTemporaryStep("Looking up location")
				return nil
			},
			lookupLocation,
		},
	)

	createNetwork := func(infra *ClusterInfrastructure) error {
		if infra.Network != nil {
			return nil
		}

		network, err := infrastructure.CreateNetwork(
			hetznerClient,
			prefixResource("network"),
			"10.0.0.0/16",
			"10.0.0.0/24",
			infra.Location.NetworkZone,
		)

		if err != nil {
			return err
		}

		infra.Network = network
		return nil
	}

	clusterInfraQueue = append(
		clusterInfraQueue,
		queues.InfrastructureQueueSteps[*ClusterInfrastructure]{
			func(*ClusterInfrastructure) error {
				stepper.StartTemporaryStep("Creating a network")
				return nil
			},
			createNetwork,
		},
	)

	err := clusterInfraQueue.Run(
		clusterInfra,
	)

	// Cluster infra could be updated in the queue even
	// in case of error (partial infrastructure)
	cluster.SetInfrastructureJSON(clusterInfra)

	return err
}
