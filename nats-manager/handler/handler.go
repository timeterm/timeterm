package handler

import (
	"context"
	"fmt"
	"os"

	"github.com/google/uuid"
	"github.com/nats-io/jsm.go"
	"github.com/nats-io/nats.go"

	"gitlab.com/timeterm/timeterm/nats-manager/manager"
)

type Handler struct {
	nc        *nats.Conn
	mgr       *manager.Manager
	streamMgr *jsm.Manager
}

func New(ctx context.Context, nc *nats.Conn, mgr *manager.Manager) (*Handler, error) {
	nc, err := nats.Connect(os.Getenv("NATS_URL"),
		nats.UserJWT(mgr.NATSCredsCBs(ctx, "superuser", "EMDEVS")),
	)
	if err != nil {
		return nil, fmt.Errorf("could not connect to NATS: %w", err)
	}

	streamMgr, err := jsm.New(nc)
	if err != nil {
		return nil, fmt.Errorf("could not create JetStream manager: %w", err)
	}

	return &Handler{
		nc:        nc,
		mgr:       mgr,
		streamMgr: streamMgr,
	}, nil
}

func (h *Handler) Close() {
	h.nc.Close()
}

func (h *Handler) ProvisionNewDevice(ctx context.Context, id uuid.UUID) (err error) {
	err = setUpDeviceConsumers(id, h.streamMgr)
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
