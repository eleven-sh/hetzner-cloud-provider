package infrastructure

import (
	"context"

	"github.com/hetznercloud/hcloud-go/hcloud"
)

func RemoveServer(
	hetznerClient *hcloud.Client,
	serverID int,
) error {

	_, err := hetznerClient.Server.Delete(
		context.TODO(),
		&hcloud.Server{
			ID: serverID,
		},
	)

	return err
}
