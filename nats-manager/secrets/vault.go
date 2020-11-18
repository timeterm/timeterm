package secrets

import (
	"errors"
	"fmt"
	"path"

	vault "github.com/hashicorp/vault/api"
	"github.com/nats-io/jwt/v2"
	"github.com/nats-io/nkeys"
)

type VaultClient struct {
	prefix string
	vault  *vault.Client
}

func (c *VaultClient) operatorSeedPath(pubKey string) string {
	return path.Join(c.prefix, "/keys/operator/", pubKey)
}

func (c *VaultClient) accountSeedPath(pubKey string) string {
	return path.Join(c.prefix, "/keys/account/", pubKey)
}

func (c *VaultClient) userSeedPath(pubKey string) string {
	return path.Join(c.prefix, "/keys/user/", pubKey)
}

func (c *VaultClient) operatorJWTPath(name string) string {
	return path.Join(c.prefix, "/jwts/operator/", name)
}

func (c *VaultClient) accountJWTPath(name string) string {
	return path.Join(c.prefix, "/jwts/account/", name)
}

func (c *VaultClient) userJWTPath(accountName, userName string) string {
	return path.Join(c.prefix, "/jwts/account/", accountName, "/user/", userName)
}

func (c *VaultClient) WriteOperatorSeed(kp nkeys.KeyPair) error {
	pubKey, err := kp.PublicKey()
	if err != nil {
		return err
	}

	pat := c.operatorSeedPath(pubKey)
	err = c.writeSeed(pat, kp)
	if err != nil {
		return fmt.Errorf("could not write operator seed: %w", err)
	}
	return nil
}

func (c *VaultClient) WriteAccountSeed(kp nkeys.KeyPair) error {
	pubKey, err := kp.PublicKey()
	if err != nil {
		return err
	}

	pat := c.accountSeedPath(pubKey)
	err = c.writeSeed(pat, kp)
	if err != nil {
		return fmt.Errorf("could not write account seed: %w", err)
	}
	return nil
}

func (c *VaultClient) WriteUserSeed(kp nkeys.KeyPair) error {
	pubKey, err := kp.PublicKey()
	if err != nil {
		return err
	}

	pat := c.userSeedPath(pubKey)
	err = c.writeSeed(pat, kp)
	if err != nil {
		return fmt.Errorf("could not write user seed: %w", err)
	}
	return nil
}

func (c *VaultClient) WriteOperatorJWT(claims *jwt.OperatorClaims, operatorPubKey string) error {
	if claims.Name == "" {
		return errors.New("operator name may not be empty")
	}
	kp, err := c.ReadOperatorSeed(operatorPubKey)
	if err != nil {
		return fmt.Errorf("could not read operator seed: %w", err)
	}

	pat := c.operatorJWTPath(claims.Name)
	return c.writeJWT(pat, claims, kp)
}

func (c *VaultClient) WriteUserJWT(claims *jwt.UserClaims, accountName string, accountPubKey string) error {
	if claims.Name == "" {
		return errors.New("user name may not be empty")
	}
	if accountName == "" {
		return errors.New("account name may not be empty")
	}
	kp, err := c.ReadAccountSeed(accountPubKey)
	if err != nil {
		return fmt.Errorf("could not read account seed: %w", err)
	}

	pat := c.userJWTPath(accountName, claims.Name)
	return c.writeJWT(pat, claims, kp)
}

func (c *VaultClient) WriteAccountJWT(claims *jwt.AccountClaims, operatorPubKey string) error {
	if claims.Name == "" {
		return errors.New("account name may not be empty")
	}
	kp, err := c.ReadOperatorSeed(operatorPubKey)
	if err != nil {
		return fmt.Errorf("could not read operator seed: %w", err)
	}

	pat := c.accountJWTPath(claims.Name)
	return c.writeJWT(pat, claims, kp)
}

func (c *VaultClient) writeJWT(pat string, claims jwt.Claims, kp nkeys.KeyPair) error {
	encoded, err := claims.Encode(kp)
	if err != nil {
		return fmt.Errorf("could not encode operator claims: %w", err)
	}

	_, err = c.vault.Logical().Write(pat, map[string]interface{}{
		"jwt": encoded,
	})
	return err
}

func (c *VaultClient) ReadOperatorSeed(pubKey string) (nkeys.KeyPair, error) {
	pat := c.operatorSeedPath(pubKey)
	kp, err := c.readSeed(pat)
	if err != nil {
		return kp, fmt.Errorf("could not read operator seed: %w", err)
	}
	return kp, nil
}

func (c *VaultClient) ReadUserSeed(pubKey string) (nkeys.KeyPair, error) {
	pat := c.userSeedPath(pubKey)
	kp, err := c.readSeed(pat)
	if err != nil {
		return kp, fmt.Errorf("could not read user seed: %w", err)
	}
	return kp, nil
}

func (c *VaultClient) ReadAccountSeed(pubKey string) (nkeys.KeyPair, error) {
	pat := c.accountSeedPath(pubKey)
	kp, err := c.readSeed(pat)
	if err != nil {
		return kp, fmt.Errorf("could not read account seed: %w", err)
	}
	return kp, nil
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
