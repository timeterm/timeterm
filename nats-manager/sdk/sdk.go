package nmsdk

import (
	"context"
	"errors"
	"fmt"

	"github.com/nats-io/nats.go"
	"gitlab.com/timeterm/timeterm/backend/pkg/natspb"
	rpcpb "gitlab.com/timeterm/timeterm/proto/go/rpc"
)

const SubjectProvisionNewDevice = "NATS-MANAGER.PROVISION-NEW-DEVICE"

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

func (c *Client) ProvisionNewDevice(ctx context.Context,
	req *rpcpb.ProvisionNewDeviceRequest,
) (*rpcpb.ProvisionNewDeviceResponseData, error) {
	var rsp rpcpb.ProvisionNewDeviceResponse

	err := c.enc.RequestWithContext(ctx, SubjectProvisionNewDevice, req, &rsp)
	if err != nil {
		return nil, err
	}

	switch data := rsp.Response.(type) {
	case *rpcpb.ProvisionNewDeviceResponse_Success:
		return data.Success, nil
	case *rpcpb.ProvisionNewDeviceResponse_Error:
		return nil, ManagerError{data.Error}
	default:
		return nil, errors.New("invalid response")
	}
}
