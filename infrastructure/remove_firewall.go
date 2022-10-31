package infrastructure

import (
	"context"

	"github.com/hetznercloud/hcloud-go/hcloud"
)

func RemoveFirewall(
	hetznerClient *hcloud.Client,
	firewallID int,
) error {

	_, err := hetznerClient.Firewall.Delete(
		context.TODO(),
		&hcloud.Firewall{
			ID: firewallID,
		},
	)

	return err
}
