package database

import (
	"context"

	"github.com/google/uuid"

)

func (w *Wrapper) GetOrganization(ctx context.Context, id uuid.UUID) (*Organization, error) {
	var organization Organization

	err := w.db.GetContext(ctx, &organization, `SELECT * FROM "organization" WHERE "id" = $1`, id)

	return &organization, err
}

func (w *Wrapper) GetStudent(ctx context.Context, id uuid.UUID) (*Student, error) {
	var student Student

	err := w.db.GetContext(ctx, &student, `SELECT * FROM "student" WHERE "id" = $1`, id)

	return &student, err
}

func (w *Wrapper) GetDevice(ctx context.Context, id uuid.UUID) (*Device, error) {
	var device Device

	err := w.db.GetContext(ctx, &device, `SELECT * FROM "device" WHERE "id" = $1`, id)

	return &device, err
}