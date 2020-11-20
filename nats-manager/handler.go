package main

import (
	"context"
	"fmt"

	"github.com/google/uuid"
	"github.com/nats-io/jsm.go"
	"github.com/nats-io/nats.go"

	"gitlab.com/timeterm/timeterm/nats-manager/secrets"
)

type handler struct {
	nc *nats.Conn
	mg *secrets.Manager
}

func (h *handler) provisionNewDevice(ctx context.Context, id uuid.UUID) (err error) {
	mgr, err := jsm.New(h.nc)
	if err != nil {
		return fmt.Errorf("could not create JetStream manager: %w", err)
	}

	err = setUpDeviceConsumers(id, mgr)
	if err != nil {
		return fmt.Errorf("could not set up device consumers: %w", err)
	}

	err = h.mg.CreateNewDeviceUser(ctx, id)
	if err != nil {
		return fmt.Errorf("could not create new device user: %w", err)
	}
	return nil
}

func (h *handler) generateDeviceCredentials(ctx context.Context, id uuid.UUID) (creds string, err error) {
	creds, err = h.mg.GenerateDeviceCredentials(ctx, id)
	if err != nil {
		return "", fmt.Errorf("could not generate credentials for device (user): %w", err)
	}
	return
}
