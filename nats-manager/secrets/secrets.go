package secrets

import (
	"fmt"
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

	pat := path.Join(c.prefix, "/operator")
	err = c.writeSeed(pat, kp)	
	if err != nil {
		return fmt.Errorf("could not write operator seed: %w", err)
	}
	return nil 
}

func (c *VaultClient) WriteAccount(kp nkeys.KeyPair) error {
	pubKey, err := kp.PublicKey()
	if err != nil {
		return err
	}

	pat := path.Join(c.prefix, "/account/", pubKey)
	err = c.writeSeed(pat, kp)
	if err != nil {
		return fmt.Errorf("could not write account seed: %w", err)
	}
	return nil
}

func (c *VaultClient) WriteUser(kp nkeys.KeyPair) error {
	pubKey, err := kp.PublicKey()
	if err != nil {
		return err
	}

	pat := path.Join(c.prefix, "/user/", pubKey)
	err = c.writeSeed(pat, kp)
	if err != nil {
		return fmt.Errorf("could not write user seed: %w", err)
	}
	return nil
}

func (c *VaultClient) ReadOperator() (nkeys.KeyPair, error) {
	pat := path.Join(c.prefix, "/operator")
	kp, err := c.readSeed(pat)
	if err != nil {
		return kp, fmt.Errorf("could not read operator seed: %w", err)
	}
	return kp, nil
}

func (c *VaultClient) ReadUser(pubKey string) (nkeys.KeyPair, error) {
	pat := path.Join(c.prefix, "/user/", pubKey)
	kp, err := c.readSeed(pat)
	if err != nil {
		return kp, fmt.Errorf("could not read user seed: %w", err)
	}
	return kp, nil
}

func (c *VaultClient) ReadAccount(pubKey string) (nkeys.KeyPair, error) {
	pat := path.Join(c.prefix, "/account/", pubKey)
	kp, err := c.readSeed(pat)
	if err != nil {
		return kp, fmt.Errorf("could not read account seed: %w", err)
	}
	return kp, nil
}

type Manager struct {
	safe *vault.Client
}

func (d *Manager) Init() error {
	return d.InitKeys()
}

func (d *Manager) InitKeys() error {
	kp, err := nkeys.CreateOperator()
	if err != nil {
		return fmt.Errorf("could not create operator keys: %w", err)
	}


}

func (c *VaultClient) writeSeed(path string, kp nkeys.KeyPair) error {
	seed, err := kp.Seed()
	if err != nil {
		return err
	}

	_, err = c.vault.Logical().Write(path, map[string]interface{}{
		"seed": seed,
	})
	return err
}

func (c *VaultClient) readSeed(path string) (nkeys.KeyPair, error) {
	var kp nkeys.KeyPair

	secret, err := c.vault.Logical().Read(path)
	if err != nil {
		return kp, err
	}

	seed, ok := secret.Data["seed"].(string)
	if !ok {
		return kp, fmt.Errorf("seed not present in secret at path %s", path)
	}

	kp, err = nkeys.FromSeed([]byte(seed))
	if err != nil {
		return kp, fmt.Errorf("could not create key pair from seed: %w", err)
	}
	return kp, nil
}
