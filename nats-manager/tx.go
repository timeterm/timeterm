package main

import (
	"github.com/nats-io/nats.go"
	"gitlab.com/timeterm/timeterm/backend/mq/natspb"
)

const (
	topicProvisionNewDevice = "provision-new-device"
)

type tx struct {
	enc *nats.EncodedConn
}

func newTx(nc *nats.Conn) error {
	enc := nats.EncodedConn{
		Conn: nc,
		Enc: natspb.NewEncoder(),
	}

	tx := tx{&enc}

	_, err := enc.QueueSubscribe(topicProvisionNewDevice, topicProvisionNewDevice, tx.provisionNewDevice)	
	return err
}

func (t *tx) provisionNewDevice(sub, reply string /* , msg *mqpb.NewDeviceMessage */) {
	
}
