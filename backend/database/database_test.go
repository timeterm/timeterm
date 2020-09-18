package database

import (
	"database/sql"
	"fmt"
	"log"
	"net"
	"os"
	"testing"

	"github.com/ory/dockertest"
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
