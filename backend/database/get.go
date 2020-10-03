package database

import (
	"context"
	"database/sql"
	"time"

	"github.com/google/uuid"
)

func (w *Wrapper) GetOrganization(ctx context.Context, id uuid.UUID) (Organization, error) {
	var organization Organization

	err := w.db.GetContext(ctx, &organization, `SELECT * FROM "organization" WHERE "id" = $1`, id)

	return organization, err
}

func (w *Wrapper) GetStudent(ctx context.Context, id uuid.UUID) (Student, error) {
	var student Student

	err := w.db.GetContext(ctx, &student, `SELECT * FROM "student" WHERE "id" = $1`, id)

	return student, err
}

func (w *Wrapper) GetDevice(ctx context.Context, id uuid.UUID) (Device, error) {
	var device Device

	err := w.db.GetContext(ctx, &device, `SELECT * FROM "device" WHERE "id" = $1`, id)

	return device, err
}

func (w *Wrapper) GetDevices(ctx context.Context) ([]Device, error) {
	var devices []Device

	err := w.db.GetContext(ctx, &devices, `SELECT * FROM "device"`)

	return devices, err
}

func (w *Wrapper) GetUserByOIDCFederation(ctx context.Context, federation OIDCFederation) (User, error) {
	var user User

	err := w.db.GetContext(ctx, &user, `
		SELECT "user".* FROM "user"
		INNER JOIN oidc_federation o ON "user".id = o.user_id
		WHERE o.oidc_subject = $1
		AND o.oidc_issuer = $2
		LIMIT 1
	`, federation.OIDCSubject, federation.OIDCIssuer)

	return user, err
}

func (w *Wrapper) GetUserByEmail(ctx context.Context, email string) (User, error) {
	var user User

	err := w.db.GetContext(ctx, &user, `SELECT * FROM "user" WHERE "email" = $1`, email)

	return user, err
}

func (w *Wrapper) GetOAuth2State(ctx context.Context, state uuid.UUID) (OAuth2State, error) {
	tx, err := w.db.Begin()
	if err != nil {
		return OAuth2State{}, err
	}
	defer func() {
		_ = tx.Rollback()
	}()

	var oauth2State OAuth2State
	err = w.db.GetContext(ctx, &oauth2State, `SELECT * FROM "oauth2_state" WHERE "state" = $1`, state)
	if err != nil {
		return oauth2State, err
	}

	_, err = w.db.ExecContext(ctx, `DELETE FROM "oauth2_state" WHERE "state" = $1`, state)
	if err != nil {
		return oauth2State, err
	}

	if oauth2State.ExpiresAt.Before(time.Now()) {
		if err = tx.Commit(); err != nil {
			return oauth2State, err
		}
		return oauth2State, sql.ErrNoRows
	}

	return oauth2State, tx.Commit()
}
