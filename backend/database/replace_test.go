package database

import (
	"context"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestWrapper_ReplaceOrganization(t *testing.T) {
	f := newFixture(t)
	defer f.Close()

	org, err := f.dbw.CreateOrganization(context.Background(), "name", "institution")
	require.NoError(t, err)

	org.Name = "new name"
	org.ZermeloInstitution = "institution asdf"

	err = f.dbw.ReplaceOrganization(context.Background(), org)
	require.NoError(t, err)

	gotOrg, err := f.dbw.GetOrganization(context.Background(), org.ID)
	require.NoError(t, err)

	assert.Equal(t, org, gotOrg)
}

func TestWrapper_ReplaceDevice(t *testing.T) {
	f := newFixture(t)
	defer f.Close()

	org1, err := f.dbw.CreateOrganization(context.Background(), "org2", "institution1")
	require.NoError(t, err)

	org2, err := f.dbw.CreateOrganization(context.Background(), "org2", "institution2")
	require.NoError(t, err)

	dev, err := f.dbw.CreateDevice(context.Background(), org1.ID, "dev1", DeviceStatusOnline)
	require.NoError(t, err)

	dev.Name = "test"
	dev.Status = DeviceStatusOffline
	dev.OrganizationID = org2.ID

	err = f.dbw.ReplaceDevice(context.Background(), dev)
	require.NoError(t, err)

	gotDev, err := f.dbw.GetDevice(context.Background(), dev.ID)
	require.NoError(t, err)

	assert.Equal(t, dev, gotDev)
}
