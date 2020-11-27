package handler

import (
	"context"
	"fmt"

	"github.com/google/uuid"
	"github.com/nats-io/jsm.go"
	"github.com/nats-io/nats.go"

	"gitlab.com/timeterm/timeterm/nats-manager/manager"
)

type Handler struct {
	nc  *nats.Conn
	mgr *manager.Manager
}

func New(nc *nats.Conn, mgr *manager.Manager) *Handler {
	return &Handler{
		nc:  nc,
		mgr: mgr,
	}
}

func (h *Handler) ProvisionNewDevice(ctx context.Context, id uuid.UUID) (err error) {
	mgr, err := jsm.New(h.nc)
	if err != nil {
		return fmt.Errorf("could not create JetStream manager: %w", err)
	}

	err = setUpDeviceConsumers(id, mgr)
	if err != nil {
		return fmt.Errorf("could not set up device consumers: %w", err)
	}

	err = h.mgr.ProvisionNewDevice(ctx, id)
	if err != nil {
		return fmt.Errorf("could not provision new device: %w", err)
	}
	return nil
}

func (h *Handler) GenerateDeviceCredentials(ctx context.Context, id uuid.UUID) (creds []byte, err error) {
	creds, err = h.mgr.GenerateDeviceCredentials(ctx, id)
	if err != nil {
		return nil, fmt.Errorf("could not generate credentials for device (user): %w", err)
	}
	return
}
