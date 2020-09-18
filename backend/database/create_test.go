package database

import (
	"context"
	"database/sql"
	"fmt"
	mrand "math/rand"
	"testing"
	"time"

	"github.com/go-logr/zapr"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"go.uber.org/zap"
)

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

func TestWrapper_CreateOrganization(t *testing.T) {
	logger, _ := zap.NewDevelopment()
	defer func() { _ = logger.Sync() }()
	log := zapr.NewLogger(logger)

	dbName := createRandomDB(t)
	dbw, err := New(connString.WithDBName(dbName).Build(), log,
		MigrationsURL("file://migrations"),
	)
	defer func() { 
		_ = dbw.Close() 
		dropDB(t, dbName)
	}()

	require.NoError(t, err)

	const orgName = "test"
	org, err := dbw.CreateOrganization(context.Background(), orgName)
	assert.NoError(t, err)
	assert.Equal(t, orgName, org.Name)
	assert.NotZero(t, org.ID)
}
