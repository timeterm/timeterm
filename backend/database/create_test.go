package database

import (
	"context"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestWrapper_CreateOrganization(t *testing.T) {
	const orgName = "test"

	f := newFixture(t)
	defer f.Close()

	org, err := f.dbw.CreateOrganization(context.Background(), orgName)

	assert.NoError(t, err)
	assert.NotZero(t, org.ID)
	assert.Equal(t, orgName, org.Name)
}

func TestWrapper_CreateDevice(t *testing.T) {
	const orgName = "test"
	const devName = "Device test"

	f := newFixture(t)
	defer f.Close()

	org, err := f.dbw.CreateOrganization(context.Background(), orgName)
	require.NoError(t, err)

	checkCreateDeviceResult := func(t *testing.T, status DeviceStatus, dev Device, err error) {
		t.Helper()

		assert.NoError(t, err)
		assert.NotZero(t, dev.ID)
		assert.Equal(t, org.ID, dev.OrganizationID)
		assert.Equal(t, devName, dev.Name)
		assert.Equal(t, status, dev.Status)
	}

	t.Run("Offline", func(t *testing.T) {
		const devStatus = DeviceStatusOffline

		dev, err := f.dbw.CreateDevice(context.Background(), org.ID, devName, devStatus)
		checkCreateDeviceResult(t, devStatus, dev, err)
	})

	t.Run("Online", func(t *testing.T) {
		const devStatus = DeviceStatusOffline

		dev, err := f.dbw.CreateDevice(context.Background(), org.ID, devName, devStatus)
		checkCreateDeviceResult(t, devStatus, dev, err)
	})
}

func TestWrapper_CreateStudent(t *testing.T) {
	const orgName = "test"

	f := newFixture(t)
	defer f.Close()

	org, err := f.dbw.CreateOrganization(context.Background(), orgName)
	require.NoError(t, err)

	student, err := f.dbw.CreateStudent(context.Background(), org.ID)
	assert.NoError(t, err)
	assert.NotZero(t, student.ID)
	assert.Equal(t, student.OrganizationID, org.ID)
}
