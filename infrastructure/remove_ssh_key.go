package infrastructure

import (
	"context"

	"github.com/hetznercloud/hcloud-go/hcloud"
)

func RemoveSSHKey(
	hetznerClient *hcloud.Client,
	sshKeyID int,
) error {

	_, err := hetznerClient.SSHKey.Delete(
		context.TODO(),
		&hcloud.SSHKey{
			ID: sshKeyID,
		},
	)

	return err
}
