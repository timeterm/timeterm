package main

import (
	"errors"
	"fmt"

	"github.com/go-logr/logr"
	"github.com/google/uuid"
	"github.com/nats-io/nats.go"
	"gitlab.com/timeterm/timeterm/backend/pkg/natspb"
	rpcpb "gitlab.com/timeterm/timeterm/proto/go/rpc"
)

const (
	topicProvisionNewDevice = "provision-new-device"
)

type tx struct {
	enc *nats.EncodedConn
	log logr.Logger
	h   *handler
}

func newTx(nc *nats.Conn, log logr.Logger, h *handler) error {
	enc := nats.EncodedConn{
		Conn: nc,
		Enc:  natspb.NewEncoder(),
	}

	tx := tx{enc: &enc, log: log, h: h}

	_, err := enc.QueueSubscribe(topicProvisionNewDevice, topicProvisionNewDevice, tx.provisionNewDevice)
	return err
}

func (t *tx) handleProvisionNewDevice(_ /* sub */, reply string, msg *rpcpb.ProvisionNewDeviceRequest) {
	rsp := new(rpcpb.ProvisionNewDeviceResponse)

	data, err := t.provisionNewDevice(msg)
	if err != nil {
		rsp.Response = &rpcpb.ProvisionNewDeviceResponse_Error{
			Error: &rpcpb.Error{
				Message: err.Error(),
			},
		}
	} else {
		rsp.Response = &rpcpb.ProvisionNewDeviceResponse_Success{Success: data}
	}

	err = t.enc.Publish(reply, rsp)
	if err != nil {
		t.log.Error(err, "could not provisionNewDevice response")
	}
}

func (t *tx) provisionNewDevice(msg *rpcpb.ProvisionNewDeviceRequest) (*rpcpb.ProvisionNewDeviceResponseData, error) {
	devID, err := uuid.Parse(msg.GetDeviceId())
	if err != nil {
		return nil, errors.New("invalid device ID")
	}

	natsCreds, err := t.h.provisionNewDevice(devID)
	if err != nil {
		return nil, fmt.Errorf("could not provision new device: %w", err)
	}

	return &rpcpb.ProvisionNewDeviceResponseData{
		NatsCreds: natsCreds,
	}, nil
}
