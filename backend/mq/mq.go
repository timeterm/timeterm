package mq

import (
	"fmt"
	"os"

	"github.com/google/uuid"
	"github.com/labstack/gommon/log"
	"github.com/nats-io/nats.go"
	nmsdk "gitlab.com/timeterm/timeterm/nats-manager/sdk"
	mqpb "gitlab.com/timeterm/timeterm/proto/go/mq"

	"gitlab.com/timeterm/timeterm/backend/pkg/natspb"
)

type Wrapper struct {
	enc *nats.EncodedConn
}

func NewWrapper() (*Wrapper, error) {
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
	defer func() {
		if err = nc.Drain(); err != nil {
			log.Error(err, "could not drain NATS connection")
		}
	}()

	return &Wrapper{
		enc: &nats.EncodedConn{
			Conn: nc,
			Enc:  natspb.NewEncoder(),
		},
	}, nil
}

func (w *Wrapper) RebootDevice(id uuid.UUID) error {
	return w.enc.Publish(fmt.Sprintf("EMDEV.%s.REBOOT", id), &mqpb.RebootMessage{
		DeviceId: id.String(),
	})
}
