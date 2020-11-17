package main

import (
	"fmt"

	"github.com/google/uuid"
	"github.com/nats-io/jsm.go"
	"github.com/nats-io/nats.go"
)

type handler struct {
	nc  *nats.Conn
	nsc *nsc
}

func (h *handler) provisionNewDevice(id uuid.UUID) (err error) {
	mgr, err := jsm.New(h.nc)
	if err != nil {
		return fmt.Errorf("could not create JetStream manager: %w", err)
	}

	err = setUpDeviceConsumers(id, mgr)
	if err != nil {
		return fmt.Errorf("could not set up device consumers: %w", err)
	}

	err = h.nsc.createNewDevUser(id)
	if err != nil {
		return fmt.Errorf("could not create new device user: %w", err)
	}
	return nil
}

func (h *handler) generateDeviceCredentials(id uuid.UUID) (creds string, err error) {
	accountName := fmt.Sprintf("EMDEV-%s", id)
	creds, err = h.nsc.generateUserCreds(accountName, accountName)
	if err != nil {
		return "", fmt.Errorf("could not generate credentials for device (user): %w", err)
	}
	return
}
