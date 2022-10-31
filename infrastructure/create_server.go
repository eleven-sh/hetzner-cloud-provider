package infrastructure

import (
	"context"
	_ "embed"
	"strings"

	"github.com/eleven-sh/agent/config"
	"github.com/hetznercloud/hcloud-go/hcloud"
)

const (
	ServerImage = "ubuntu-22.04"
)

var (
	//go:embed init_server.sh
	serverInitScript string
)

type Server struct {
	ID                int                      `json:"id"`
	Type              string                   `json:"type"`
	PublicIPAddress   string                   `json:"tmp_public_ip_address"`
	InitScriptResults *InitServerScriptResults `json:"init_script_results"`
}

func CreateServer(
	hetznerClient *hcloud.Client,
	locationID int,
	networkID int,
	serverType string,
	firewallID int,
	name string,
	sshKeyID int,
) (returnedServer *Server, returnedError error) {

	updatedServerInitScript := strings.ReplaceAll(
		serverInitScript,
		"${ELEVEN_CONFIG_DIR}",
		config.ElevenConfigDirPath,
	)

	updatedServerInitScript = strings.ReplaceAll(
		updatedServerInitScript,
		"${ELEVEN_AGENT_CONFIG_DIR}",
		config.ElevenAgentConfigDirPath,
	)

	createServerResp, _, err := hetznerClient.Server.Create(
		context.TODO(),
		hcloud.ServerCreateOpts{
			Name: name,
			ServerType: &hcloud.ServerType{
				Name: serverType,
			},
			Firewalls: []*hcloud.ServerCreateFirewall{
				{
					Firewall: hcloud.Firewall{
						ID: firewallID,
					},
				},
			},
			Image: &hcloud.Image{
				Name: ServerImage,
			},
			Location: &hcloud.Location{
				ID: locationID,
			},
			Networks: []*hcloud.Network{
				{
					ID: networkID,
				},
			},
			SSHKeys: []*hcloud.SSHKey{
				{
					ID: sshKeyID,
				},
			},
			UserData: updatedServerInitScript,
		},
	)

	if err != nil {
		returnedError = err
		return
	}

	_, errChan := hetznerClient.Action.WatchOverallProgress(
		context.TODO(),
		append(
			[]*hcloud.Action{createServerResp.Action},
			createServerResp.NextActions...,
		),
	)

	if err = <-errChan; err != nil {
		returnedError = err
		return
	}

	returnedServer = &Server{
		ID:              createServerResp.Server.ID,
		Type:            createServerResp.Server.ServerType.Name,
		PublicIPAddress: createServerResp.Server.PublicNet.IPv4.IP.String(),
	}
	return
}
