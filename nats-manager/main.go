package main

import (
	"context"
	"errors"
	"fmt"
	"io/ioutil"
	"os"
	"os/signal"
	"path"
	"syscall"

	"github.com/go-logr/logr"
	"github.com/go-logr/zapr"
	"github.com/nats-io/nats.go"
	"go.uber.org/zap"
)

func main() {
	exitCode := 0
	defer os.Exit(exitCode)

	logger, _ := zap.NewDevelopment()
	defer func() { _ = logger.Sync() }()

	log := zapr.NewLogger(logger)
	defer log.Info("shutdown complete")

	if err := realMain(log); err != nil {
		log.Error(err, "error running nats-manager")
		exitCode = 1
	}
}

func realMain(log logr.Logger) error {
	cfg, err := loadConfig()
	if err != nil {
		return fmt.Errorf("could not load configuration: %w", err)
	}

	nsc, err := newNsc(log, cfg.dataDir, cfg.nscPath)
	if err != nil {
		return err
	}

	log.Info("connecting with NATS", "url", cfg.natsURL)
	nc, err := nats.Connect(cfg.natsURL,
		nats.UserCredentials(path.Join(nsc.nkeysPath(), "creds/TIMETERM/NATS-MANAGER/NATS-MANAGER.creds")),
	)
	if err != nil {
		return fmt.Errorf("could not connect to NATS: %w", err)
	}
	defer func() {
		err = nc.Drain()
		if err != nil {
			log.Error(err, "error draining NATS connection on shutdown")
		}
	}()
	log.Info("connected to NATS")

	ctx, cancel := contextWithShutdown(context.Background())
	defer cancel()

	err = runTx(ctx, nc, log, &handler{
		nc:  nc,
		nsc: nsc,
	})
	if err != nil {
		return fmt.Errorf("could not run transport: %w", err)
	}
	log.Info("shutting down")

	return nil
}

func contextWithShutdown(parent context.Context) (ctx context.Context, cancel func()) {
	sigs := make(chan os.Signal, 1)
	signal.Notify(sigs, syscall.SIGTERM, syscall.SIGINT)

	ctx, cancel = context.WithCancel(parent)

	go func() {
		defer signal.Stop(sigs)
		defer cancel()

		select {
		case <-sigs:
		case <-ctx.Done():
		}
	}()

	return
}

func needsInit(dataDir string) (bool, error) {
	files, err := ioutil.ReadDir(dataDir)
	if err != nil {
		return false, err
	}
	return len(files) == 0, nil
}

type config struct {
	nscPath string
	natsURL string
	dataDir string
}

func loadConfig() (*config, error) {
	natsURL := os.Getenv("NATS_URL")
	if natsURL == "" {
		return nil, errors.New("environment variable NATS_URL is not set")
	}

	dataDir := os.Getenv("NATS_MANAGER_DATA_DIR")
	if dataDir == "" {
		return nil, errors.New("environment variable NATS_MANAGER_DATA_DIR is not set")
	}

	nscPath := os.Getenv("NSC_PATH")
	if nscPath == "" {
		nscPath = "nsc"
	}

	return &config{
		natsURL: natsURL,
		dataDir: dataDir,
		nscPath: nscPath,
	}, nil
}
