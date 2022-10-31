package infrastructure

import (
	"context"
	"fmt"

	"github.com/hetznercloud/hcloud-go/hcloud"
)

func CloseServerPort(
	hetznerClient *hcloud.Client,
	firewallID int,
	portToClose string,
) error {

	firewall, _, err := hetznerClient.Firewall.GetByID(
		context.TODO(),
		firewallID,
	)

	if err != nil {
		return err
	}

	if firewall == nil {
		return fmt.Errorf("invalid firewall %d", firewallID)
	}

	updatedRules := []hcloud.FirewallRule{}

	for _, rule := range firewall.Rules {
		if *rule.Port == portToClose {
			continue
		}

		updatedRules = append(updatedRules, rule)
	}

	updateFirewallResp, _, err := hetznerClient.Firewall.SetRules(
		context.TODO(),
		firewall,
		hcloud.FirewallSetRulesOpts{
			Rules: updatedRules,
		},
	)

	if err != nil {
		return err
	}

	_, errChan := hetznerClient.Action.WatchOverallProgress(
		context.TODO(),
		updateFirewallResp,
	)

	if err = <-errChan; err != nil {
		return err
	}

	return nil
}
