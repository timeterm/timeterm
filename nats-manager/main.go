package main

import (
	"context"
	"errors"
	"fmt"
	"io/ioutil"
	"os"
	"os/signal"
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
	nc, err := nats.Connect(os.Getenv("NATS_URL"))
	if err != nil {
		return fmt.Errorf("could not connect to NATS: %w", err)
	}
	defer func() {
		err = nc.Drain()
		if err != nil {
			log.Error(err, "error draining NATS connection on shutdown")
		}
	}()

	nscConfig := &nscConfig{
		dataDir: os.Getenv("NATS_MANAGER_DATA_DIR"),
	}
	if nscConfig.dataDir == "" {
		return errors.New("environment variable NATS_MANAGER_DATA_DIR is not set")
	}

	needsInit, err := needsInit(nscConfig.dataDir)
	if err != nil {
		return fmt.Errorf("could not check if already initialized: %w", err)
	}

	if needsInit {
		err = nscInitCmd(nscConfig.dataDir).Run()
		if err != nil {
			return fmt.Errorf("could not init nsc: %w", err)
		}
	}

	ctx, cancel := contextWithShutdown(context.Background())
	defer cancel()

	err = runTx(ctx, nc, log, &handler{
		nc:     nc,
		nscCfg: nscConfig,
	})
	if err != nil {
		return fmt.Errorf("could not run transport: %w", err)
	}
	log.Info("shutting down")

	return nil
}

func contextWithShutdown(parent context.Context) (ctx context.Context, cancel func()) {
	sigs := make(chan os.Signal, 1)
	signal.Notify(sigs, syscall.SIGTERM, syscall.SIGKILL)

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
	return len(files) > 0, nil
}
