package main

import (
	"fmt"

	"github.com/google/uuid"
	"github.com/nats-io/jsm.go"
	"github.com/nats-io/nats.go"
)

type handler struct {
	nc      *nats.Conn
	dataDir string
}

func (h *handler) provisionNewDevice(id uuid.UUID) (natsCreds string, err error) {
	mgr, err := jsm.New(h.nc)
	if err != nil {
		return "", fmt.Errorf("could not create JetStream manager: %w", err)
	}

	err = setUpDeviceConsumers(id, mgr)
	if err != nil {
		return "", fmt.Errorf("could not set up device consumers: %w", err)
	}

	natsCreds, err = createNewDevUser(id, &nscConfig{
		dataDir: h.dataDir,
	})
	if err != nil {
		return "", fmt.Errorf("could not create new device user: %w", err)
	}
	return natsCreds, nil
}
