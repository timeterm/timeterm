package database

import (
	"database/sql"
	"fmt"
	mrand "math/rand"
	"net"
	"os"
	"testing"
	"time"

	"github.com/go-logr/logr"
	"github.com/go-logr/zapr"
	"github.com/ory/dockertest"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"go.uber.org/zap"
)

type connStringBuilder struct {
	user     string
	password string
	address  string
	dbName   string
	sslMode  string
}

func (b connStringBuilder) WithUser(user string) connStringBuilder {
	b.user = user
	return b
}

func (b connStringBuilder) WithPassword(password string) connStringBuilder {
	b.password = password
	return b
}

func (b connStringBuilder) WithAddress(address string) connStringBuilder {
	b.address = address
	return b
}

func (b connStringBuilder) WithDBName(dbName string) connStringBuilder {
	b.dbName = dbName
	return b
}

func (b connStringBuilder) WithSSLMode(sslMode string) connStringBuilder {
	b.sslMode = sslMode
	return b
}

func (b connStringBuilder) Build() string {
	connStr := fmt.Sprintf("postgres://%s:%s@%s/%s", b.user, b.password, b.address, b.dbName)
	if b.sslMode != "" {
		connStr += "?sslmode=" + b.sslMode
	}
	return connStr
}

var connString connStringBuilder

func TestMain(m *testing.M) {
	logger, _ := zap.NewDevelopment()
	defer func() { _ = logger.Sync() }()
	slog := logger.Sugar()

	connString = connString.WithUser("postgres").
		WithPassword("postgres").
		WithDBName("postgres").
		WithSSLMode("disable")

	slog.Info("Connecting with Docker")
	pool, err := dockertest.NewPool("")
	if err != nil {
		slog.Fatalf("Could not connect to Docker: %s", err)
	}

	slog.Info("Starting Postgres...")
	startTime := time.Now()
	resource, err := pool.Run("postgres", "12.3", []string{
		"POSTGRES_USER=postgres",
		"POSTGRES_PASSWORD=postgres",
		"POSTGRES_DB=timeterm",
	})
	if err != nil {
		slog.Fatalf("Could not start resource: %s", err)
	}

	// Exponential backoff-retry, because the application in the container might not be ready to accept connections yet
	if err := pool.Retry(func() error {
		connString = connString.WithAddress(net.JoinHostPort("localhost", resource.GetPort("5432/tcp")))

		db, err := sql.Open("postgres", connString.Build())
		if err != nil {
			return err
		}
		return db.Ping()
	}); err != nil {
		slog.Fatalf("Could not connect to Docker: %s", err)
	}
	slog.Infof("Postgres started in %s, running tests", time.Since(startTime))

	code := m.Run()

	slog.Info("Tests done, terminating Postgres")
	if err := pool.Purge(resource); err != nil {
		slog.Fatalf("Could not purge resource: %s", err)
	}
	slog.Info("Postgres terminated")

	os.Exit(code)
}

func createRandomDB(t *testing.T) string {
	db, err := sql.Open("postgres", connString.Build())
	require.NoError(t, err)

	name := fmt.Sprintf("timeterm_random_%d", mrand.New(mrand.NewSource(time.Now().UnixNano())).Uint32())
	_, err = db.Exec("CREATE DATABASE " + name)
	require.NoError(t, err)

	return name
}

func forceDropDB(t *testing.T, name string) {
	db, err := sql.Open("postgres", connString.Build())
	require.NoError(t, err)

	_, err = db.Exec(`
		UPDATE pg_database SET datallowconn = 'false' WHERE datname = $1;
	`, name)
	require.NoError(t, err)

	_, err = db.Exec(`
		SELECT pg_terminate_backend(pid)
		FROM pg_stat_activity
		WHERE datname = $1;
	`, name)
	require.NoError(t, err)

	_, err = db.Exec("DROP DATABASE " + name)
	require.NoError(t, err)
}

type fixture struct {
	t      *testing.T
	logger logr.Logger
	dbName string
	dbw    *Wrapper
}

func (f fixture) Close() {
	assert.NoError(f.t, f.dbw.Close())
	forceDropDB(f.t, f.dbName)
}

func newFixture(t *testing.T) fixture {
	logger, _ := zap.NewDevelopment()
	defer func() { _ = logger.Sync() }()
	log := zapr.NewLogger(logger)

	dbName := createRandomDB(t)
	dbw, err := New(connString.WithDBName(dbName).Build(), log,
		MigrationsURL("file://migrations"),
	)
	require.NoError(t, err)

	return fixture{
		t:      t,
		logger: log,
		dbName: dbName,
		dbw:    dbw,
	}
}
