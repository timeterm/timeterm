package main

import (
	"context"
	"os"

	"github.com/go-logr/zapr"
	_ "github.com/joho/godotenv/autoload"
	_ "github.com/lib/pq"
	"github.com/nats-io/nats.go"
	"go.uber.org/zap"

	"gitlab.com/timeterm/timeterm/backend/api"
	"gitlab.com/timeterm/timeterm/backend/database"
)

func main() {
	logger, _ := zap.NewDevelopment()
	defer func() { _ = logger.Sync() }()

	log := zapr.NewLogger(logger)

	log.Info("starting")
	db, err := database.New(os.Getenv("DATABASE_URL"), log,
		database.WithJanitor(true),
	)
	if err != nil {
		log.Error(err, "could not open database")
		os.Exit(1)
	}

	nc, err := nats.Connect(os.Getenv("NATS_URL"))
	if err != nil {
		log.Error(err, "could not connect to NATS")
		os.Exit(1)
	}

	server, err := api.NewServer(db, log, nc)
	if err != nil {
		log.Error(err, "could not create API server")
		os.Exit(1)
	}

	err = server.Run(context.Background())
	log.Error(err, "error running API server")
	os.Exit(1)
}
