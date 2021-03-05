package secrets

import (
	"crypto/rand"
	"encoding/base64"
	"errors"
	"fmt"
	"io"

	"github.com/google/uuid"
	vault "github.com/hashicorp/vault/api"
	"google.golang.org/protobuf/proto"

	devcfgpb "gitlab.com/timeterm/timeterm/proto/go/devcfg"
)

type Wrapper struct {
	c *vault.Client

	mount  string
	prefix string
}

func New(mount, prefix string) (*Wrapper, error) {
	client, err := vault.NewClient(vault.DefaultConfig())
	if err != nil {
		return nil, err
	}
	return &Wrapper{
		c:      client,
		mount:  mount,
		prefix: prefix,
	}, nil
}

func (w *Wrapper) createNetworkingServiceSecretPath(id uuid.UUID) string {
	return fmt.Sprintf("%s/data/%s/networking/service/%s", w.mount, w.prefix, id)
}

func (w *Wrapper) createOrganizationZermeloTokenSecretPath(organizationID uuid.UUID) string {
	return fmt.Sprintf("%s/data/%s/zermelo/tokens/organization/%s", w.mount, w.prefix, organizationID)
}

func (w *Wrapper) createOrganizationLogsKeySecretPath(organizationID uuid.UUID) string {
	return fmt.Sprintf("%s/data/%s/logskey/%s", w.mount, w.prefix, organizationID)
}

func (w *Wrapper) GetNetworkingService(id uuid.UUID) (*devcfgpb.NetworkingService, error) {
	secretPath := w.createNetworkingServiceSecretPath(id)
	secret, err := w.c.Logical().Read(secretPath)
	if err != nil {
		return nil, err
	}
	if secret == nil {
		return nil, nil
	}

	secretData, ok := secret.Data["data"].(map[string]interface{})
	if !ok {
		return nil, errors.New("invalid secret data (may not be present)")
	}
	bytes, ok := secretData["config"].(string)
	if !ok {
		return nil, errors.New("config not present in secret")
	}

	decoded, err := base64.StdEncoding.DecodeString(bytes)
	if err != nil {
		return nil, err
	}

	var networkingService devcfgpb.NetworkingService

	err = proto.Unmarshal(decoded, &networkingService)
	if err != nil {
		return nil, err
	}
	return &networkingService, nil
}

func (w *Wrapper) DeleteNetworkingService(id uuid.UUID) error {
	secretPath := w.createNetworkingServiceSecretPath(id)
	_, err := w.c.Logical().Delete(secretPath)
	return err
}

func (w *Wrapper) UpsertNetworkingService(id uuid.UUID, cfg *devcfgpb.NetworkingService) error {
	bytes, err := proto.Marshal(cfg)
	if err != nil {
		return err
	}

	encoded := base64.StdEncoding.EncodeToString(bytes)

	secretPath := w.createNetworkingServiceSecretPath(id)

	_, err = w.c.Logical().Write(secretPath, map[string]interface{}{
		"data": map[string]interface{}{
			"config": encoded,
		},
	})
	return err
}

func (w *Wrapper) GetOrganizationZermeloToken(organizationID uuid.UUID) ([]byte, error) {
	secretPath := w.createOrganizationZermeloTokenSecretPath(organizationID)
	secret, err := w.c.Logical().Read(secretPath)
	if err != nil {
		return nil, err
	}
	if secret == nil {
		return nil, err
	}

	secretData, ok := secret.Data["data"].(map[string]interface{})
	if !ok {
		return nil, errors.New("invalid secret data (may not be present)")
	}
	token, ok := secretData["token"].(string)
	if !ok {
		return nil, errors.New("token not present in secret")
	}

	return []byte(token), nil
}

func (w *Wrapper) UpsertOrganizationZermeloToken(id uuid.UUID, token []byte) error {
	secretPath := w.createOrganizationZermeloTokenSecretPath(id)

	_, err := w.c.Logical().Write(secretPath, map[string]interface{}{
		"data": map[string]interface{}{
			"token": string(token),
		},
	})
	return err
}

func (w *Wrapper) GetOrganizationLogsKeySecret(organizationID uuid.UUID) ([]byte, error) {
	secretPath := w.createOrganizationLogsKeySecretPath(organizationID)
	secret, err := w.c.Logical().Read(secretPath)
	if err != nil {
		return nil, err
	}
	if secret == nil {
		return nil, err
	}

	secretData, ok := secret.Data["data"].(map[string]interface{})
	if !ok {
		return nil, errors.New("invalid secret data (may not be present)")
	}
	bytes, ok := secretData["key"].(string)
	if !ok {
		return nil, errors.New("logs key not present in secret")
	}

	key, err := base64.StdEncoding.DecodeString(bytes)
	if err != nil {
		return nil, fmt.Errorf("could not decode key")
	}

	return []byte(key), nil
}

func (w *Wrapper) NewOrganizationLogsKeySecret(organizationID uuid.UUID) error {
	// The key is an AES(-256) key, so we're going for 32 bytes.
	key := make([]byte, 32)
	if _, err := io.ReadFull(rand.Reader, key); err != nil {
		return fmt.Errorf("could not generate key: %w", err)
	}
	return w.UpsertOrganizationLogsKeySecret(organizationID, key)
}

func (w *Wrapper) UpsertOrganizationLogsKeySecret(organizationID uuid.UUID, key []byte) error {
	secretPath := w.createOrganizationLogsKeySecretPath(organizationID)

	_, err := w.c.Logical().Write(secretPath, map[string]interface{}{
		"data": map[string]interface{}{
			"key": base64.StdEncoding.EncodeToString(key),
		},
	})
	if err != nil {
		return fmt.Errorf("could not write logs key to Vault: %w", err)
	}
	return nil
}

func (w *Wrapper) DeleteOrganizationLogsKeySecret(organizationID uuid.UUID) error {
	secretPath := w.createOrganizationLogsKeySecretPath(organizationID)
	_, err := w.c.Logical().Delete(secretPath)
	return err
}
