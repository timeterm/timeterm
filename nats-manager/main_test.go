package main_test

import (
	"context"
	"testing"

	"github.com/google/uuid"
	"github.com/nats-io/nats.go"
	"github.com/stretchr/testify/require"
	rpcpb "gitlab.com/timeterm/timeterm/proto/go/rpc"

	nmsdk "gitlab.com/timeterm/timeterm/nats-manager/sdk"
)

func TestNatsManager(t *testing.T) {
	nc, err := nats.Connect("nats://localhost:4222")
	require.NoError(t, err)

	deviceID := uuid.New()
	client := nmsdk.NewClient(nc)

	data, err := client.ProvisionNewDevice(context.Background(), &rpcpb.ProvisionNewDeviceRequest{
		DeviceId: deviceID.String(),
	})
	require.NoError(t, err)

	t.Logf("Got NATS creds: \n%s\n", data.GetNatsCreds())
}
