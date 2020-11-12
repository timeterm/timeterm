package main

import (
	"context"
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
	logger, _ := zap.NewDevelopment()
	defer func() { _ = logger.Sync() }()

	log := zapr.NewLogger(logger)
	defer log.Info("shutdown complete")

	nc, err := nats.Connect(os.Getenv("NATS_URL"))
	if err != nil {
		logFatal(log, err, "could not connect to NATS")
	}
	defer func() {
		err = nc.Drain()
		if err != nil {
			log.Error(err, "error draining NATS connection on shutdown")
		}
	}()

	dataDir := os.Getenv("NATS_MANAGER_DATA_DIR")
	needsInit, err := needsInit(dataDir)
	if err != nil {
		logFatal(log, err, "could not check if already initialized")
	}

	if needsInit {
		err = nscInitCmd(path.Join(dataDir, "store")).Run()
		if err != nil {
			logFatal(log, err, "could not init nsc")
		}
	}

	hdlr := handler{
		nc:      nc,
		dataDir: dataDir,
	}

	ctx, cancel := contextWithShutdown(context.Background())
	defer cancel()

	err = runTx(ctx, nc, log, &hdlr)
	if err != nil {
		logFatal(log, err, "could not run transport")
	}
}

func logFatal(log logr.Logger, err error, msg string, keysAndValues ...interface{}) {
	log.Error(err, "fatal: "+msg, keysAndValues...)
	os.Exit(1)
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
