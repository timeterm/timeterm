package static

import (
	"context"
	"fmt"
	"os"

	"github.com/go-logr/logr"
	"github.com/nats-io/jsm.go"
	"github.com/nats-io/nats.go"

	"gitlab.com/timeterm/timeterm/nats-manager/manager"
)

type streams map[string]accountStreams
type accountStreams map[string]streamConfig

type streamConfig struct {
	options []jsm.StreamOption
}

func ConfigureStreams(ctx context.Context, log logr.Logger, mgr *manager.Manager) error {
	log.Info("setting up static streams")

	streams := streams{
		"EMDEVS": {
			"EMDEV-DISOWN-TOKEN": {
				options: []jsm.StreamOption{
					jsm.FileStorage(),
					jsm.Subjects("EMDEV.*.DISOWN-TOKEN"),
				},
			},
		},
	}

	// We'll just assume that the stream doesn't exist if the LoadStream call errors
	// because at the time of writing jsm.go does not tell you reliably what the error actually is.
	for accountName, streams := range streams {
		nc, err := nats.Connect(os.Getenv("NATS_URL"),
			nats.UserJWT(mgr.NATSCredsCBs(ctx, "superuser", accountName)),
		)
		if err != nil {
			return fmt.Errorf("can not configure streams for account %s: could not connect to NATS: %w",
				accountName, err,
			)
		}

		mgr, err := jsm.New(nc)
		if err != nil {
			return fmt.Errorf("can not configure streams for account %s: could not create NATS manager: %w",
				accountName, err,
			)
		}

		for streamName, stream := range streams {
			log := log.WithValues("streamName", streamName, "accountName", accountName)

			if strm, err := mgr.LoadStream(streamName); err != nil {
				log.Info("creating stream")

				if _, err = mgr.NewStream(streamName, stream.options...); err != nil {
					return fmt.Errorf("could not create JetStream stream: %w", err)
				}

				log.Info("created stream")
			} else {
				log.Info("updating stream")

				if err = strm.UpdateConfiguration(strm.Configuration(), stream.options...); err != nil {
					return fmt.Errorf("coudl not update JetStream stream configuration: %w", err)
				}

				log.Info("updated stream")
			}
		}
	}

	log.Info("done setting up static streams")

	return nil
}
