package main

import (
	"os"
	"os/signal"
	"syscall"

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
		log.Error(err, "could not connect to NATS")
		os.Exit(1)
	}

	hdlr := handler{
		nc:      nc,
		dataDir: os.Getenv("NATS_MANAGER_DATA_DIR"),
	}

	err = newTx(nc, log, &hdlr)
	if err != nil {
		log.Error(err, "could not create transport")
		os.Exit(1)
	}

	sigs := make(chan os.Signal, 1)
	signal.Notify(sigs, syscall.SIGTERM, syscall.SIGKILL)

	<-sigs
}
