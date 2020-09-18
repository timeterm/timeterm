package database

import (
	"database/sql"
	"fmt"
	mrand "math/rand"
	"net"
	"net/url"
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

const postgresVersion = "12.3"

var connURI connURIBuilder

type connURIBuilder struct {
	prev     *connURIBuilder
	user     string
	password string
	address  string
	dbName   string
	opts     url.Values
}

func (b connURIBuilder) newLayer() connURIBuilder {
	return connURIBuilder{
		prev: &b,
		opts: make(url.Values),
	}
}

func (b connURIBuilder) WithUser(user string) connURIBuilder {
	b.user = user
	return b
}

func (b connURIBuilder) WithPassword(password string) connURIBuilder {
	b.password = password
	return b
}

func (b connURIBuilder) WithAddress(address string) connURIBuilder {
	b.address = address
	return b
}

func (b connURIBuilder) WithDBName(dbName string) connURIBuilder {
	b.dbName = dbName
	return b
}

func (b connURIBuilder) withOpt(k, v string) connURIBuilder {
	newb := b.newLayer()
	newb.opts[k] = []string{v}
	return newb
}

func (b connURIBuilder) WithSSLMode(sslMode string) connURIBuilder {
	return b.withOpt("sslmode", sslMode)
}

func (b connURIBuilder) cloneOpts() url.Values {
	opts := make(url.Values)
	for k, vs := range b.opts {
		opts[k] = vs
	}
	return opts
}

func (b connURIBuilder) inherit() connURIBuilder {
	b.opts = b.cloneOpts()

	if b.prev != nil {
		prev := b.prev.inherit()

		for opt, vs := range prev.opts {
			b.opts[opt] = make([]string, len(vs))
			for i, v := range vs {
				b.opts[opt][i] = v
			}
		}

		if b.user == "" {
			b.user = prev.user
		}
		if b.password == "" {
			b.password = prev.password
		}
		if b.address == "" {
			b.address = prev.address
		}
		if b.dbName == "" {
			b.dbName = prev.dbName
		}
	}

	return b
}

func (b connURIBuilder) Build() string {
	final := b.inherit()

	return (&url.URL{
		Scheme:   "postgres",
		User:     url.UserPassword(final.user, final.password),
		Host:     final.address,
		Path:     final.dbName,
		RawQuery: final.opts.Encode(),
	}).String()
}

func TestMain(m *testing.M) {
	logger, _ := zap.NewDevelopment()
	defer func() { _ = logger.Sync() }()
	slog := logger.Sugar()

	connURI = connURI.
		WithUser("postgres").
		WithPassword("postgres").
		WithDBName("postgres").
		WithSSLMode("disable")

	slog.Info("Connecting with Docker")
	pool, err := dockertest.NewPool("")
	if err != nil {
		slog.Fatalf("Could not connect to Docker: %v", err)
	}

	slog.Info("Starting Postgres...")
	startTime := time.Now()
	resource, err := pool.Run("postgres", postgresVersion, []string{
		"POSTGRES_USER=postgres",
		"POSTGRES_PASSWORD=postgres",
		"POSTGRES_DB=postgres",
	})
	if err != nil {
		slog.Fatalf("Could not start resource: %v", err)
	}

	err = resource.Expire(3600)
	if err != nil {
		slog.Fatalf("Could not set expiration for Postgres: %v", err)
	}

	// Retry with exponential backoff, because the application in the container might
	// not be ready to accept connections yet immediately after starting.
	if err := pool.Retry(func() error {
		connURI = connURI.WithAddress(net.JoinHostPort("localhost", resource.GetPort("5432/tcp")))
		uri := connURI.Build()

		db, err := sql.Open("postgres", uri)
		if err != nil {
			return err
		}
		return db.Ping()
	}); err != nil {
		slog.Fatalf("Could not connect to Docker: %v", err)
	}
	slog.Infof("Postgres started in %s, running tests", time.Since(startTime))

	code := m.Run()

	slog.Info("Tests done, terminating Postgres")
	if err := pool.Purge(resource); err != nil {
		slog.Fatalf("Could not purge resource: %v", err)
	}
	slog.Info("Postgres terminated successfully")

	os.Exit(code)
}

// createRandomDB creates a new random database with a unique name, safe for concurrently running tests.
func createRandomDB(t *testing.T) string {
	db, err := sql.Open("postgres", connURI.Build())
	require.NoError(t, err)

	// Grab a random number and create the database name with it.
	randNum := mrand.New(mrand.NewSource(time.Now().UnixNano())).Uint32()
	name := fmt.Sprintf("timeterm_random_%d", randNum)

	// Create the database.
	_, err = db.Exec("CREATE DATABASE " + name)
	require.NoError(t, err)

	return name
}

// forceDropDB forces the drop of a database.
func forceDropDB(t *testing.T, name string) {
	db, err := sql.Open("postgres", connURI.Build())
	require.NoError(t, err)

	// Don't allow anyone to connect anymore.
	_, err = db.Exec(`
		UPDATE pg_database SET datallowconn = 'false' WHERE datname = $1;
	`, name)
	require.NoError(t, err)

	// Terminate client connections that are still alive.
	_, err = db.Exec(`
		SELECT pg_terminate_backend(pid)
		FROM pg_stat_activity
		WHERE datname = $1;
	`, name)
	require.NoError(t, err)

	// Drop the database.
	_, err = db.Exec("DROP DATABASE " + name)
	require.NoError(t, err)
}

// The fixture contains the necessary information to run a database integration test.
type fixture struct {
	t      *testing.T
	logger logr.Logger
	dbName string
	dbw    *Wrapper
}

// Close frees all the resources that the fixture has used (currently a database).
func (f fixture) Close() {
	assert.NoError(f.t, f.dbw.Close())
	forceDropDB(f.t, f.dbName)
}

// newFixture creates a new fixture with a new random database and a logger.
func newFixture(t *testing.T) fixture {
	logger, _ := zap.NewDevelopment()
	defer func() { _ = logger.Sync() }()
	log := zapr.NewLogger(logger)

	dbName := createRandomDB(t)
	dbw, err := New(connURI.WithDBName(dbName).Build(), log,
		// Set the migrations URL to ./migrations, because the test is run in its own folder.
		WithMigrationsURL("file://migrations"),
	)
	require.NoError(t, err)

	return fixture{
		t:      t,
		logger: log,
		dbName: dbName,
		dbw:    dbw,
	}
}
