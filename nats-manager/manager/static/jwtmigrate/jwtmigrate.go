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
	"gitlab.com/timeterm/timeterm/nats-manager/manager/static/jwtpatch"
	"gitlab.com/timeterm/timeterm/nats-manager/secrets"
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

type operatorMigration struct {
	nameRegex string
	patches   *jwtpatch.OperatorClaimsPatches
}

type accountMigration struct {
	nameRegex         string
	operatorNameRegex string
	patches           *jwtpatch.AccountClaimsPatches
}

type userMigration struct {
	nameRegex         string
	accountNameRegex  string
	operatorNameRegex string
	patches           *jwtpatch.UserClaimsPatches
}

type accountCreate struct {
	patches *jwtpatch.AccountClaimsPatches
}

type userCreate struct {
	accountName string
	patches     *jwtpatch.UserClaimsPatches
}

type Migration struct {
	Name    string
	Version int

	CreateAccounts map[string]accountCreate
	CreateUsers    map[string]userCreate

	OperatorsUp []*operatorMigration
	AccountsUp  []*accountMigration
	UsersUp     []*userMigration
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
func (m Migrations) Run(log logr.Logger, dbw *database.Wrapper, mgr *manager.Manager, st *secrets.Store) error {
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

	if len(m) <= currentVer {
		return fmt.Errorf(
			"current migration version (%d) is higher than what this version of nats-manager can handle (%d)",
			currentVer, len(m),
		)
	}

	for _, migration := range m[currentVer:] {
		if err := migration.Run(log, dbw, mgr, st); err != nil {
			return fmt.Errorf("could not run migration %d (%q): %w", migration.Version, migration.Name, err)
		}
	}

	return nil
}

func (m Migration) Run(log logr.Logger, dbw *database.Wrapper, mgr *manager.Manager, st *secrets.Store) error {
	ctx := context.Background()

	log.Info("running migration", "version", m.Version, "name", m.Name)

	for name, acc := range m.CreateAccounts {
		if _, err := mgr.NewAccount(ctx, name, func(c *jwt.AccountClaims) {
			jwtpatch.PatchAccountClaims(c, acc.patches)
			setMigrationVersionInAccount(c, m.Version)
		}); err != nil {
			return fmt.Errorf("could not create account %s: %w", name, err)
		}
	}

	for name, user := range m.CreateUsers {
		if _, err := mgr.NewUser(ctx, name, user.accountName, func(c *jwt.UserClaims) {
			jwtpatch.PatchUserClaims(c, user.patches)
			setMigrationVersionInUser(c, m.Version)
		}); err != nil {
			return fmt.Errorf("could not create user %s under account %s: %w", name, user.accountName, err)
		}
	}

	for _, opm := range m.OperatorsUp {
		if err := dbw.WalkOperatorSubjectsRe(ctx, opm.nameRegex, func(subject string) bool {
			tok, err := st.ReadOperatorJWT(subject)
			if err != nil {
				log.Error(err, "could not read operator JWT", "subject", subject)
				return false
			}

			if v, ok := getMigrationVersionFromOperator(tok); !ok || v == m.Version-1 {
				jwtpatch.PatchOperatorClaims(tok, opm.patches)
				setMigrationVersionInOperator(tok, m.Version)
			}

			if err = st.WriteOperatorJWT(tok, subject); err != nil {
				log.Error(err, "could not write operator JWT", "subject", subject)
			}

			return true
		}); err != nil {
			return fmt.Errorf("could not walk operator subjects: %w", err)
		}
	}

	for _, acm := range m.AccountsUp {
		if err := dbw.WalkAccountSubjectsRe(ctx, acm.nameRegex, acm.operatorNameRegex, func(subject string) bool {
			tok, err := st.ReadAccountJWT(subject)
			if err != nil {
				log.Error(err, "could not read account JWT", "subject", subject)
				return false
			}

			if v, ok := getMigrationVersionFromAccount(tok); !ok || v == m.Version-1 {
				jwtpatch.PatchAccountClaims(tok, acm.patches)
				setMigrationVersionInAccount(tok, m.Version)
			}

			if err = st.WriteAccountJWT(tok, subject); err != nil {
				log.Error(err, "could not write account JWT", "subject", subject)
			}

			return true
		}); err != nil {
			return fmt.Errorf("could not walk account subjects: %w", err)
		}
	}

	for _, usm := range m.UsersUp {
		if err := dbw.WalkUserSubjectsRe(
			ctx,
			usm.nameRegex,
			usm.accountNameRegex,
			usm.operatorNameRegex,
			func(subject string) bool {
				tok, err := st.ReadUserJWT(subject)
				if err != nil {
					log.Error(err, "could not read user JWT", "subject", subject)
					return false
				}

				if v, ok := getMigrationVersionFromUser(tok); !ok || v == m.Version-1 {
					jwtpatch.PatchUserClaims(tok, usm.patches)
					setMigrationVersionInUser(tok, m.Version)
				}

				if err = st.WriteUserJWT(tok, subject); err != nil {
					log.Error(err, "could not write user JWT", "subject", subject)
				}

				return true
			},
		); err != nil {
			return fmt.Errorf("could not walk user subjects: %w", err)
		}
	}

	return nil
}
