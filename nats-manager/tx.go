package main

import (
	"context"
	"errors"
	"fmt"
	"runtime/debug"

	"github.com/go-logr/logr"
	"github.com/google/uuid"
	"github.com/nats-io/nats.go"
	"gitlab.com/timeterm/timeterm/backend/pkg/natspb"
	rpcpb "gitlab.com/timeterm/timeterm/proto/go/rpc"

	nmsdk "gitlab.com/timeterm/timeterm/nats-manager/sdk"
)

type tx struct {
	enc *nats.EncodedConn
	log logr.Logger
	h   *handler
}

func runTx(ctx context.Context, nc *nats.Conn, log logr.Logger, h *handler) error {
	enc := nats.EncodedConn{
		Conn: nc,
		Enc:  natspb.NewEncoder(),
	}

	tx := tx{enc: &enc, log: log, h: h}

	sub, err := enc.QueueSubscribe(
		nmsdk.SubjectProvisionNewDevice,
		nmsdk.SubjectProvisionNewDevice,
		tx.handleProvisionNewDevice,
	)
	if err != nil {
		return err
	}
	defer func() {
		err = sub.Drain()
		if err != nil {
			log.Error(err, "error draining subscription", "topic", sub.Subject)
		}
	}()

	<-ctx.Done()

	return nil
}

func (t *tx) handlePanic() {
	if r := recover(); r != nil {
		err, ok := r.(error)
		if !ok {
			err = nil
		}

		args := []interface{}{"stack", string(debug.Stack())}
		if err == nil {
			args = append(args, "r", r)
		}
		t.log.Error(err, "recovered from a panic", args...)
	}
}

func (t *tx) handleProvisionNewDevice(_ /* sub */, reply string, msg *rpcpb.ProvisionNewDeviceRequest) {
	defer t.handlePanic()

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
		var logArgs []interface{}
		var nerr nscError
		if errors.As(err, &nerr) {
			logArgs = append(logArgs, "log", nerr.Log())
		}
		t.log.Error(err, "could not provision new device", logArgs...)

		return nil, fmt.Errorf("could not provision new device: %w", err)
	}

	return &rpcpb.ProvisionNewDeviceResponseData{
		NatsCreds: natsCreds,
	}, nil
}
