package main

import (
	"context"
	"errors"
	"fmt"
	"os"
	"os/signal"

	"github.com/go-logr/logr"
	"github.com/go-logr/zapr"
	_ "github.com/joho/godotenv/autoload"
	_ "github.com/lib/pq"
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

	secr, err := secrets.New(os.Getenv("VAULT_MOUNT"), os.Getenv("VAULT_PREFIX"))
	if err != nil {
		return fmt.Errorf("could not create secrets wrapper: %w", err)
	}

	server, err := api.NewServer(log, db, secr)
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
