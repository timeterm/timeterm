package main_test

import (
	"context"
	"testing"

	"github.com/google/uuid"
	"github.com/nats-io/nats.go"
	"github.com/stretchr/testify/require"

	nmsdk "gitlab.com/timeterm/timeterm/nats-manager/pkg/sdk"
)

func TestNatsManager(t *testing.T) {
	nc, err := nats.Connect("nats://localhost:4222")
	require.NoError(t, err)

	deviceID := uuid.New()
	client := nmsdk.NewClient(nc)

	err = client.ProvisionNewDevice(context.Background(), deviceID)
	require.NoError(t, err)

	creds, err := client.GenerateDeviceCredentials(context.Background(), deviceID)
	require.NoError(t, err)

	t.Logf("Got NATS creds: \n%s\n", creds)
}
