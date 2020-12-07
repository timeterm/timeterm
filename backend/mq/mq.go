package mq

import (
	"fmt"
	"os"

	"github.com/go-logr/logr"
	"github.com/google/uuid"
	"github.com/nats-io/nats.go"
	nmsdk "gitlab.com/timeterm/timeterm/nats-manager/sdk"
	mqpb "gitlab.com/timeterm/timeterm/proto/go/mq"

	"gitlab.com/timeterm/timeterm/backend/pkg/natspb"
)

type Wrapper struct {
	log logr.Logger
	enc *nats.EncodedConn
}

func NewWrapper(log logr.Logger) (*Wrapper, error) {
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
	}, nil
}

func (w *Wrapper) RebootDevice(id uuid.UUID) error {
	log := w.log.WithValues("deviceId", id)

	subj := fmt.Sprintf("EMDEV.%s.REBOOT", id)
	log.V(1).WithValues("subject", subj).Info("publishing reboot message")
	err := w.enc.Publish(subj, &mqpb.RebootMessage{
		DeviceId: id.String(),
	})
	if err != nil {
		log.V(1).WithValues("subject", subj).Error(err, "publishing failed")
	} else {
		log.V(1).WithValues("subject", subj).Info("publishing succeeded")
	}

	return err
}
