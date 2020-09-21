package database

import (
	"context"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestWrapper_DeleteOrganization(t *testing.T) {
	f := newFixture(t)
	defer f.Close()

	org, err := f.dbw.CreateOrganization(context.Background(), "test", "example")
	require.NoError(t, err)

	err = f.dbw.DeleteOrganization(context.Background(), org.ID)
	assert.NoError(t, err)

	_, err = f.dbw.GetOrganization(context.Background(), org.ID)
	assert.Error(t, err)
}

func TestWrapper_DeleteDevice(t *testing.T) {
	f := newFixture(t)
	defer f.Close()

	org, err := f.dbw.CreateOrganization(context.Background(), "test", "example")
	require.NoError(t, err)

	dev, err := f.dbw.CreateDevice(context.Background(), org.ID, "test", DeviceStatusOnline)
	require.NoError(t, err)

	err = f.dbw.DeleteDevice(context.Background(), dev.ID)
	assert.NoError(t, err)

	_, err = f.dbw.GetDevice(context.Background(), dev.ID)
	assert.Error(t, err)
}

func TestWrapper_DeleteStudent(t *testing.T) {
	f := newFixture(t)
	defer f.Close()

	org, err := f.dbw.CreateOrganization(context.Background(), "test", "example")
	require.NoError(t, err)

	student, err := f.dbw.CreateStudent(context.Background(), org.ID)
	require.NoError(t, err)

	err = f.dbw.DeleteStudent(context.Background(), student.ID)
	assert.NoError(t, err)

	_, err = f.dbw.GetStudent(context.Background(), student.ID)
	assert.Error(t, err)
}
