package infrastructure

import (
	"context"

	"github.com/hetznercloud/hcloud-go/hcloud"
)

func DetachFirewallFromServer(
	hetznerClient *hcloud.Client,
	firewallID int,
	serverID int,
) error {

	detachFirewallResp, _, err := hetznerClient.Firewall.RemoveResources(
		context.TODO(),
		&hcloud.Firewall{
			ID: firewallID,
		},
		[]hcloud.FirewallResource{
			{
				Type: hcloud.FirewallResourceTypeServer,
				Server: &hcloud.FirewallResourceServer{
					ID: serverID,
				},
			},
		},
	)

	if err != nil {
		return err
	}

	_, errChan := hetznerClient.Action.WatchOverallProgress(
		context.TODO(),
		detachFirewallResp,
	)

	if err = <-errChan; err != nil {
		return err
	}

	return nil
}
