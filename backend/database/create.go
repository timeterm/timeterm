package database

import (
	"context"

	"github.com/google/uuid"
)

type Organization struct {
	ID   uuid.UUID `db:"id"`
	Name string    `db:"name"`
}

type Student struct {
	ID             uuid.UUID `db:"id"`
	OrganizationID uuid.UUID `db:"organization_id"`
}

type Device struct {
	ID             uuid.UUID `db:"id"`
	OrganizationID uuid.UUID `db:"organization_id"`
	Name           string    `db:"name"`
	Status         string    `db:"status"`
}

func (w *Wrapper) CreateOrganization(ctx context.Context, name string) (Organization, error) {
	org := Organization{
		Name: name,
	}

	row := w.db.QueryRowContext(ctx, `INSERT INTO "organization" ("id", "name") VALUES (DEFAULT, $1) RETURNING "id"`, name)

	return org, row.Scan(&org.ID)
}

func (w *Wrapper) CreateStudent(ctx context.Context, organizationID uuid.UUID) (Student, error) {
	std := Student{
		OrganizationID: organizationID,
	}

	row := w.db.QueryRowContext(ctx, `INSERT INTO "student" ("id", "organizationID") VALUES (DEFAULT, $1) RETURNING "id"`, organizationID)

	return std, row.Scan(&std.ID)
}

func (w *Wrapper) CreateDevice(ctx context.Context, organizationID uuid.UUID, name string, status string) (Device, error) {
	dev := Device{
		OrganizationID: organizationID,
		Name:           name,
		Status:         status,
	}

	row := w.db.QueryRowContext(ctx, `INSERT INTO "device" ("id", "organizationID", "name", "status") VALUES (DEFAULT, $1, $2, $3) RETURNING "id"`, organizationID, name, status)

	return dev, row.Scan(&dev.ID)
}
