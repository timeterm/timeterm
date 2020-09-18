package database

import (
	"context"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestWrapper_CreateOrganization(t *testing.T) {
	const orgName = "test"
	f := newFixture(t)
	defer f.Close()

	org, err := f.dbw.CreateOrganization(context.Background(), orgName)

	assert.NoError(t, err)
	assert.Equal(t, orgName, org.Name)
	assert.NotZero(t, org.ID)
}
