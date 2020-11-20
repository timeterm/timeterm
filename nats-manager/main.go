package main

import (
	"context"
	"errors"
	"fmt"
	"io/ioutil"
	"os"
	"os/signal"
	"sync/atomic"
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
	"gitlab.com/timeterm/timeterm/nats-manager/secrets"
	"gitlab.com/timeterm/timeterm/nats-manager/secrets/static"
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

	vc, err := vault.NewClient(vault.DefaultConfig())
	if err != nil {
		return fmt.Errorf("could not create Vault client: %w", err)
	}

	firstRun := false
	dbw, err := database.New(
		cfg.databaseURL,
		log,
		database.WithMigrationHooks(database.MigrationHooks{
			After: []database.MigrationHook{
				func(from, to uint) error {
					if from == 0 {
						firstRun = true
					}
					return nil
				},
			},
		}),
	)
	if err != nil {
		return fmt.Errorf("could not connect to database: %w", err)
	}
	defer func() {
		if err = dbw.Close(); err != nil {
			log.Error(err, "could not close database")
		}
	}()

	vcc := secrets.NewVaultClient(cfg.vaultPrefix, vc)
	mgr, err := secrets.NewManager(log, vcc, dbw, secrets.DefaultOperatorConfig())
	if err != nil {
		return fmt.Errorf("could not create secrets manager: %w", err)
	}

	ctx, cancel := contextWithShutdown(context.Background())
	defer cancel()

	if firstRun {
		log.Info("first run, initializing")

		if err = mgr.Init(context.Background()); err != nil {
			return fmt.Errorf("could not init secrets manager: %w", err)
		}

		if err = static.ConfigureStaticUsers(ctx, mgr); err != nil {
			return fmt.Errorf("could not configure static users: %w", err)
		}

		log.Info("initialized")
	}

	srv := api.NewServer(log, vcc)

	eg, ctx := errgroup.WithContext(ctx)
	eg.Go(func() error {
		if err := srv.Serve(ctx, cfg.apiAddress); err != nil {
			return fmt.Errorf("could not serve API: %w", err)
		}
		return nil
	})
	eg.Go(func() error {
		if err := actuallyRunTx(ctx, cfg, mgr, log); err != nil {
			return fmt.Errorf("could not run transport: %w", err)
		}
		return nil
	})
	return eg.Wait()
}

func actuallyRunTx(ctx context.Context, cfg *config, mgr *secrets.Manager, log logr.Logger) error {
	creds, err := mgr.GenerateUserCredentials(ctx, "backend", "BACKEND")
	if err != nil {
		return fmt.Errorf("failed to generate backend user credentials: %w", err)
	}

	f, err := ioutil.TempFile("", "*.jwt")
	if err != nil {
		return fmt.Errorf("could not write NATS credentials: %w", err)
	}
	if _, err = f.WriteString(creds); err != nil {
		return fmt.Errorf("could not write NATS credentials: %w", err)
	}
	if err = f.Close(); err != nil {
		return fmt.Errorf("could not close NATS credentials file: %w", err)
	}

	nc, err := tryConnectNATS(ctx, log, cfg.natsURL, nats.UserCredentials(f.Name()))
	if err != nil {
		return fmt.Errorf("could not connect to NATS: %w", err)
	}
	defer func() {
		if err = nc.Drain(); err != nil {
			log.Error(err, "error draining NATS connection on shutdown")
		}
	}()

	if err = runTx(ctx, nc, log, &handler{
		nc: nc,
		mg: mgr,
	}); err != nil {
		return fmt.Errorf("could not run transport: %w", err)
	}
	log.Info("shutting down")

	return nil
}

func tryConnectNATS(ctx context.Context, log logr.Logger, url string, opts ...nats.Option) (*nats.Conn, error) {
	connected := make(chan *nats.Conn)
	stop := int32(0)

	go func() {
		tick := time.NewTicker(time.Second * 5)
		defer tick.Stop()

		for atomic.LoadInt32(&stop) == 0 {
			select {
			case <-ctx.Done():
				return
			case <-tick.C:
			}

			nc, err := nats.Connect(url, opts...)
			if err == nil {
				log.Info("connected to NATS")
				connected <- nc
				return
			} else {
				log.Error(err, "error connecting to NATS (will likely retry unless shutting down)")
			}
		}
	}()

	select {
	case <-ctx.Done():
		atomic.StoreInt32(&stop, 1)
		return nil, ctx.Err()
	case nc := <-connected:
		return nc, nil
	}
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

type config struct {
	apiAddress   string
	natsURL      string
	databaseURL  string
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

	operatorName := os.Getenv("OPERATOR_NAME")
	if operatorName == "" {
		return nil, errors.New("environment variable OPERATOR_NAME is not set")
	}

	return &config{
		apiAddress:   apiAddress,
		natsURL:      natsURL,
		databaseURL:  databaseURL,
		vaultPrefix:  vaultPrefix,
		operatorName: operatorName,
	}, nil
}
