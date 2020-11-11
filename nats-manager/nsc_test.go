package main

import (
	"fmt"
	"testing"

	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
)

func TestNscAddUserCmd(t *testing.T) {
	devID := uuid.New()
	devUser := fmt.Sprintf("EMDEV-%s", devID)

	cmd := nscAddUserCmd(devUser, devUser, createDevUserConfig(devID))

	assert.Equal(t, []string{
		"nsc", "add", "user",
		"--name", devUser,
		"--account", devUser,
		"--allow-sub", fmt.Sprintf("EMDEV.%s.REBOOT", devID),
		"--allow-pub", fmt.Sprintf("$JS.API.CONSUMER.MSG.NEXT.EMDEV-DISOWN-TOKEN.EMDEV-%s-EMDEV-DISOWN-TOKEN", devID),
		"--allow-pub", fmt.Sprintf("$JS.ACK.EMDEV-DISOWN-TOKEN.EMDEV-%s-EMDEV-DISOWN-TOKEN.>", devID),
	}, cmd.Args)
}
