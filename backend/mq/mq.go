package mq

import (
	"context"
	"fmt"
	"os"
	"sync"
	"time"

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

	debounces sync.Map
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

func (w *Wrapper) NetworkingConfigUpdated(organizationID uuid.UUID) {
	w.GetNetworkConfigUpdatedDebounce(organizationID)()
}

func (w *Wrapper) GetNetworkConfigUpdatedDebounce(organizationID uuid.UUID) func() {
	if d, ok := w.debounces.Load(organizationID); ok {
		return d.(func())
	}

	ctx, cancel := context.WithCancel(context.Background())
	time.AfterFunc(30*time.Second, func() {
		cancel()
		w.debounces.Delete(organizationID)
	})

	debfn := debounce(ctx, func() {
		log := w.log.WithValues("organizationId", organizationID)

		if err := w.dbw.WalkDevices(context.Background(), organizationID, func(d *database.Device) bool {
			log = log.WithValues("deviceId", d.ID)
			if err := w.RetrieveNewNetworkingConfig(d.ID); err != nil {
				log.Error(err, "could not send message to device to retrieve new networking config")
			}
			return true
		}); err != nil {
			log.Error(err, "could not walk devices in organization (to send RetrieveNewNetworkingConfig messages)")
		}
	}, time.Second)

	d, _ := w.debounces.LoadOrStore(organizationID, debfn)
	return d.(func())
}

func debounce(ctx context.Context, f func(), d time.Duration) func() {
	var mu sync.Mutex
	var t *time.Timer

	go func() {
		<-ctx.Done()
		mu.Lock()
		defer mu.Unlock()

		if t != nil && t.Stop() {
			go f()
		}
	}()

	return func() {
		select {
		case <-ctx.Done():
			return
		default:
		}

		mu.Lock()
		defer mu.Unlock()

		if t != nil {
			t.Stop()
		}
		t = time.AfterFunc(d, f)
	}
}
