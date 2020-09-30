package database

import (
	"context"

	"github.com/google/uuid"
)

type Organization struct {
	ID                 uuid.UUID
	Name               string
	ZermeloInstitution string
}

type Student struct {
	ID             uuid.UUID
	OrganizationID uuid.UUID
}

type DeviceStatus string

const (
	DeviceStatusOnline  DeviceStatus = "online"
	DeviceStatusOffline DeviceStatus = "offline"
)

type Device struct {
	ID             uuid.UUID
	OrganizationID uuid.UUID
	Name           string
	Status         DeviceStatus
}

type Devices struct {
	Devices []Device
}

func (w *Wrapper) CreateOrganization(ctx context.Context, name string, zermeloInstitution string) (Organization, error) {
	org := Organization{
		Name:               name,
		ZermeloInstitution: zermeloInstitution,
	}

	row := w.db.QueryRowContext(ctx, `INSERT INTO "organization" ("id", "name", "zermelo_institution") VALUES (DEFAULT, $1, $2) RETURNING "id"`, name, zermeloInstitution)

	return org, row.Scan(&org.ID)
}

func (w *Wrapper) CreateStudent(ctx context.Context, organizationID uuid.UUID) (Student, error) {
	std := Student{
		OrganizationID: organizationID,
	}

	row := w.db.QueryRowContext(ctx, `INSERT INTO "student" ("id", "organization_id") VALUES (DEFAULT, $1) RETURNING "id"`, organizationID)

	return std, row.Scan(&std.ID)
}

func (w *Wrapper) CreateDevice(ctx context.Context, organizationID uuid.UUID, name string, status DeviceStatus) (Device, error) {
	dev := Device{
		OrganizationID: organizationID,
		Name:           name,
		Status:         status,
	}

	row := w.db.QueryRowContext(ctx, `INSERT INTO "device" ("id", "organization_id", "name", "status") VALUES (DEFAULT, $1, $2, $3) RETURNING "id"`, organizationID, name, status)

	return dev, row.Scan(&dev.ID)
}
