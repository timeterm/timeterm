package database

import (
	"context"
	"database/sql"
	"strings"
	"time"

	sq "github.com/Masterminds/squirrel"
	"github.com/google/uuid"
	"github.com/lib/pq"
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

func min(x, y uint64) uint64 {
	if x < y {
		return x
	}
	return y
}

func or(x *uint64, y uint64) uint64 {
	if x != nil {
		return *x
	}
	return y
}

type Pagination struct {
	Offset, Limit, Total uint64
}

type PaginatedDevices struct {
	Pagination
	Devices []*Device
}

type GetDevicesOpts struct {
	OrganizationID uuid.UUID
	Limit          *uint64
	Offset         *uint64
	NameSearch     *string
}

func cleanSearch(s string) string {
	return strings.NewReplacer("%", "\\%", "_", "\\_").Replace(s)
}

func (w *Wrapper) GetDevices(ctx context.Context, opts GetDevicesOpts) (PaginatedDevices, error) {
	devs := PaginatedDevices{
		Pagination: Pagination{
			Limit:  min(or(opts.Limit, 50), 100),
			Offset: or(opts.Offset, 0),
		},
	}

	conds := sq.And{
		sq.Eq{"organization_id": opts.OrganizationID},
	}
	if opts.NameSearch != nil {
		conds = append(conds, sq.Expr("name LIKE '%' || ? || '%'", cleanSearch(*opts.NameSearch)))
	}

	buildQuery := func(b sq.SelectBuilder) sq.SelectBuilder {
		return b.
			From("device").
			Where(conds).
			PlaceholderFormat(sq.Dollar)
	}

	devsSql, args, err := buildQuery(sq.Select(`*`)).
		Limit(devs.Pagination.Limit).
		Offset(devs.Pagination.Offset).
		OrderBy("name ASC").
		ToSql()
	if err != nil {
		return devs, err
	}

	err = w.db.SelectContext(ctx, &devs.Devices, devsSql, args...)
	if err != nil {
		return devs, err
	}

	totalSql, args, err := buildQuery(sq.Select("COUNT(*)")).ToSql()
	if err != nil {
		return devs, err
	}

	err = w.db.GetContext(ctx, &devs.Total, totalSql, args...)
	if err != nil {
		return devs, err
	}

	return devs, nil
}

type GetStudentsOpts struct {
	OrganizationID uuid.UUID
	Limit          *uint64
	Offset         *uint64
}

type PaginatedStudents struct {
	Pagination
	Students []Student
}

func (w *Wrapper) GetStudents(ctx context.Context, opts GetStudentsOpts) (PaginatedStudents, error) {
	students := PaginatedStudents{
		Pagination: Pagination{
			Limit:  min(or(opts.Limit, 50), 100),
			Offset: or(opts.Offset, 0),
		},
	}

	conds := sq.And{
		sq.Eq{"student.organization_id": opts.OrganizationID},
	}

	buildQuery := func(b sq.SelectBuilder) sq.SelectBuilder {
		return b.
			From("student").
			Where(conds).
			PlaceholderFormat(sq.Dollar)
	}

	devsSql, args, err := buildQuery(sq.Select(`student.*, (COUNT(student_card.*) > 0) AS has_card_associated`)).
		LeftJoin("student_card ON student_card.student_id = student.id").
		Limit(students.Pagination.Limit).
		Offset(students.Pagination.Offset).
		OrderBy("zermelo_user ASC").
		GroupBy("student.id").
		ToSql()
	if err != nil {
		return students, err
	}

	err = w.db.SelectContext(ctx, &students.Students, devsSql, args...)
	if err != nil {
		return students, err
	}

	totalSql, args, err := buildQuery(sq.Select("COUNT(*)")).ToSql()
	if err != nil {
		return students, err
	}

	err = w.db.GetContext(ctx, &students.Total, totalSql, args...)
	if err != nil {
		return students, err
	}

	return students, nil
}

func (w *Wrapper) GetUserByID(ctx context.Context, id uuid.UUID) (User, error) {
	var user User

	err := w.db.GetContext(ctx, &user, `SELECT * FROM "user" WHERE "id" = $1`, id)

	return user, err
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

func (w *Wrapper) GetUserByToken(ctx context.Context, token uuid.UUID) (User, error) {
	var user User

	tokenHash, err := hashToken(token)
	if err != nil {
		return user, err
	}

	err = w.db.GetContext(ctx, &user, `
		SELECT u.* FROM "user_token"
		INNER JOIN "user" u on u.id = user_token.user_id
		WHERE "user_token"."token_hash" = $1 AND "expires_at" > now()
	`, tokenHash)

	return user, err
}

func (w *Wrapper) AreDevicesInOrganization(ctx context.Context,
	organizationID uuid.UUID,
	ids ...uuid.UUID,
) (bool, error) {
	var amountInOrganization int

	err := w.db.GetContext(ctx, &amountInOrganization, `
		SELECT COUNT(*) FROM "device"
		WHERE "id" = ANY($1)
		AND "organization_id" = $2
	`, pq.Array(ids), organizationID)

	return amountInOrganization == len(ids), err
}

func (w *Wrapper) AreStudentsInOrganization(ctx context.Context,
	organizationID uuid.UUID,
	ids ...uuid.UUID,
) (bool, error) {
	var amountInOrganization int

	err := w.db.GetContext(ctx, &amountInOrganization, `
		SELECT COUNT(*) FROM "student"
		WHERE "id" = ANY($1)
		AND "organization_id" = $2
	`, pq.Array(ids), organizationID)

	return amountInOrganization == len(ids), err
}

func (w *Wrapper) GetNetworkingService(ctx context.Context, id uuid.UUID) (NetworkingService, error) {
	var networkingService NetworkingService

	err := w.db.GetContext(ctx, &networkingService, `SELECT * FROM "networking_service" WHERE "id" = $1`, id)

	return networkingService, err
}

type GetNetworkingServicesOpts struct {
	OrganizationID uuid.UUID
	Limit          *uint64
	Offset         *uint64
}

type PaginatedNetworkingServices struct {
	Pagination
	NetworkingServices []NetworkingService
}

func (w *Wrapper) GetNetworkingServices(ctx context.Context, opts GetNetworkingServicesOpts) (PaginatedNetworkingServices, error) {
	netwServices := PaginatedNetworkingServices{
		Pagination: Pagination{
			Limit:  min(or(opts.Limit, 50), 100),
			Offset: or(opts.Offset, 0),
		},
	}

	conds := sq.And{
		sq.Eq{"organization_id": opts.OrganizationID},
	}

	buildQuery := func(b sq.SelectBuilder) sq.SelectBuilder {
		return b.
			From("networking_service").
			Where(conds).
			PlaceholderFormat(sq.Dollar)
	}

	netwServSql, args, err := buildQuery(sq.Select(`*`)).
		Limit(netwServices.Pagination.Limit).
		Offset(netwServices.Pagination.Offset).
		OrderBy("name ASC").
		ToSql()
	if err != nil {
		return netwServices, err
	}

	err = w.db.SelectContext(ctx, &netwServices.NetworkingServices, netwServSql, args...)
	if err != nil {
		return netwServices, err
	}

	totalSql, args, err := buildQuery(sq.Select("COUNT(*)")).ToSql()
	if err != nil {
		return netwServices, err
	}

	err = w.db.GetContext(ctx, &netwServices.Total, totalSql, args...)
	if err != nil {
		return netwServices, err
	}

	return netwServices, nil
}

func (w *Wrapper) GetAllNetworkingServices(ctx context.Context, organizationID uuid.UUID) ([]NetworkingService, error) {
	var networkingServices []NetworkingService

	err := w.db.SelectContext(ctx, &networkingServices, `
		SELECT * FROM "networking_service" WHERE "organization_id" = $1
	`, organizationID)

	return networkingServices, err
}

func (w *Wrapper) GetOrganizationByDeviceRegistrationToken(ctx context.Context, token uuid.UUID) (Organization, error) {
	var org Organization

	hash, err := hashToken(token)
	if err != nil {
		return org, err
	}

	err = w.db.GetContext(ctx, &org, `
		SELECT organization.* FROM device_registration_token
		INNER JOIN organization ON organization.id = device_registration_token.organization_id
		WHERE device_registration_token.token_hash = $1
	`, hash)

	return org, err
}

func (w *Wrapper) GetDeviceByToken(ctx context.Context, token uuid.UUID) (Device, error) {
	var dev Device

	hash, err := hashToken(token)
	if err != nil {
		return dev, err
	}

	err = w.db.GetContext(ctx, &dev, `
		SELECT device.* from device_token
		INNER JOIN device on device.id = device_token.device_id
		WHERE device_token.token_hash = $1
	`, hash)

	return dev, err
}

func (w *Wrapper) GetStudentByCard(ctx context.Context, uid []byte, organizationID uuid.UUID) (Student, error) {
	var student Student

	hash, err := hashBytes(uid)
	if err != nil {
		return student, err
	}

	err = w.db.GetContext(ctx, &student, `
		SELECT student.* FROM student_card
		INNER JOIN student ON student.id = student_card.student_id
		WHERE student_card.id_hash = $1 AND student_card.organization_id = $2
	`, hash, organizationID)

	return student, err
}

type GetAdminMessagesOpts struct {
	OrganizationID uuid.UUID
	Limit          *uint64
	FromTimestamp  *time.Time
}

func (w *Wrapper) GetAdminMessages(ctx context.Context, opts GetAdminMessagesOpts) ([]AdminMessage, error) {
	limit := min(or(opts.Limit, 50), 100)

	var fromTimestamp time.Time
	if opts.FromTimestamp != nil {
		fromTimestamp = *opts.FromTimestamp
	} else {
		fromTimestamp = time.Now()
	}

	conds := sq.And{
		sq.Eq{"organization_id": opts.OrganizationID},
	}

	buildQuery := func(b sq.SelectBuilder) sq.SelectBuilder {
		return b.
			From("admin_message").
			Where(conds).
			PlaceholderFormat(sq.Dollar)
	}

	devsSql, args, err := buildQuery(sq.Select(`*`)).
		Where(sq.Lt{"logged_at": fromTimestamp}).
		Limit(limit).
		OrderBy("logged_at DESC").
		ToSql()
	if err != nil {
		return nil, err
	}

	var messages []AdminMessage
	err = w.db.SelectContext(ctx, &messages, devsSql, args...)
	if err != nil {
		return nil, err
	}

	return messages, nil
}

func (w *Wrapper) GetAdminMessage(ctx context.Context, tm time.Time, organizationID uuid.UUID) (AdminMessage, error) {
	var message AdminMessage

	err := w.db.GetContext(ctx, &message, `
		SELECT * FROM admin_message
		WHERE logged_at = $1 AND organization_id = $2
	`, tm, organizationID)

	return message, err
}
