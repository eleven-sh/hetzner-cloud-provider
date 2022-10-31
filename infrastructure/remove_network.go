package infrastructure

import (
	"context"

	"github.com/hetznercloud/hcloud-go/hcloud"
)

func RemoveNetwork(
	hetznerClient *hcloud.Client,
	networkID int,
) error {

	_, err := hetznerClient.Network.Delete(
		context.TODO(),
		&hcloud.Network{
			ID: networkID,
		},
	)

	return err
}
