package database

import (
	"context"
	"testing"

	"github.com/go-logr/zapr"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"go.uber.org/zap"
)

func TestWrapper_CreateOrganization(t *testing.T) {
	logger, _ := zap.NewDevelopment()
	defer func() { _ = logger.Sync() }()
	log := zapr.NewLogger(logger)

	dbw, err := New("postgres://postgres:postgres@localhost/timeterm?sslmode=disable", log,
		MigrationsURL("file://migrations"),
	)
	require.NoError(t, err)

	const orgName = "test"
	org, err := dbw.CreateOrganization(context.Background(), orgName)
	assert.NoError(t, err)
	assert.Equal(t, orgName, org.Name)
	assert.NotZero(t, org.ID)
}
