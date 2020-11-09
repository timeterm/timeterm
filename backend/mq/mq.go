package mq

import (
	"fmt"

	"github.com/google/uuid"
	"github.com/nats-io/nats.go"
	mqpb "gitlab.com/timeterm/timeterm/proto/go/mq"

	"gitlab.com/timeterm/timeterm/backend/pkg/natspb"
)

type Wrapper struct {
	enc *nats.EncodedConn
}

func NewWrapper(nc *nats.Conn) *Wrapper {
	return &Wrapper{
		enc: &nats.EncodedConn{
			Conn: nc,
			Enc:  natspb.NewEncoder(),
		},
	}
}

func (w *Wrapper) RebootDevice(id uuid.UUID) error {
	return w.enc.Publish(fmt.Sprintf("FEDEV.%s.REBOOT", id), &mqpb.RebootMessage{
		DeviceId: id.String(),
	})
}
