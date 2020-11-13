package main_test

import (
	"testing"
	"time"

	"github.com/google/uuid"
	"github.com/nats-io/nats.go"
	"github.com/stretchr/testify/require"
	"gitlab.com/timeterm/timeterm/backend/pkg/natspb"
	rpcpb "gitlab.com/timeterm/timeterm/proto/go/rpc"
)

func TestNatsManager(t *testing.T) {
	nc, err := nats.Connect("nats://localhost:4222")
	require.NoError(t, err)

	enc := &nats.EncodedConn{
		Conn: nc,
		Enc:  natspb.NewEncoder(),
	}

	deviceID := uuid.New()
	var rsp rpcpb.ProvisionNewDeviceResponse

	err = enc.Request("NATS-MANAGER.PROVISION-NEW-DEVICE", &rpcpb.ProvisionNewDeviceRequest{
		DeviceId: deviceID.String(),
	}, &rsp, time.Second*3000)
	require.NoError(t, err)

	switch data := rsp.Response.(type) {
	case *rpcpb.ProvisionNewDeviceResponse_Error:
		t.Fatalf("Got error response: %v", data.Error.GetMessage())
	case *rpcpb.ProvisionNewDeviceResponse_Success:
		t.Logf("Got NATS creds: \n%s\n", data.Success.GetNatsCreds())
	}
}
