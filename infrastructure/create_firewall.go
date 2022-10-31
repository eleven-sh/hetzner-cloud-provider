package infrastructure

import (
	"context"

	"github.com/hetznercloud/hcloud-go/hcloud"
)

type Firewall struct {
	ID               int  `json:"id"`
	AttachedToServer bool `json:"attached_to_server"`
}

func CreateFirewall(
	hetznerClient *hcloud.Client,
	name string,
	rules []hcloud.FirewallRule,
) (returnedFirewall *Firewall, returnedError error) {

	createFirewallResp, _, err := hetznerClient.Firewall.Create(
		context.TODO(),
		hcloud.FirewallCreateOpts{
			Name:  name,
			Rules: rules,
		},
	)

	if err != nil {
		returnedError = err
		return
	}

	_, errChan := hetznerClient.Action.WatchOverallProgress(
		context.TODO(),
		createFirewallResp.Actions,
	)

	if err = <-errChan; err != nil {
		returnedError = err
		return
	}

	returnedFirewall = &Firewall{
		ID:               createFirewallResp.Firewall.ID,
		AttachedToServer: false,
	}
	return
}
