package nmsdk

import (
	"context"
	"errors"
	"fmt"

	"github.com/google/uuid"
	"github.com/nats-io/nats.go"
	"gitlab.com/timeterm/timeterm/backend/pkg/natspb"
	rpcpb "gitlab.com/timeterm/timeterm/proto/go/rpc"
)

const (
	SubjectProvisionNewDevice        = "NATS-MANAGER.PROVISION-NEW-DEVICE"
	SubjectGenerateDeviceCredentials = "NATS-MANAGER.GENERATE-DEVICE-CREDENTIALS"
)

type ManagerError struct {
	Data *rpcpb.Error
}

func (e ManagerError) Error() string {
	return fmt.Sprintf("nats-manager server error: %s", e.Data.GetMessage())
}

type Client struct {
	enc *nats.EncodedConn
}

func NewClient(nc *nats.Conn) *Client {
	return &Client{
		enc: &nats.EncodedConn{
			Enc:  natspb.NewEncoder(),
			Conn: nc,
		},
	}
}

func (c *Client) ProvisionNewDevice(ctx context.Context, id uuid.UUID) error {
	var rsp rpcpb.ProvisionNewDeviceResponse

	err := c.enc.RequestWithContext(ctx, SubjectProvisionNewDevice, &rpcpb.ProvisionNewDeviceRequest{
		DeviceId: id.String(),
	}, &rsp)
	if err != nil {
		return err
	}

	switch data := rsp.Response.(type) {
	case *rpcpb.ProvisionNewDeviceResponse_Success:
		return nil
	case *rpcpb.ProvisionNewDeviceResponse_Error:
		return ManagerError{data.Error}
	default:
		return errors.New("invalid response")
	}
}

func (c *Client) GenerateDeviceCredentials(ctx context.Context, id uuid.UUID) (string, error) {
	var rsp rpcpb.GenerateDeviceCredentialsResponse

	err := c.enc.RequestWithContext(ctx, SubjectGenerateDeviceCredentials, &rpcpb.GenerateDeviceCredentialsRequest{
		DeviceId: id.String(),
	}, &rsp)
	if err != nil {
		return "", err
	}

	switch data := rsp.Response.(type) {
	case *rpcpb.GenerateDeviceCredentialsResponse_Sucess:
		return data.Sucess.GetNatsCreds(), nil
	case *rpcpb.GenerateDeviceCredentialsResponse_Error:
		return "", ManagerError{data.Error}
	default:
		return "", errors.New("invalid response")
	}
}
