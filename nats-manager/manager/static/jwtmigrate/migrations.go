package jwtmigrate

import (
	"github.com/go-logr/logr"

	"gitlab.com/timeterm/timeterm/nats-manager/database"
	"gitlab.com/timeterm/timeterm/nats-manager/manager"
	"gitlab.com/timeterm/timeterm/nats-manager/manager/static/jwtpatch"
	"gitlab.com/timeterm/timeterm/nats-manager/secrets"
)

var migrations = Migrations{
	{
		Name:    "Initial",
		Version: 1,
		CreateAccounts: map[string]accountCreate{
			"BACKEND": {},
		},
		CreateUsers: map[string]userCreate{
			"backend": {
				accountName: "BACKEND",
				patches: &jwtpatch.UserClaimsPatches{
					UserPatches: jwtpatch.UserPatches{
						PermissionsPatches: jwtpatch.PermissionsPatches{
							Pub: &jwtpatch.PermissionPatches{
								Allow: jwtpatch.StringListPatches{
									Add: []string{">"},
								},
							},
							Sub: &jwtpatch.PermissionPatches{
								Allow: jwtpatch.StringListPatches{
									Add: []string{">"},
								},
							},
						},
					},
				},
			},
		},
	},
}

func RunStaticMigrations(log logr.Logger, dbw *database.Wrapper, mgr *manager.Manager, st *secrets.Store) error {
	return migrations.Run(log, dbw, mgr, st)
}
