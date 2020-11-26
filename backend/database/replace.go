package database

import (
	"context"

	"github.com/google/uuid"
)

func (w *Wrapper) ReplaceOrganization(ctx context.Context, org Organization) error {
	_, err := w.db.ExecContext(ctx,
		`UPDATE "organization" SET "name" = $1, "zermelo_institution" = $2 WHERE "id" = $3`,
		org.Name, org.ZermeloInstitution, org.ID,
	)

	return err
}

func (w *Wrapper) ReplaceDevice(ctx context.Context, dev Device) error {
	_, err := w.db.ExecContext(ctx,
		`UPDATE "device" SET "name" = $1, "organization_id" = $2, "status" = $3 WHERE "id" = $4`,
		dev.Name, dev.OrganizationID, dev.Status, dev.ID,
	)

	return err
}

func (w *Wrapper) ReplaceStudent(ctx context.Context, s Student) error {
	_, err := w.db.ExecContext(ctx,
		`UPDATE "student" SET "zermelo_user" = $1, "organization_id" = $2 WHERE "id" = $3`,
		s.ZermeloUser, s.OrganizationID, s.ID,
	)

	return err
}

func (w *Wrapper) ReplaceUser(ctx context.Context, user User) error {
	_, err := w.db.ExecContext(ctx,
		`UPDATE "user" SET "organization_id"= $1, "email" = $2, "name" = $3 WHERE "id" = $4`,
		user.OrganizationID, user.Email, user.Name, user.ID,
	)

	return err
}

func (w *Wrapper) ReplaceDeviceHeartbeat(ctx context.Context, id uuid.UUID) error {
	_, err := w.db.ExecContext(ctx,
		`UPDATE "device" SET "last_heartbeat" = TIME.NOW() WHERE "id" = $1`,
		id,
	)

	return err
}
