package secrets

import (
	"path"

	vault "github.com/hashicorp/vault/api"
	"github.com/nats-io/nkeys"
)

type VaultClient struct {
	prefix string
	vault  *vault.Client
}

func (c *VaultClient) WriteOperator(kp nkeys.KeyPair) error {
	pubKey, err := kp.PublicKey()
	if err != nil {
		return err
	}

	seed, err := kp.Seed()
	if err != nil {
		return err
	}

	pat := path.Join(c.prefix, "/operator/", pubKey)
	_, err = c.vault.Logical().Write(pat, map[string]interface{}{
		"seed": seed,
	})

	return err
}

type Manager struct {
	safe *vault.Client
}

func (d *Manager) Init() error {
	return d.InitKeys()
}

func (d *Manager) InitKeys() error {
	d.safe.Logical().Write("/timeterm/nats-manager")
}
