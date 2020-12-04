package database

import (
	"context"
	"database/sql"
	"time"

	"github.com/google/uuid"
	"golang.org/x/crypto/sha3"
)

const DefaultTokenExpiration = time.Hour * 24

type Organization struct {
	ID                 uuid.UUID
	Name               string
	ZermeloInstitution string
}

type Student struct {
	ID                uuid.UUID
	OrganizationID    uuid.UUID
	ZermeloUser       sql.NullString
	HasCardAssociated bool
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

type NetworkingService struct {
	ID             uuid.UUID
	OrganizationID uuid.UUID
	Name           string
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
		INSERT INTO "organization" ("name", "zermelo_institution") 
		VALUES ($1, $2) 
		RETURNING "id"
	`, name, zermeloInstitution)

	return org, row.Scan(&org.ID)
}

func (w *Wrapper) CreateStudent(ctx context.Context, organizationID uuid.UUID) (Student, error) {
	std := Student{
		OrganizationID: organizationID,
	}

	row := w.db.QueryRowContext(ctx, `
		INSERT INTO "student" ("organization_id") 
		VALUES ($1) 
		RETURNING "id"
	`, organizationID)

	return std, row.Scan(&std.ID)
}

func (w *Wrapper) CreateNetworkingService(
	ctx context.Context,
	organizationID uuid.UUID,
	name string,
) (NetworkingService, error) {
	ns := NetworkingService{
		OrganizationID: organizationID,
		Name:           name,
	}

	row := w.db.QueryRowContext(ctx, `
		INSERT INTO "networking_service" (organization_id, name)
		VALUES ($1, $2)
		RETURNING "id"
	`, organizationID, name)

	return ns, row.Scan(&ns.ID)
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
		INSERT INTO "device" (organization_id, name, status)
		VALUES ($1, $2, $3)
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
		INSERT INTO "oauth2_state" ("issuer", "redirect_url")
		VALUES ($1, $2)
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
		INSERT INTO "user" (name, email, organization_id) 
		VALUES ($1, $2, $3)
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
		INSERT INTO "organization" (name, zermelo_institution)
		VALUES ('', '')
		RETURNING "id"
	`)
	if err != nil {
		return user, err
	}

	err = tx.GetContext(ctx, &user.ID, `
		INSERT INTO "user" (name, organization_id, email)
		VALUES ($1, $2, $3)
		RETURNING id
	`, user.Name, user.OrganizationID, user.Email)
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

func (w *Wrapper) CreateUserToken(ctx context.Context, userID uuid.UUID) (uuid.UUID, error) {
	token := uuid.New()
	tokenHash, err := hashToken(token)
	if err != nil {
		return token, err
	}

	_, err = w.db.ExecContext(ctx, `
		INSERT INTO "user_token" ("token_hash", "user_id")
		VALUES ($1, $2)
	`, tokenHash, userID)

	return token, err
}

func (w *Wrapper) CreateDeviceToken(ctx context.Context, deviceID uuid.UUID) (uuid.UUID, error) {
	token := uuid.New()
	tokenHash, err := hashToken(token)
	if err != nil {
		return token, err
	}

	_, err = w.db.ExecContext(ctx, `
		INSERT INTO "device_token" ("token_hash", "device_id")
		VALUES ($1, $2)
	`, tokenHash, deviceID)

	return token, err
}

func (w *Wrapper) CreateDeviceRegistrationToken(ctx context.Context, organizationID uuid.UUID) (uuid.UUID, error) {
	token := uuid.New()
	tokenHash, err := hashToken(token)
	if err != nil {
		return token, err
	}

	_, err = w.db.ExecContext(ctx, `
		INSERT INTO "device_registration_token" ("token_hash", "organization_id")
		VALUES ($1, $2)
	`, tokenHash, organizationID)

	return token, err
}

func hashToken(token uuid.UUID) ([]byte, error) {
	return hashBytes(token[:])
}

func hashBytes(b []byte) ([]byte, error) {
	h := sha3.NewShake256()
	_, err := h.Write(b)
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
