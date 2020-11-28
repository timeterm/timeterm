package nmsdk

import (
	"errors"
	"fmt"
	"os"
	"path"

	vault "github.com/hashicorp/vault/api"
	"github.com/nats-io/jwt/v2"
	"github.com/nats-io/nats.go"
)

func makeAppCredsPath(mount, prefix, app string) string {
	return path.Join(mount, "/data/", prefix, "/apps/creds/", app)
}

type AppCredsRetriever struct {
	vaultMount  string
	vaultPrefix string
	appName     string
	vault       *vault.Client
}

func NewAppCredsRetrieverFromEnv(appName string) (AppCredsRetriever, error) {
	prefix := os.Getenv("NATS_MANAGER_VAULT_PREFIX")
	mount := os.Getenv("NATS_MANAGER_VAULT_MOUNT")

	vc, err := vault.NewClient(vault.DefaultConfig())
	if err != nil {
		return AppCredsRetriever{}, err
	}

	return AppCredsRetriever{
		vaultMount:  mount,
		vaultPrefix: prefix,
		appName:     appName,
		vault:       vc,
	}, nil
}

func NewAppCredsRetriever(appName, vaultMount, vaultPrefix string, vault *vault.Client) AppCredsRetriever {
	return AppCredsRetriever{
		appName:     appName,
		vaultMount:  vaultMount,
		vaultPrefix: vaultPrefix,
		vault:       vault,
	}
}

func wipeBytes(bs []byte) {
	for i := range bs {
		bs[i] = 'X'
	}
}

func (r AppCredsRetriever) NatsCredsCBs() (nats.UserJWTHandler, nats.SignatureHandler) {
	getCreds := func() ([]byte, error) {
		secret, err := r.vault.Logical().Read(makeAppCredsPath(r.vaultMount, r.vaultPrefix, r.appName))
		if err != nil {
			return nil, fmt.Errorf("could not read app credentials from Vault: %w", err)
		}

		secretData, ok := secret.Data["data"].(map[string]interface{})
		if !ok {
			return nil, fmt.Errorf("invalid secret data: %w", err)
		}
		creds, ok := secretData["creds"].(string)
		if !ok {
			return nil, errors.New("could not extract credentials from secret")
		}
		return []byte(creds), nil
	}

	jwtCB := func() (string, error) {
		creds, err := getCreds()
		if err != nil {
			return "", err
		}
		defer wipeBytes(creds)
		return jwt.ParseDecoratedJWT(creds)
	}

	signCB := func(nonce []byte) ([]byte, error) {
		creds, err := getCreds()
		if err != nil {
			return nil, err
		}
		defer wipeBytes(creds)

		nkey, err := jwt.ParseDecoratedUserNKey(creds)
		if err != nil {
			return nil, err
		}
		defer nkey.Wipe()

		return nkey.Sign(nonce)
	}

	return jwtCB, signCB
}
