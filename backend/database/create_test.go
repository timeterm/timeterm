package database

import (
	"context"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestWrapper_CreateOrganization(t *testing.T) {
	const orgName = "test"
	const orgZermeloInstitution = "example"

	f := newFixture(t)
	defer f.Close()

	org, err := f.dbw.CreateOrganization(context.Background(), orgName, orgZermeloInstitution)

	assert.NoError(t, err)
	assert.NotZero(t, org.ID)
	assert.Equal(t, orgName, org.Name)
	assert.Equal(t, orgZermeloInstitution, org.ZermeloInstitution)
}

func TestWrapper_CreateDevice(t *testing.T) {
	const orgName = "test"
	const orgZermeloInstitution = "example"
	const devName = "Device test"

	f := newFixture(t)
	defer f.Close()

	org, err := f.dbw.CreateOrganization(context.Background(), orgName, orgZermeloInstitution)
	require.NoError(t, err)

	dev, tok, err := f.dbw.CreateDevice(context.Background(), org.ID, devName)

	assert.NoError(t, err)
	assert.NotZero(t, dev.ID)
	assert.NotZero(t, tok)
	assert.Equal(t, org.ID, dev.OrganizationID)
	assert.Equal(t, devName, dev.Name)
}

func TestWrapper_CreateStudent(t *testing.T) {
	const orgName = "test"
	const orgZermeloInstitution = "example"

	f := newFixture(t)
	defer f.Close()

	org, err := f.dbw.CreateOrganization(context.Background(), orgName, orgZermeloInstitution)
	require.NoError(t, err)

	student, err := f.dbw.CreateStudent(context.Background(), Student{
		OrganizationID: org.ID,
	})
	assert.NoError(t, err)
	assert.NotZero(t, student.ID)
	assert.Equal(t, student.OrganizationID, org.ID)
}
