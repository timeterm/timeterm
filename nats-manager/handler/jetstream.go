package handler

import (
	"fmt"
	"time"

	"github.com/google/uuid"
	"github.com/nats-io/jsm.go"
)

func setUpDeviceConsumers(devID uuid.UUID, mgr *jsm.Manager) error {
	consumerName := fmt.Sprintf("EMDEV-%s-EMDEV-RETRIEVE-NEW-NETWORKING-CONFIG", devID)
	wantDisownTokenSubject := fmt.Sprintf("EMDEV.%s.RETRIEVE-NEW-NETWORKING-CONFIG", devID)

	_, err := mgr.NewConsumer("EMDEV-RETRIEVE-NEW-NETWORKING-CONFIG",
		jsm.DurableName(consumerName),
		jsm.FilterStreamBySubject(wantDisownTokenSubject),
		jsm.AckWait(time.Second*30),
		jsm.AcknowledgeExplicit(),
		jsm.DeliverAllAvailable(),
	)
	if err != nil {
		return fmt.Errorf("could not set up EMDEV-RETRIEVE-NEW-NETWORKING-CONFIG consumer: %w", err)
	}
	return nil
}
