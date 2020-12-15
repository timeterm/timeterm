package database

import (
	"context"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestWrapper_GetOrganization(t *testing.T) {
	f := newFixture(t)
	defer f.Close()

	want, err := f.dbw.CreateOrganization(context.Background(), "test", "example")
	require.NoError(t, err)

	got, err := f.dbw.GetOrganization(context.Background(), want.ID)
	assert.NoError(t, err)
	assert.Equal(t, got, want)
}

func TestWrapper_GetDevice(t *testing.T) {
	f := newFixture(t)
	defer f.Close()

	org, err := f.dbw.CreateOrganization(context.Background(), "test", "example")
	require.NoError(t, err)

	want, _, err := f.dbw.CreateDevice(context.Background(), org.ID, "example device")
	require.NoError(t, err)

	got, err := f.dbw.GetDevice(context.Background(), want.ID)
	assert.NoError(t, err)
	assert.Equal(t, got, want)
}

func TestWrapper_GetStudent(t *testing.T) {
	f := newFixture(t)
	defer f.Close()

	org, err := f.dbw.CreateOrganization(context.Background(), "test", "example")
	require.NoError(t, err)

	want, err := f.dbw.CreateStudent(context.Background(), Student{
		OrganizationID: org.ID,
	})
	require.NoError(t, err)

	got, err := f.dbw.GetStudent(context.Background(), want.ID)
	assert.NoError(t, err)
	assert.Equal(t, got, want)
}
