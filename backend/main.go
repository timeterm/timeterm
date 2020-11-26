package main

import (
	"context"
	"errors"
	"fmt"
	"io/ioutil"
	"net/http"
	"os"
	"os/signal"

	"github.com/go-logr/logr"
	"github.com/go-logr/zapr"
	_ "github.com/joho/godotenv/autoload"
	_ "github.com/lib/pq"
	"github.com/nats-io/jwt/v2"
	"github.com/nats-io/nats.go"
	"go.uber.org/zap"

	"gitlab.com/timeterm/timeterm/backend/api"
	"gitlab.com/timeterm/timeterm/backend/database"
	_ "gitlab.com/timeterm/timeterm/backend/pkg/natspb"
	"gitlab.com/timeterm/timeterm/backend/secrets"
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

	secr, err := secrets.New()
	if err != nil {
		log.Error(err, "could not create a secret wrapper")
		os.Exit(1)
	}

	nc, err := nats.Connect(os.Getenv("NATS_URL"), nats.UserJWT(natsCredsCBs()))
	if err != nil {
		return fmt.Errorf("could not connect to NATS: %w", err)
	}
	defer func() {
		if err = nc.Drain(); err != nil {
			log.Error(err, "could not drain NATS connection")
		}
	}()

	server, err := api.NewServer(db, log, nc, secr)
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

func wipeBytes(bs []byte) {
	for i := range bs {
		bs[i] = 'X'
	}
}

func natsCredsCBs() (nats.UserJWTHandler, nats.SignatureHandler) {
	getCreds := func() ([]byte, error) {
		endpoint := os.Getenv("NATS_GET_CREDS_ENDPOINT")
		if endpoint == "" {
			return nil, errors.New("environment variable NATS_GET_CREDS_ENDPOINT is not set")
		}

		rsp, err := http.Get(endpoint)
		if err != nil {
			return nil, fmt.Errorf("could not request NATS credentials: %w", err)
		}
		defer func() { _ = rsp.Body.Close() }()

		return ioutil.ReadAll(rsp.Body)
	}

	jwtCB := func() (string, error) {
		creds, err := getCreds()
		if err != nil {
			return "", err
		}
		defer wipeBytes(creds)
		return jwt.ParseDecoratedJWT(creds)
	}

	signCB := func(nonce []byte) ([]byte, error) {
		creds, err := getCreds()
		if err != nil {
			return nil, err
		}
		defer wipeBytes(creds)

		nkey, err := jwt.ParseDecoratedUserNKey(creds)
		if err != nil {
			return nil, err
		}
		defer nkey.Wipe()

		return nkey.Sign(nonce)
	}

	return jwtCB, signCB
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
