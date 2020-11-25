package static

import (
	"fmt"

	"github.com/go-logr/logr"
	"github.com/nats-io/jsm.go"
	"github.com/nats-io/nats.go"
)

func ConfigureStreams(log logr.Logger, nc *nats.Conn) error {
	log.Info("setting up static streams")

	mgr, err := jsm.New(nc)
	if err != nil {
		return fmt.Errorf("could not create JetStream manager: %w", err)
	}

	streams := map[string][]jsm.StreamOption{
		"EMDEV-DISOWN-TOKEN": {
			jsm.Subjects("EMDEV.*.DISOWN-TOKEN"),
		},
	}

	// We'll just assume that the stream doesn't exist if the LoadStream call errors
	// because at the time of writing jsm.go does not tell you reliably what the error actually is.
	for name, opts := range streams {
		log := log.WithValues("streamName", name)

		if strm, err := mgr.LoadStream(name); err != nil {
			log.Info("creating stream")

			if _, err = mgr.NewStream(name, opts...); err != nil {
				return fmt.Errorf("could not create JetStream stream: %w", err)
			}

			log.Info("created stream")
		} else {
			log.Info("updating stream")

			if err = strm.UpdateConfiguration(strm.Configuration(), opts...); err != nil {
				return fmt.Errorf("coudl not update JetStream stream configuration: %w", err)
			}

			log.Info("updated stream")
		}
	}

	log.Info("done setting up static streams")

	return nil
}
