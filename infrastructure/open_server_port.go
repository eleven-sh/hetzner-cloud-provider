package infrastructure

import (
	"context"
	"fmt"
	"net"

	"github.com/hetznercloud/hcloud-go/hcloud"
)

func OpenServerPort(
	hetznerClient *hcloud.Client,
	firewallID int,
	portToOpen string,
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

	_, sourceIP, err := net.ParseCIDR("0.0.0.0/0")

	if err != nil {
		return err
	}

	updatedRules := firewall.Rules
	updatedRules = append(updatedRules, hcloud.FirewallRule{
		Direction: hcloud.FirewallRuleDirectionIn,
		Protocol:  hcloud.FirewallRuleProtocolTCP,
		Port:      &portToOpen,
		SourceIPs: []net.IPNet{
			{
				IP:   sourceIP.IP,
				Mask: sourceIP.Mask,
			},
		},
	})

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
