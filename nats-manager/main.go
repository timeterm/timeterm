package main

import (
	"log"
	"os"

	"github.com/nats-io/nats.go"
	"gitlab.com/timeterm/timeterm/backend/pkg/natspb"
)

func main() {
	nc, err := nats.Connect(os.Getenv("NATS_URL"))
	if err != nil {
		log.Fatalf("Could not connect to NATS: %v", err)
	}

	enc := nats.EncodedConn{
		Conn: nc,
		Enc:  natspb.NewEncoder(),
	}

	enc.QueueSubscribe("FEDEV.new", "FEDEV.new", func(subject, reply string /* o *mqpb.NewDevMessage */) {
		enc.Publish(reply, nil)
	})
}
