package main

import (
	"fmt"
	"time"

	"github.com/google/uuid"
	"github.com/nats-io/jsm.go"
	"github.com/nats-io/nats.go"
)

type handler struct {
	nc *nats.Conn
}

func (h *handler) provisionNewDevice(id uuid.UUID) error {
	mgr, err := jsm.New(h.nc)
	if err != nil {
		return err
	}

	consumerName := fmt.Sprintf("EMDEV-%s", id)
	wantDisownTokenSubject := fmt.Sprintf("EMDEV.%s.DISOWN-TOKEN", id)
	_, err = mgr.NewConsumer("EMDEV-DISOWN-TOKEN",
		jsm.DurableName(consumerName),
		jsm.FilterStreamBySubject(wantDisownTokenSubject),
		jsm.AckWait(time.Second*30),
		jsm.AcknowledgeExplicit(),
		jsm.DeliverAllAvailable(),
	)
	if err != nil {
		return err
	}

	return nil
}
