package mq

import (
	"context"
	"fmt"
	"os"

	"github.com/go-logr/logr"
	"github.com/google/uuid"
	"github.com/nats-io/nats.go"
	nmsdk "gitlab.com/timeterm/timeterm/nats-manager/sdk"
	mqpb "gitlab.com/timeterm/timeterm/proto/go/mq"

	"gitlab.com/timeterm/timeterm/backend/database"
	"gitlab.com/timeterm/timeterm/backend/pkg/natspb"
)

type Wrapper struct {
	log logr.Logger
	enc *nats.EncodedConn
	dbw *database.Wrapper
}

func NewWrapper(log logr.Logger, dbw *database.Wrapper) (*Wrapper, error) {
	acr, err := nmsdk.NewAppCredsRetrieverFromEnv("backend-emdevs")
	if err != nil {
		return nil, fmt.Errorf("could not create (NATS) app credentials retriever: %w", err)
	}

	nc, err := nats.Connect(os.Getenv("NATS_URL"),
		nats.UserJWT(acr.NatsCredsCBs()),
		// Never stop trying to reconnect.
		nats.MaxReconnects(-1),
	)
	if err != nil {
		return nil, fmt.Errorf("could not connect to NATS: %w", err)
	}

	return &Wrapper{
		log: log.WithName("MqWrapper"),
		enc: &nats.EncodedConn{
			Conn: nc,
			Enc:  natspb.NewEncoder(),
		},
		dbw: dbw,
	}, nil
}

func (w *Wrapper) RebootDevice(id uuid.UUID) error {
	log := w.log.WithValues("deviceId", id)

	subj := fmt.Sprintf("EMDEV.%s.REBOOT", id)
	log = log.V(1).WithValues("subject", subj)
	log.Info("publishing reboot message")
	err := w.enc.Publish(subj, new(mqpb.RebootMessage))
	if err != nil {
		log.Error(err, "publishing failed")
	} else {
		log.Info("publishing succeeded")
	}

	return err
}

func (w *Wrapper) NetworkingConfigUpdated(organizationID uuid.UUID) {
	log := w.log.WithValues("organizationId", organizationID)

	go func() {
		if err := w.dbw.WalkDevices(context.Background(), organizationID, func(d *database.Device) bool {
			log = log.WithValues("deviceId", d.ID)
			if err := w.RetrieveNewNetworkingConfig(d.ID); err != nil {
				log.Error(err, "could not send message to device to retrieve new networking config")
			}
			return true
		}); err != nil {
			log.Error(err, "could not walk devices in organization (to send RetrieveNewNetworkingConfig messages)")
		}
	}()
}

func (w *Wrapper) RetrieveNewNetworkingConfig(deviceID uuid.UUID) error {
	log := w.log.WithValues("deviceId", deviceID)

	subj := fmt.Sprintf("EMDEV.%s.RETRIEVE-NEW-NETWORKING-CONFIG", deviceID)
	log = log.V(1).WithValues("subject", subj)
	log.Info("publishing new networking config retrieval message")
	err := w.enc.Publish(subj, new(mqpb.RebootMessage))
	if err != nil {
		log.Error(err, "publishing failed")
	} else {
		log.Info("publishing succeeded")
	}

	return err
}
