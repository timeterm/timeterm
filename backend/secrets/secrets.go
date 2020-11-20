package secrets

import (
	"encoding/base64"
	"errors"
	"fmt"

	"github.com/google/uuid"
	"github.com/hashicorp/vault/api"
	"google.golang.org/protobuf/proto"

	devcfgpb "gitlab.com/timeterm/timeterm/proto/go/devcfg"
)

type Wrapper struct {
	c *api.Client
}

func New() (*Wrapper, error) {
	client, err := api.NewClient(api.DefaultConfig())
	if err != nil {
		return nil, err
	}
	return &Wrapper{client}, nil
}

func createNetworkingServiceConfigSecretPath(id uuid.UUID) string {
	return fmt.Sprintf("/timeterm/timeterm/ethernet/config/%s", id)
}

func (w *Wrapper) GetNetworkingServiceConfig(id uuid.UUID) (*devcfgpb.NetworkingService, error) {
	secretPath := createNetworkingServiceConfigSecretPath(id)
	secret, err := w.c.Logical().Read(secretPath)
	if err != nil {
		return nil, err
	}

	bytes, ok := secret.Data["config"].(string)
	if !ok {
		return nil, errors.New("could not retrieve config from secret")
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
	secretPath := createNetworkingServiceConfigSecretPath(id)
	_, err := w.c.Logical().Delete(secretPath)
	return err
}

func (w *Wrapper) UpsertEthernetConfig(id uuid.UUID, cfg *devcfgpb.NetworkingService) error {
	bytes, err := proto.Marshal(cfg)
	if err != nil {
		return err
	}

	encoded := base64.StdEncoding.EncodeToString(bytes)

	secretPath := createNetworkingServiceConfigSecretPath(id)

	_, err = w.c.Logical().Write(secretPath, map[string]interface{}{
		"config": encoded,
	})
	return err
}
