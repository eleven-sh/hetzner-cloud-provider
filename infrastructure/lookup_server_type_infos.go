package infrastructure

import (
	"context"
	"errors"

	"github.com/hetznercloud/hcloud-go/hcloud"
)

var (
	ErrInvalidServerType = errors.New("ErrInvalidServerType")
)

type ServerTypeInfos struct {
	ID   int    `json:"id"`
	Name string `json:"name"`
}

func LookupServerTypeInfos(
	hetznerClient *hcloud.Client,
	serverTypeName string,
	locationName string,
) (returnedServerTypeInfos *ServerTypeInfos, returnedError error) {

	serverType, _, err := hetznerClient.ServerType.GetByName(
		context.TODO(),
		serverTypeName,
	)

	if err != nil {
		returnedError = err
		return
	}

	if serverType == nil {
		returnedError = ErrInvalidServerType
		return
	}

	serverTypeExistsInLocation := false
	for _, price := range serverType.Pricings {
		if price.Location.Name == locationName {
			serverTypeExistsInLocation = true
		}
	}

	if !serverTypeExistsInLocation {
		returnedError = ErrInvalidServerType
		return
	}

	returnedServerTypeInfos = &ServerTypeInfos{
		ID:   serverType.ID,
		Name: serverType.Name,
	}
	return
}
