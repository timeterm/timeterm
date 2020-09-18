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
	"github.com/lib/pq"
	"github.com/ory/dockertest"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"go.uber.org/zap"
)

const postgresVersion = "12.3"

// connURI contains the 'base' connection URI for connecting to the running
// Postgres instance. It should not be modified by tests. It is configured in TestMain.
// Tests can use the methods of the connURIBuilder to set options to their content.
var connURI connURIBuilder

// connURIBuilder is a Postgres connection URI builder, helping with setting common values.
// Because we don't want to affect the global connURI, none of the methods has a pointer receiver.
// You can see the connURIBuilder as a fancy (singly) linked list: it keeps a reference to the previous 'layer',
// and creates a new layer if an option is set (so we don't completely have to clone the url.Values every time).
type connURIBuilder struct {
	prev     *connURIBuilder
	user     string
	password string
	address  string
	dbName   string
	opts     url.Values
}

// newLayer creates a new connURIBuilder layer with a reference to the receiver as previous layer.
func (b connURIBuilder) newLayer() connURIBuilder {
	return connURIBuilder{
		prev: &b,
		opts: make(url.Values),
	}
}

// WithUser sets the user which is used to log into Postgres with.
func (b connURIBuilder) WithUser(user string) connURIBuilder {
	b.user = user
	return b
}

// WithPassword sets the password which is used to log into Postgres with.
func (b connURIBuilder) WithPassword(password string) connURIBuilder {
	b.password = password
	return b
}

// WithAddress sets the address which the running Postgres instance is at.
// Should either be the host or a host:port combination.
func (b connURIBuilder) WithAddress(address string) connURIBuilder {
	b.address = address
	return b
}

// WithDBName sets the name of the database to connect with.
func (b connURIBuilder) WithDBName(dbName string) connURIBuilder {
	b.dbName = dbName
	return b
}

func (b connURIBuilder) withOpt(k, v string) connURIBuilder {
	newb := b.newLayer()
	newb.opts[k] = []string{v}
	return newb
}

// WithDBName sets the sslmode option.
// Valid values are: disable, allow, prefer, require, verify-ca, verify-full.
// See https://www.postgresql.org/docs/current/libpq-ssl.html
func (b connURIBuilder) WithSSLMode(sslMode string) connURIBuilder {
	return b.withOpt("sslmode", sslMode)
}

// cloneOpts makes a shallow copy of the options of the connURIBuilder.
func (b connURIBuilder) cloneOpts() url.Values {
	opts := make(url.Values)
	for k, vs := range b.opts {
		opts[k] = vs
	}
	return opts
}

// inherit creates a new connURIBuilder containing all settings from b and its ancestors.
// If a value is not set by b and b has an ancestor, the value is retrieved from the ancestor.
// In the case of options (such as sslmode), slices are not concatenated and/or deduplicated for simplicity.
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

// Build builds the connection URI.
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

	// Delete the database container after 1 hour, the tests really shouldn't run for that long
	// and the container shouldn't be leaked even if the tests are killed or stopped by uncatchable signals.
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
	_, err = db.Exec(fmt.Sprintf("CREATE DATABASE %s", pq.QuoteIdentifier(name)))
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
	_, err = db.Exec(fmt.Sprintf("DROP DATABASE %s", pq.QuoteIdentifier(name)))
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
