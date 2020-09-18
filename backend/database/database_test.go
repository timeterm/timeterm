package database

import (
	"database/sql"
	"fmt"
	"log"
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
	connString = connString.WithUser("postgres").
		WithPassword("postgres").
		WithDBName("postgres").
		WithSSLMode("disable")

	pool, err := dockertest.NewPool("")
	if err != nil {
		log.Fatalf("Could not connect to Docker: %s", err)
	}

	resource, err := pool.Run("postgres", "12.3", []string{
		"POSTGRES_USER=postgres",
		"POSTGRES_PASSWORD=postgres",
		"POSTGRES_DB=timeterm",
	})
	if err != nil {
		log.Fatalf("Could not start resource: %s", err)
	}

	// exponential backoff-retry, because the application in the container might not be ready to accept connections yet
	if err := pool.Retry(func() error {
		connString = connString.WithAddress(net.JoinHostPort("localhost", resource.GetPort("5432/tcp")))

		db, err := sql.Open("postgres", connString.Build())
		if err != nil {
			return err
		}
		return db.Ping()
	}); err != nil {
		log.Fatalf("Could not connect to Docker: %s", err)
	}

	code := m.Run()

	if err := pool.Purge(resource); err != nil {
		log.Fatalf("Could not purge resource: %s", err)
	}

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

func dropDB(t *testing.T, name string) {
	db, err := sql.Open("postgres", connString.Build())
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
	dropDB(f.t, f.dbName)
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
