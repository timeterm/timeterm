package main

import (
	"fmt"

	"github.com/google/uuid"
)

func createDevUserConfig(id uuid.UUID) userConfig {
	return userConfig{
		streams: []streamConsumer{
			{
				stream:   "EMDEV-DISOWN-TOKEN",
				consumer: fmt.Sprintf("EMDEV-%s-EMDEV-DISOWN-TOKEN", id),
			},
		},
		other: []aclEntry{
			{
				topic: fmt.Sprintf("EMDEV.%s.REBOOT", id),
				op:    topicOpSub,
			},
		},
	}
}
