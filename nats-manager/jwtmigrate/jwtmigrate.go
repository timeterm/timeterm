package jwtmigrate

import (
	"context"
	"fmt"
	"sort"
	"strconv"
	"strings"

	"github.com/go-logr/logr"
	"github.com/nats-io/jwt/v2"

	"gitlab.com/timeterm/timeterm/nats-manager/database"
	"gitlab.com/timeterm/timeterm/nats-manager/manager"
	"gitlab.com/timeterm/timeterm/nats-manager/pkg/jwtpatch"
)

const jwtTagPrefix = "timeterm.migration_version="

func getMigrationVersionFromTagList(t jwt.TagList) (int, bool) {
	for _, tag := range t {
		if strings.HasPrefix(tag, jwtTagPrefix) {
			versionStr := tag[len(jwtTagPrefix)-1:]
			version, err := strconv.Atoi(versionStr)
			if err != nil {
				return version, true
			}
		}
	}
	return 0, false
}

func getMigrationVersionFromOperator(c *jwt.OperatorClaims) (int, bool) {
	return getMigrationVersionFromTagList(c.Tags)
}

func getMigrationVersionFromAccount(c *jwt.AccountClaims) (int, bool) {
	return getMigrationVersionFromTagList(c.Tags)
}

func getMigrationVersionFromUser(c *jwt.UserClaims) (int, bool) {
	return getMigrationVersionFromTagList(c.Tags)
}

func removeMigrationVersionTags(t *jwt.TagList) {
	var remove []string
	for _, tag := range *t {
		if strings.HasPrefix(tag, jwtTagPrefix) {
			remove = append(remove, tag)
		}
	}
	t.Remove(remove...)
}

func setMigrationVersionInTagList(t *jwt.TagList, version int) {
	removeMigrationVersionTags(t)
	t.Add(fmt.Sprintf("%s%d", jwtTagPrefix, version))
}

func setMigrationVersionInOperator(c *jwt.OperatorClaims, version int) {
	setMigrationVersionInTagList(&c.Tags, version)
}

func setMigrationVersionInAccount(c *jwt.AccountClaims, version int) {
	setMigrationVersionInTagList(&c.Tags, version)
}

func setMigrationVersionInUser(c *jwt.UserClaims, version int) {
	setMigrationVersionInTagList(&c.Tags, version)
}

type OperatorMigration struct {
	NameRegex string
	Patch     func(log logr.Logger, r OperatorRef, claims *jwt.OperatorClaims)
}

type AccountMigration struct {
	NameRegex         string
	OperatorNameRegex string
	Patch             func(log logr.Logger, r AccountRef, claims *jwt.AccountClaims)
}

type UserMigration struct {
	NameRegex         string
	AccountNameRegex  string
	OperatorNameRegex string
	Patch             func(log logr.Logger, r UserRef, claims *jwt.UserClaims)
}

type OperatorRef struct {
	Name string
}

func (o OperatorRef) Account(name string) AccountRef {
	return AccountRef{Name: name, Operator: o}
}

type AccountRef struct {
	Name     string
	Operator OperatorRef
}

func (a AccountRef) User(name string) UserRef {
	return UserRef{Name: name, Account: a}
}

type AccountCreate struct {
	Patches *jwtpatch.AccountClaimsPatches
}

type UserRef struct {
	Name    string
	Account AccountRef
}

type UserCreate struct {
	Patches *jwtpatch.UserClaimsPatches
}

type (
	AccountCreates map[AccountRef]AccountCreate
	UserCreates    map[UserRef]UserCreate
)

type Migration struct {
	Name    string
	Version int

	CreateAccounts map[AccountRef]AccountCreate
	CreateUsers    map[UserRef]UserCreate

	OperatorsUp []*OperatorMigration
	AccountsUp  []*AccountMigration
	UsersUp     []*UserMigration
}

type Migrations []Migration

func (m Migrations) Len() int {
	return len(m)
}

func (m Migrations) Less(i, j int) bool {
	return m[i].Version < m[j].Version
}

func (m Migrations) Swap(i, j int) {
	m[i], m[j] = m[j], m[i]
}

func (m Migrations) Validate() error {
	sort.Sort(m)

	prevVersion := 0
	for _, migration := range m {
		if migration.Version != prevVersion+1 {
			return fmt.Errorf("expected to find migration version %d, got %d", migration.Version, prevVersion+1)
		}
	}

	return nil
}

// Run doesn't take a context because migrations should not be aborted (generally).
func (m Migrations) Run(log logr.Logger, dbw *database.Wrapper, mgr *manager.Manager) error {
	log = log.WithName("Migrations")
	log.Info("migrating")
	defer log.Info("done migrating")

	sort.Sort(m)
	if err := m.Validate(); err != nil {
		return fmt.Errorf("error validating migrations: %w", err)
	}

	currentVer, err := dbw.GetJWTMigrationVersion(context.Background())
	if err != nil {
		return fmt.Errorf("could not get current JWT migration version: %w", err)
	}

	if len(m) < currentVer {
		return fmt.Errorf(
			"current migration version (%d) is higher than what this version of nats-manager can handle (%d)",
			currentVer, len(m),
		)
	} else if len(m) == currentVer {
		log.Info("no migrations to perform")
		return nil
	}

	for _, migration := range m[currentVer:] {
		if err = migration.Run(log, dbw, mgr); err != nil {
			return fmt.Errorf("could not run migration %d (%q): %w", migration.Version, migration.Name, err)
		}
		if err = dbw.SetJWTMigrationVersion(context.Background(), migration.Version); err != nil {
			return fmt.Errorf("could not update migration version to %d: %w", migration.Version, err)
		}
	}

	return nil
}

func (m Migration) Run(log logr.Logger, dbw *database.Wrapper, mgr *manager.Manager) error {
	ctx := context.Background()

	log = log.WithValues("migrationVersion", m.Version, "name", m.Name)
	log.Info("running migration")

	for ref, acc := range m.CreateAccounts {
		if _, err := mgr.NewAccount(ctx, ref.Name, ref.Operator.Name, func(c *jwt.AccountClaims) {
			jwtpatch.PatchAccountClaims(c, acc.Patches)
			setMigrationVersionInAccount(c, m.Version)
		}); err != nil {
			return fmt.Errorf("could not create account %s with operator %s: %w", ref.Name, ref.Operator.Name, err)
		}
	}

	for ref, user := range m.CreateUsers {
		if _, err := mgr.NewUser(ctx, ref.Name, ref.Account.Name, ref.Account.Operator.Name, func(c *jwt.UserClaims) {
			jwtpatch.PatchUserClaims(c, user.Patches)
			setMigrationVersionInUser(c, m.Version)
		}); err != nil {
			return fmt.Errorf("could not create user %s under account %s with operator %s: %w",
				ref.Name, ref.Account.Name, ref.Account.Operator.Name, err,
			)
		}
	}

	for _, opm := range m.OperatorsUp {
		if err := dbw.WalkOperatorSubjectsRe(ctx, opm.NameRegex, func(op database.Operator) bool {
			if err := mgr.UpdateOperator(ctx, op.Name, func(c *jwt.OperatorClaims) {
				if v, ok := getMigrationVersionFromOperator(c); !ok || v == m.Version-1 {
					opm.Patch(log, OperatorRef{Name: op.Name}, c)
					setMigrationVersionInOperator(c, m.Version)
				}
			}); err != nil {
				log.Error(err, "could not update operator",
					"name", op.Name,
					"subject", op.Subject,
				)
			}

			return true
		}); err != nil {
			return fmt.Errorf("could not walk operator subjects: %w", err)
		}
	}

	for _, acm := range m.AccountsUp {
		if err := dbw.WalkAccountSubjectsRe(ctx, acm.NameRegex, acm.OperatorNameRegex, func(acc database.Account) bool {
			if err := mgr.UpdateAccount(ctx, acc.Name, acc.OperatorName, func(c *jwt.AccountClaims) {
				if v, ok := getMigrationVersionFromAccount(c); !ok || v == m.Version-1 {
					acm.Patch(log, AccountRef{
						Name: acc.Name,
						Operator: OperatorRef{
							Name: acc.OperatorName,
						},
					}, c)
					setMigrationVersionInAccount(c, m.Version)
				}
			}); err != nil {
				log.Error(err, "could not update account",
					"name", acc.Name,
					"operatorName", acc.OperatorName,
					"subject", acc.Subject,
				)
			}

			return true
		}); err != nil {
			return fmt.Errorf("could not walk account subjects: %w", err)
		}
	}

	for _, usm := range m.UsersUp {
		if err := dbw.WalkUserSubjectsRe(
			ctx,
			usm.NameRegex,
			usm.AccountNameRegex,
			usm.OperatorNameRegex,
			func(user database.User) bool {
				if err := mgr.UpdateUser(ctx, user.Name, user.AccountName, user.OperatorName, func(c *jwt.UserClaims) {
					if v, ok := getMigrationVersionFromUser(c); !ok || v == m.Version-1 {
						usm.Patch(log, UserRef{
							Name: user.Name,
							Account: AccountRef{
								Name: user.AccountName,
								Operator: OperatorRef{
									Name: user.OperatorName,
								},
							},
						}, c)
						setMigrationVersionInUser(c, m.Version)
					}
				}); err != nil {
					log.Error(err, "could not update user",
						"name", user.Name,
						"accountName", user.AccountName,
						"operatorName", user.OperatorName,
						"subject", user.Subject,
					)
				}

				return true
			},
		); err != nil {
			return fmt.Errorf("could not walk user subjects: %w", err)
		}
	}

	return nil
}
