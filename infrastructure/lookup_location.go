package infrastructure

import (
	"context"
	"fmt"

	"github.com/hetznercloud/hcloud-go/hcloud"
)

type Location struct {
	ID          int                `json:"id"`
	Name        string             `json:"name"`
	NetworkZone hcloud.NetworkZone `json:"network_zone"`
}

func LookupLocation(
	hetznerClient *hcloud.Client,
	name string,
) (returnedLocation *Location, returnedError error) {

	location, _, err := hetznerClient.Location.GetByName(
		context.TODO(),
		name,
	)

	if err != nil {
		returnedError = err
		return
	}

	if location == nil {
		returnedError = fmt.Errorf("invalid location %s", name)
		return
	}

	returnedLocation = &Location{
		ID:          location.ID,
		Name:        location.Name,
		NetworkZone: location.NetworkZone,
	}
	return
}
