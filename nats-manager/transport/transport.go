package transport

import (
	"context"
	"errors"
	"fmt"
	"runtime/debug"
	"time"

	"github.com/go-logr/logr"
	"github.com/google/uuid"
	"github.com/nats-io/nats.go"
	"gitlab.com/timeterm/timeterm/backend/pkg/natspb"
	rpcpb "gitlab.com/timeterm/timeterm/proto/go/rpc"

	"gitlab.com/timeterm/timeterm/nats-manager/handler"
	nmsdk "gitlab.com/timeterm/timeterm/nats-manager/sdk"
)

const requestHandleTimeout = time.Second * 30

type Transport struct {
	enc *nats.EncodedConn
	log logr.Logger
	h   *handler.Handler
}

func New(nc *nats.Conn, log logr.Logger, h *handler.Handler) *Transport {
	return &Transport{
		enc: &nats.EncodedConn{
			Conn: nc,
			Enc:  natspb.NewEncoder(),
		},
		log: log,
		h:   h,
	}
}

func (t *Transport) Run(ctx context.Context) error {
	if _, err := t.enc.QueueSubscribe(
		nmsdk.SubjectProvisionNewDevice,
		nmsdk.SubjectProvisionNewDevice,
		t.handleProvisionNewDevice,
	); err != nil {
		return err
	}

	if _, err := t.enc.QueueSubscribe(
		nmsdk.SubjectGenerateDeviceCredentials,
		nmsdk.SubjectGenerateDeviceCredentials,
		t.handleGenerateDeviceCredentials,
	); err != nil {
		return err
	}

	if err := t.enc.Flush(); err != nil {
		return err
	}

	<-ctx.Done()
	return ctx.Err()
}

func (t *Transport) handlePanic() {
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

func (t *Transport) handleProvisionNewDevice(_ /* sub */, reply string, msg *rpcpb.ProvisionNewDeviceRequest) {
	defer t.handlePanic()

	rsp := new(rpcpb.ProvisionNewDeviceResponse)
	ctx, cancel := context.WithTimeout(context.Background(), requestHandleTimeout)
	defer cancel()

	err := t.provisionNewDevice(ctx, msg)
	if err != nil {
		rsp.Response = &rpcpb.ProvisionNewDeviceResponse_Error{
			Error: &rpcpb.Error{
				Message: err.Error(),
			},
		}
	} else {
		rsp.Response = &rpcpb.ProvisionNewDeviceResponse_Success{Success: &rpcpb.Empty{}}
	}

	err = t.enc.Publish(reply, rsp)
	if err != nil {
		t.log.Error(err, "could not publish provisionNewDevice response")
	}
}

func (t *Transport) handleGenerateDeviceCredentials(
	_, /* sub */
	reply string,
	msg *rpcpb.GenerateDeviceCredentialsRequest,
) {
	defer t.handlePanic()

	rsp := new(rpcpb.GenerateDeviceCredentialsResponse)
	ctx, cancel := context.WithTimeout(context.Background(), requestHandleTimeout)
	defer cancel()

	creds, err := t.generateDeviceCredentials(ctx, msg)
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

func (t *Transport) provisionNewDevice(ctx context.Context, msg *rpcpb.ProvisionNewDeviceRequest) error {
	devID, err := uuid.Parse(msg.GetDeviceId())
	if err != nil {
		return errors.New("invalid device ID")
	}

	err = t.h.ProvisionNewDevice(ctx, devID)
	if err != nil {
		t.log.Error(err, "could not provision new device")

		return fmt.Errorf("could not provision new device: %w", err)
	}

	return nil
}

func (t *Transport) generateDeviceCredentials(
	ctx context.Context,
	msg *rpcpb.GenerateDeviceCredentialsRequest,
) (string, error) {
	devID, err := uuid.Parse(msg.GetDeviceId())
	if err != nil {
		return "", errors.New("invalid device ID")
	}

	creds, err := t.h.GenerateDeviceCredentials(ctx, devID)
	if err != nil {
		t.log.Error(err, "could not generate device credentials")

		return "", fmt.Errorf("could not generate credentials for device: %w", err)
	}

	return creds, nil
}
