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

func NewVaultClient(prefix string, c *vault.Client) *VaultClient {
	return &VaultClient{
		prefix: prefix,
		vault:  c,
	}
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

func (c *VaultClient) readJWT(pat string) (string, error) {
	secret, err := c.vault.Logical().Read(pat)
	if err != nil {
		return "", err
	}

	token, ok := secret.Data["jwt"].(string)
	if !ok {
		return "", fmt.Errorf("jwt not present in secret at path %s", pat)
	}
	return token, nil
}

func (c *VaultClient) ReadOperatorJWT(name string) (*jwt.OperatorClaims, error) {
	pat := c.operatorJWTPath(name)
	token, err := c.readJWT(pat)
	if err != nil {
		return nil, fmt.Errorf("could not read operator JWT at path %s from Vault: %w", pat, err)
	}

	claims, err := jwt.DecodeOperatorClaims(token)
	if err != nil {
		return nil, fmt.Errorf("could not decode operator claims: %w", err)
	}
	return claims, nil
}

func (c *VaultClient) ReadAccountJWT(name string) (*jwt.AccountClaims, error) {
	pat := c.accountJWTPath(name)
	token, err := c.readJWT(pat)
	if err != nil {
		return nil, fmt.Errorf("could not read account JWT at path %s from Vault: %w", pat, err)
	}

	claims, err := jwt.DecodeAccountClaims(token)
	if err != nil {
		return nil, fmt.Errorf("could not decode account claims: %w", err)
	}
	return claims, nil
}

func (c *VaultClient) ReadUserJWT(accountName, userName string) (*jwt.UserClaims, error) {
	pat := c.userJWTPath(accountName, userName)
	token, err := c.readJWT(pat)
	if err != nil {
		return nil, fmt.Errorf("could not read user JWT at path %s from Vault: %w", pat, err)
	}

	claims, err := jwt.DecodeUserClaims(token)
	if err != nil {
		return nil, fmt.Errorf("could not decode user claims: %w", err)
	}
	return claims, nil
}

func (c *VaultClient) ReadOperatorSeed(pubKey string) (nkeys.KeyPair, error) {
	pat := c.operatorSeedPath(pubKey)
	return c.readSeed(pat)
}

func (c *VaultClient) ReadUserSeed(pubKey string) (nkeys.KeyPair, error) {
	pat := c.userSeedPath(pubKey)
	return c.readSeed(pat)
}

func (c *VaultClient) ReadAccountSeed(pubKey string) (nkeys.KeyPair, error) {
	pat := c.accountSeedPath(pubKey)
	return c.readSeed(pat)
}

func (c *VaultClient) writeSeed(path string, kp nkeys.KeyPair) error {
	seed, err := kp.Seed()
	if err != nil {
		return err
	}

	// Converting seed to a string is safe (or should really be safe)
	// because it should be in Base 32.
	seedStr := string(seed)
	_, err = c.vault.Logical().Write(path, map[string]interface{}{
		"seed": seedStr,
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
