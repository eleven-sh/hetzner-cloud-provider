package infrastructure

import (
	"context"
	"net"

	"github.com/hetznercloud/hcloud-go/hcloud"
)

type Network struct {
	ID int `json:"id"`
}

func CreateNetwork(
	hetznerClient *hcloud.Client,
	name string,
	cidrBlock string,
	subnetCIDRBlock string,
	networkZone hcloud.NetworkZone,
) (returnedNetwork *Network, returnedError error) {

	_, ipRange, err := net.ParseCIDR(cidrBlock)

	if err != nil {
		returnedError = err
		return
	}

	_, subnetIPRange, err := net.ParseCIDR(subnetCIDRBlock)

	if err != nil {
		returnedError = err
		return
	}

	createNetworkResp, _, err := hetznerClient.Network.Create(
		context.TODO(),
		hcloud.NetworkCreateOpts{
			Name: name,
			IPRange: &net.IPNet{
				IP:   ipRange.IP,
				Mask: ipRange.Mask,
			},
			Subnets: []hcloud.NetworkSubnet{
				{
					Type:        hcloud.NetworkSubnetTypeCloud,
					NetworkZone: networkZone,
					IPRange: &net.IPNet{
						IP:   subnetIPRange.IP,
						Mask: subnetIPRange.Mask,
					},
				},
			},
		},
	)

	if err != nil {
		returnedError = err
		return
	}

	returnedNetwork = &Network{
		ID: createNetworkResp.ID,
	}
	return
}
