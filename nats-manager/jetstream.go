package main

import (
	"fmt"
	"time"

	"github.com/google/uuid"
	"github.com/nats-io/jsm.go"
)

func setUpDeviceConsumers(devID uuid.UUID, mgr *jsm.Manager) error {
	consumerName := fmt.Sprintf("EMDEV-%s-EMDEV-DISOWN-TOKEN", devID)
	wantDisownTokenSubject := fmt.Sprintf("EMDEV.%s.DISOWN-TOKEN", devID)

	_, err := mgr.NewConsumer("EMDEV-DISOWN-TOKEN",
		jsm.DurableName(consumerName),
		jsm.FilterStreamBySubject(wantDisownTokenSubject),
		jsm.AckWait(time.Second*30),
		jsm.AcknowledgeExplicit(),
		jsm.DeliverAllAvailable(),
	)
	if err != nil {
		return fmt.Errorf("could not set up EMDEV-DISOWN-TOKEN consumer: %w", err)
	}
	return nil
}
