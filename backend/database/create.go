package database

import (
	"context"
	"database/sql"
	"time"

	"github.com/google/uuid"
	"golang.org/x/crypto/sha3"
)

type Organization struct {
	ID                 uuid.UUID
	Name               string
	ZermeloInstitution string
}

type Student struct {
	ID             uuid.UUID
	OrganizationID uuid.UUID
	ZermeloUser    string
}

type OAuth2State struct {
	State       uuid.UUID
	Issuer      string
	RedirectURL string
	CreatedAt   time.Time
	ExpiresAt   time.Time
}

type OIDCFederation struct {
	OIDCIssuer   string
	OIDCSubject  string
	OIDCAudience string
	UserID       uuid.UUID
}

type User struct {
	ID             uuid.UUID
	Name           string
	Email          string
	OrganizationID uuid.UUID
}

type DeviceStatus string

const (
	DeviceStatusNotActivated DeviceStatus = "not_activated"
	DeviceStatusOK           DeviceStatus = "ok"
)

type Device struct {
	ID             uuid.UUID
	OrganizationID uuid.UUID
	Name           string
	Status         DeviceStatus
	LastHeartbeat  sql.NullTime
}

func (w *Wrapper) CreateOrganization(ctx context.Context,
	name string,
	zermeloInstitution string,
) (Organization, error) {
	org := Organization{
		Name:               name,
		ZermeloInstitution: zermeloInstitution,
	}

	row := w.db.QueryRowContext(ctx, `
		INSERT INTO "organization" ("id", "name", "zermelo_institution") 
		VALUES (DEFAULT, $1, $2) 
		RETURNING "id"
	`, name, zermeloInstitution)

	return org, row.Scan(&org.ID)
}

func (w *Wrapper) CreateStudent(ctx context.Context, organizationID uuid.UUID) (Student, error) {
	std := Student{
		OrganizationID: organizationID,
	}

	row := w.db.QueryRowContext(ctx, `
		INSERT INTO "student" ("id", "organization_id") 
		VALUES (DEFAULT, $1) 
		RETURNING "id"
	`, organizationID)

	return std, row.Scan(&std.ID)
}

func (w *Wrapper) CreateDevice(ctx context.Context,
	organizationID uuid.UUID,
	name string,
	status DeviceStatus,
) (Device, uuid.UUID, error) {
	dev := Device{
		OrganizationID: organizationID,
		Name:           name,
		Status:         status,
	}

	token := uuid.New()
	tokenHash, err := hashToken(token)
	if err != nil {
		return dev, token, err
	}

	tx, err := w.db.Beginx()
	if err != nil {
		return dev, token, err
	}
	defer func() { _ = tx.Rollback() }()

	err = tx.GetContext(ctx, &dev.ID, `
		INSERT INTO "device" (id, organization_id, name, status)
		VALUES (DEFAULT, $1, $2, $3)
		RETURNING "id"
	`, dev.OrganizationID, dev.Name, dev.Status)
	if err != nil {
		return dev, token, err
	}

	_, err = tx.ExecContext(ctx, `
		INSERT INTO "device_token" ("token_hash", "device_id")
		VALUES ($1, $2)
	`, tokenHash, dev.ID)
	if err != nil {
		return dev, token, err
	}

	return dev, token, tx.Commit()
}

func (w *Wrapper) CreateOAuth2State(ctx context.Context, issuer, redirectURL string) (OAuth2State, error) {
	state := OAuth2State{
		Issuer:      issuer,
		RedirectURL: redirectURL,
	}

	row := w.db.QueryRowContext(ctx, `
		INSERT INTO "oauth2_state" ("state", "issuer", "redirect_url")
		VALUES (DEFAULT, $1, $2)
		RETURNING "state"
	`, issuer, redirectURL)

	return state, row.Scan(&state.State)
}

func (w *Wrapper) CreateOIDCFederation(ctx context.Context, federation OIDCFederation) (OIDCFederation, error) {
	_, err := w.db.ExecContext(ctx, `
		INSERT INTO "oidc_federation" (oidc_subject, oidc_issuer, oidc_audience, user_id)
		VALUES ($1, $2, $3, $4)
	`, federation.OIDCSubject, federation.OIDCIssuer, federation.OIDCAudience, federation.UserID)

	return federation, err
}

func (w *Wrapper) CreateUser(ctx context.Context, name, email string, organizationID uuid.UUID) (User, error) {
	user := User{
		Name:           name,
		Email:          email,
		OrganizationID: organizationID,
	}

	row := w.db.QueryRowContext(ctx, `
		INSERT INTO "user" (id, name, email, organization_id) 
		VALUES (DEFAULT, $1, $2, $3)
		RETURNING "id"
	`, name, email, organizationID)

	return user, row.Scan(&user.ID)
}

func (w *Wrapper) CreateNewUser(ctx context.Context, name, email string, federation OIDCFederation) (User, error) {
	user := User{
		Name:  name,
		Email: email,
	}

	tx, err := w.db.Beginx()
	if err != nil {
		return user, err
	}
	defer func() { _ = tx.Rollback() }()

	err = tx.GetContext(ctx, &user.OrganizationID, `
		INSERT INTO "organization" (id, name, zermelo_institution)
		VALUES (DEFAULT, '', '')
		RETURNING "id"
	`)
	if err != nil {
		return user, err
	}

	_, err = tx.ExecContext(ctx, `
		INSERT INTO "oidc_federation" (oidc_subject, oidc_issuer, oidc_audience, user_id)
		VALUES ($1, $2, $3, $4)
	`, federation.OIDCSubject, federation.OIDCIssuer, federation.OIDCAudience, user.ID)
	if err != nil {
		return user, err
	}

	return user, tx.Commit()
}

func (w *Wrapper) CreateToken(ctx context.Context, userID uuid.UUID) (uuid.UUID, error) {
	token := uuid.New()
	tokenHash, err := hashToken(token)
	if err != nil {
		return token, err
	}

	_, err = w.db.ExecContext(ctx, `
		INSERT INTO "user_token" ("token_hash", "user_id", "created_at", "expires_at")
		VALUES ($1, $2, DEFAULT, DEFAULT)
	`, tokenHash, userID)

	return token, err
}

func hashToken(token uuid.UUID) ([]byte, error) {
	h := sha3.NewShake256()
	_, err := h.Write(token[:])
	if err != nil {
		return nil, err
	}

	hash := make([]byte, 64)
	_, err = h.Read(hash)
	if err != nil {
		return nil, err
	}

	return hash, err
}
