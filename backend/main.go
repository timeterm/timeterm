package main

import (
	"context"
	"errors"
	"os"
	"os/signal"

	"github.com/go-logr/logr"
	"github.com/go-logr/zapr"
	_ "github.com/joho/godotenv/autoload"
	_ "github.com/lib/pq"
	"github.com/nats-io/nats.go"
	"go.uber.org/zap"

	"gitlab.com/timeterm/timeterm/backend/api"
	"gitlab.com/timeterm/timeterm/backend/database"
	"gitlab.com/timeterm/timeterm/backend/mq"
	_ "gitlab.com/timeterm/timeterm/backend/mq/natspb"
)

func main() {
	logger, _ := zap.NewDevelopment()
	defer func() { _ = logger.Sync() }()

	log := zapr.NewLogger(logger)
	defer log.Info("shutdown complete")

	log.Info("starting")
	db, err := database.New(os.Getenv("DATABASE_URL"), log,
		database.WithJanitor(true),
	)
	if err != nil {
		log.Error(err, "could not open database")
		os.Exit(1)
	}
	defer func() {
		if err = db.Close(); err != nil {
			log.Error(err, "could not close database")
		}
	}()

	nc, err := nats.Connect(os.Getenv("NATS_URL"))
	if err != nil {
		log.Error(err, "could not connect to NATS")
		os.Exit(1)
	}

	mqw := mq.NewWrapper(nc)
	server, err := api.NewServer(db, log, mqw)
	if err != nil {
		log.Error(err, "could not create API server")
		os.Exit(1)
	}

	ctx, cancel := contextWithTermination(context.Background(), log)
	defer cancel()

	err = server.Run(ctx)
	if !errors.Is(err, context.Canceled) {
		log.Error(err, "error running API server")
		os.Exit(1)
	}
}

func contextWithTermination(ctx context.Context, log logr.Logger) (context.Context, func()) {
	ctx, cancel := context.WithCancel(ctx)

	go func() {
		sigs := make(chan os.Signal, 1)
		signal.Notify(sigs, os.Interrupt)
		defer signal.Stop(sigs)

		select {
		case <-sigs:
			log.Info("shutting down")

			cancel()
		case <-ctx.Done():
		}
	}()

	return ctx, cancel
}
