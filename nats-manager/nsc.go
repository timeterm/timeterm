package main

import (
	"fmt"
	"os/exec"

	"github.com/google/uuid"
)

func streamConsumerACLs(c streamConsumer) []string {
	// See https://github.com/nats-io/jetstream#acls
	return []string{
		fmt.Sprintf("$JS.API.CONSUMER.MSG.NEXT.%s.%s", c.stream, c.consumer),
		fmt.Sprintf("$JS.ACK.%s.%s.>", c.stream, c.consumer),
	}
}

func createNewDevUser(id uuid.UUID) (string, error) {
	accountName := fmt.Sprintf("fedev-%s", id)

	nscAddAccount(accountName)

	nscAddUser(accountName, []streamConsumer{
		{
			stream:   "EMDEV-DISOWN-TOKEN",
			consumer: fmt.Sprintf("FEDEV-%s", id),
		},
	})

	return "", nil
}

type streamConsumer struct {
	stream, consumer string
}

func nscAddAccount(name string) error {
	args := []string{"add", "account", "--name", name}

	return exec.Command("nsc", args...).Run()
}

func nscAddUser(name string, allowStreams []streamConsumer) error {
	args := []string{"add", "user", "--name", name}
	for _, sc := range allowStreams {
		for _, pubTopic := range streamConsumerACLs(sc) {
			args = append(args, "--allow-pub", pubTopic)
		}
	}

	return exec.Command("nsc", args...).Run()
}
