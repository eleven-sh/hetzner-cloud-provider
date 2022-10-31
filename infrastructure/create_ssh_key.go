package infrastructure

import (
	"context"
	"crypto/ed25519"
	"crypto/rand"
	"encoding/pem"

	"github.com/hetznercloud/hcloud-go/hcloud"
	"github.com/mikesmitty/edkey"
	"golang.org/x/crypto/ssh"
)

type SSHKey struct {
	ID         int    `json:"id"`
	Name       string `json:"name"`
	PEMContent string `json:"pem_content"`
}

func CreateSSHKey(
	hetznerClient *hcloud.Client,
	name string,
) (returnedSSHKey *SSHKey, returnedError error) {

	publicKey, privateKey, err := generateEd25519Keys()

	if err != nil {
		return nil, err
	}

	createSSHKeyResp, _, err := hetznerClient.SSHKey.Create(
		context.TODO(),
		hcloud.SSHKeyCreateOpts{
			Name:      name,
			PublicKey: publicKey,
		},
	)

	if err != nil {
		returnedError = err
		return
	}

	returnedSSHKey = &SSHKey{
		ID:         createSSHKeyResp.ID,
		Name:       createSSHKeyResp.Name,
		PEMContent: privateKey,
	}
	return
}

func generateEd25519Keys() (
	returnedPubkey string,
	returnedPrivateKey string,
	returnedError error,
) {
	pubKey, privKey, err := ed25519.GenerateKey(rand.Reader)

	if err != nil {
		returnedError = err
		return
	}

	publicKey, err := ssh.NewPublicKey(pubKey)

	if err != nil {
		returnedError = err
		return
	}

	pemKey := &pem.Block{
		Type:  "OPENSSH PRIVATE KEY",
		Bytes: edkey.MarshalED25519PrivateKey(privKey),
	}

	privateKey := pem.EncodeToMemory(pemKey)
	authorizedKey := ssh.MarshalAuthorizedKey(publicKey)

	return string(authorizedKey), string(privateKey), nil
}
