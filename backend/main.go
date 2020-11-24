package main

import (
	"context"
	"errors"
	"fmt"
	"io"
	"io/ioutil"
	"net/http"
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
	_ "gitlab.com/timeterm/timeterm/backend/pkg/natspb"
)

func main() {
	exitCode := 0
	defer os.Exit(exitCode)

	logger, _ := zap.NewDevelopment()
	defer func() { _ = logger.Sync() }()

	log := zapr.NewLogger(logger)
	defer log.Info("shutdown complete")

	if err := realMain(log); err != nil {
		log.Error(err, "error running backend")
		exitCode = 1
	}
}

func realMain(log logr.Logger) error {
	log.Info("starting")
	db, err := database.New(os.Getenv("DATABASE_URL"), log,
		database.WithJanitor(true),
	)
	if err != nil {
		return fmt.Errorf("could not open database: %w", err)
	}
	defer func() {
		if err = db.Close(); err != nil {
			log.Error(err, "could not close database")
		}
	}()

	log.Info("retrieving NATS credentials")
	credsFile, err := getNATSCreds()
	if err != nil {
		return  fmt.Errorf("error retrieving NATS credentials: %w", err)
	}
	log.Info("NATS credentials retrieved")

	nc, err := nats.Connect(os.Getenv("NATS_URL"), nats.UserCredentials(credsFile))
	if err != nil {
		return fmt.Errorf("could not connect to NATS: %w", err)
	}
	defer func() {
		if err = nc.Drain(); err != nil {
			log.Error(err, "could not drain NATS connection")
		}
	}()

	server, err := api.NewServer(db, log, nc)
	if err != nil {
		return fmt.Errorf("could not create API server: %w", err)
	}

	ctx, cancel := contextWithTermination(context.Background(), log)
	defer cancel()

	err = server.Run(ctx)
	if !errors.Is(err, context.Canceled) {
		return fmt.Errorf("error running API server")
	}
	return nil
}

func getNATSCreds() (string, error) {
	endpoint := os.Getenv("NATS_GET_CREDS_ENDPOINT")
	if endpoint == "" {
		return "", errors.New("environment variable NATS_GET_CREDS_ENDPOINT is not set")
	}

	rsp, err := http.Get(endpoint)
	if err != nil {
		return "", fmt.Errorf("could not request NATS credentials: %w", err)
	}
	defer func() { _ = rsp.Body.Close() }()

	f, err := ioutil.TempFile("", "*.creds")
	if err != nil {
		return "", fmt.Errorf("could not create temporary credentials file: %w", err)
	}
	defer func() { _ = f.Close() }()

	if _, err = io.Copy(f, rsp.Body); err != nil {
		return "", fmt.Errorf("could not copy response body to temporary file: %w", err)
	}

	return f.Name(), nil
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
