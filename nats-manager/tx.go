package main

import (
	"context"
	"errors"
	"fmt"
	"runtime/debug"

	"github.com/go-logr/logr"
	"github.com/golang/protobuf/ptypes/empty"
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

	_, err := enc.QueueSubscribe(
		nmsdk.SubjectProvisionNewDevice,
		nmsdk.SubjectProvisionNewDevice,
		tx.handleProvisionNewDevice,
	)
	if err != nil {
		return err
	}

	_, err = enc.QueueSubscribe(
		nmsdk.SubjectGenerateDeviceCredentials,
		nmsdk.SubjectGenerateDeviceCredentials,
		tx.handleGenerateDeviceCredentials,
	)
	if err != nil {
		return err
	}

	err = nc.Flush()
	if err != nil {
		return err
	}

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

	err := t.provisionNewDevice(msg)
	if err != nil {
		rsp.Response = &rpcpb.ProvisionNewDeviceResponse_Error{
			Error: &rpcpb.Error{
				Message: err.Error(),
			},
		}
	} else {
		rsp.Response = &rpcpb.ProvisionNewDeviceResponse_Success{Success: &empty.Empty{}}
	}

	err = t.enc.Publish(reply, rsp)
	if err != nil {
		t.log.Error(err, "could not publish provisionNewDevice response")
	}
}

func (t *tx) handleGenerateDeviceCredentials(_ /* sub */, reply string, msg *rpcpb.GenerateDeviceCredentialsRequest) {
	defer t.handlePanic()

	rsp := new(rpcpb.GenerateDeviceCredentialsResponse)

	creds, err := t.generateDeviceCredentials(msg)
	if err != nil {
		rsp.Response = &rpcpb.GenerateDeviceCredentialsResponse_Error{
			Error: &rpcpb.Error{
				Message: err.Error(),
			},
		}
	} else {
		rsp.Response = &rpcpb.GenerateDeviceCredentialsResponse_Sucess{
			Sucess: &rpcpb.DeviceCredentials{
				NatsCreds: creds,
			},
		}
	}

	err = t.enc.Publish(reply, rsp)
	if err != nil {
		t.log.Error(err, "could not publish generateDeviceCredentials response")
	}
}

func (t *tx) provisionNewDevice(msg *rpcpb.ProvisionNewDeviceRequest) error {
	devID, err := uuid.Parse(msg.GetDeviceId())
	if err != nil {
		return errors.New("invalid device ID")
	}

	err = t.h.provisionNewDevice(devID)
	if err != nil {
		t.log.Error(err, "could not provision new device", errorLogArgs(err)...)

		return fmt.Errorf("could not provision new device: %w", err)
	}

	return nil
}

func (t *tx) generateDeviceCredentials(msg *rpcpb.GenerateDeviceCredentialsRequest) (string, error) {
	devID, err := uuid.Parse(msg.GetDeviceId())
	if err != nil {
		return "", errors.New("invalid device ID")
	}

	creds, err := t.h.generateDeviceCredentials(devID)
	if err != nil {
		t.log.Error(err, "could not generate device credentials", errorLogArgs(err)...)

		return "", fmt.Errorf("could not generate credentials for device: %w", err)
	}

	return creds, nil
}
