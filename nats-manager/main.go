package main

import (
	"bufio"
	"context"
	"errors"
	"fmt"
	"os"
	"os/signal"
	"strings"
	"syscall"
	"time"

	"github.com/go-logr/logr"
	"github.com/go-logr/zapr"
	vault "github.com/hashicorp/vault/api"
	_ "github.com/joho/godotenv/autoload"
	"github.com/nats-io/nats.go"
	"go.uber.org/zap"
	"golang.org/x/sync/errgroup"

	"gitlab.com/timeterm/timeterm/nats-manager/api"
	"gitlab.com/timeterm/timeterm/nats-manager/database"
	"gitlab.com/timeterm/timeterm/nats-manager/handler"
	"gitlab.com/timeterm/timeterm/nats-manager/manager"
	"gitlab.com/timeterm/timeterm/nats-manager/manager/static"
	"gitlab.com/timeterm/timeterm/nats-manager/secrets"
	"gitlab.com/timeterm/timeterm/nats-manager/transport"
)

func main() {
	start := time.Now()

	exitCode := 0
	defer os.Exit(exitCode)

	logger, _ := zap.NewDevelopment()
	defer func() { _ = logger.Sync() }()

	log := zapr.NewLogger(logger)
	defer log.Info("shutdown complete")

	logArt(log, ttMsg)

	if err := realMain(log, start); err != nil {
		if !errors.Is(err, context.Canceled) {
			log.Error(err, "error running nats-manager")
			exitCode = 1
		}
	}
}

func realMain(log logr.Logger, start time.Time) error {
	cfg, err := loadConfig()
	if err != nil {
		return fmt.Errorf("could not load configuration: %w", err)
	}

	vc, err := vault.NewClient(vault.DefaultConfig())
	if err != nil {
		return fmt.Errorf("could not create Vault client: %w", err)
	}

	dbw, err := database.New(cfg.databaseURL, log)
	if err != nil {
		return fmt.Errorf("could not connect to database: %w", err)
	}
	defer func() {
		if err = dbw.Close(); err != nil {
			log.Error(err, "could not close database")
		}
	}()

	sst := secrets.NewStore(log, cfg.vaultMount, cfg.vaultPrefix, vc)
	mgr, err := manager.New(log, sst, dbw, manager.DefaultOperatorConfig())
	if err != nil {
		return fmt.Errorf("could not create secrets manager: %w", err)
	}

	ctx, cancel := contextWithShutdown(context.Background(), log)
	defer cancel()

	if err = mgr.Init(ctx); err != nil {
		return fmt.Errorf("could not init secrets manager: %w", err)
	}

	if err = static.RunJWTMigrations(log, dbw, mgr); err != nil {
		return fmt.Errorf("could not run static JWT migrations: %w", err)
	}

	log.Info("started", "took", time.Since(start))

	go mgr.CheckJWTs(ctx)

	eg, ctx := errgroup.WithContext(ctx)
	eg.Go(func() error {
		srv := api.NewServer(log, sst, mgr)
		if err := srv.Serve(ctx, cfg.apiAddress); err != nil {
			return fmt.Errorf("could not serve API: %w", err)
		}
		return nil
	})
	eg.Go(func() error {
		return retry(ctx, func() {
			tryRunTransport(ctx, log, cfg, mgr)
		})
	})
	return eg.Wait()
}

// retry runs f indefinitely until ctx is canceled.
func retry(ctx context.Context, f func()) error {
	for {
		select {
		case <-ctx.Done():
			return ctx.Err()
		default:
			f()
		}
	}
}

func tryRunTransport(ctx context.Context, log logr.Logger, cfg *config, mgr *manager.Manager) {
	nc, err := trySetUpNATS(ctx, log, cfg, mgr)
	if err != nil {
		log.Error(err, "error setting up NATS")
		return
	}
	defer func() {
		if err = nc.Drain(); err != nil {
			log.Error(err, "error draining NATS connection on shutdown")
		}
	}()

	connectedCb := func(nc *nats.Conn) {
		log.Info("connected to NATS, configuring static streams")
		if err := static.ConfigureStreams(context.Background(), log, mgr); err != nil {
			log.Error(err, "error setting up static streams")
			return
		}
	}
	nc.SetDisconnectErrHandler(func(conn *nats.Conn, err error) {
		log.Error(err, "disconnected from NATS")
	})
	nc.SetClosedHandler(func(conn *nats.Conn) {
		log.Info("NATS connection closed")
	})
	nc.SetErrorHandler(func(conn *nats.Conn, subscription *nats.Subscription, err error) {
		log.Error(err, "error in NATS connection")
	})
	nc.SetReconnectHandler(connectedCb)
	connectedCb(nc)

	hdlr, err := handler.New(ctx, nc, mgr)
	if err != nil {
		log.Error(err, "could not create handler")
	}
	defer hdlr.Close()

	tx := transport.New(nc, log, hdlr)
	if err := tx.Run(ctx); err != nil {
		if !errors.Is(err, context.Canceled) {
			log.Error(err, "error running transport")
		}
	}
}

func trySetUpNATS(ctx context.Context, log logr.Logger, cfg *config, mgr *manager.Manager) (*nats.Conn, error) {
	nc, err := tryConnectNATS(ctx, log, cfg.natsURL,
		nats.UserJWT(mgr.NATSCredsCBs(ctx, "nats-manager", "BACKEND")),
		// Never stop trying to reconnect.
		nats.MaxReconnects(-1),
	)
	if err != nil {
		return nil, fmt.Errorf("could not connect to NATS: %w", err)
	}
	return nc, err
}

func tryConnectNATS(ctx context.Context, log logr.Logger, url string, opts ...nats.Option) (*nats.Conn, error) {
	connected := make(chan *nats.Conn, 1)
	stopped := make(chan struct{})

	go func() {
		tick := time.NewTicker(time.Second * 5)
		defer tick.Stop()
		defer close(stopped)
		firstTry := true

		for {
			if !firstTry {
				select {
				case <-ctx.Done():
					return
				case <-tick.C:
				}
			} else {
				firstTry = false
			}

			nc, err := nats.Connect(url, opts...)
			if err == nil {
				connected <- nc
				log.Info("connected to NATS")
				return
			}
			log.Error(err, "error connecting to NATS (will likely retry unless shutting down)")
		}
	}()

	select {
	case <-stopped:
		return nil, ctx.Err()
	case nc := <-connected:
		return nc, nil
	}
}

func contextWithShutdown(parent context.Context, log logr.Logger) (ctx context.Context, cancel func()) {
	sigs := make(chan os.Signal, 1)
	signal.Notify(sigs, syscall.SIGTERM, syscall.SIGINT)

	ctx, cancel = context.WithCancel(parent)

	go func() {
		defer signal.Stop(sigs)
		defer cancel()

		select {
		case <-sigs:
			log.Info("shutdown requested by OS")
		case <-ctx.Done():
		}
	}()

	return
}

type config struct {
	apiAddress   string
	natsURL      string
	databaseURL  string
	vaultMount   string
	vaultPrefix  string
	operatorName string
}

func loadConfig() (*config, error) {
	apiAddress := os.Getenv("API_ADDRESS")
	if apiAddress == "" {
		return nil, errors.New("environment variable API_ADDRESS is not set")
	}

	natsURL := os.Getenv("NATS_URL")
	if natsURL == "" {
		return nil, errors.New("environment variable NATS_URL is not set")
	}

	databaseURL := os.Getenv("DATABASE_URL")
	if databaseURL == "" {
		return nil, errors.New("environment variable DATABASE_URL is not set")
	}

	vaultPrefix := os.Getenv("VAULT_PREFIX")
	if vaultPrefix == "" {
		return nil, errors.New("environment variable VAULT_PREFIX is not set")
	}

	vaultMount := os.Getenv("VAULT_MOUNT")
	if vaultMount == "" {
		return nil, errors.New("environment variable VAULT_MOUNT is not set")
	}

	operatorName := os.Getenv("OPERATOR_NAME")
	if operatorName == "" {
		return nil, errors.New("environment variable OPERATOR_NAME is not set")
	}

	return &config{
		apiAddress:   apiAddress,
		natsURL:      natsURL,
		databaseURL:  databaseURL,
		vaultMount:   vaultMount,
		vaultPrefix:  vaultPrefix,
		operatorName: operatorName,
	}, nil
}

const ttMsg = `
 ╭──────────────╮
 ╰─────╮ ╭────╮ │
       │ │  ╭─╯ ╰─╮
       │ │  ╰─╮ ╭─╯
       │ │    │ ╰─╮
       ╰─╯    ╰───╯
       nats-manager

`

func logArt(l logr.Logger, s string) {
	scan := bufio.NewScanner(strings.NewReader(s))
	for scan.Scan() {
		l.Info(scan.Text())
	}
}
