package static

import (
	"github.com/go-logr/logr"

	"gitlab.com/timeterm/timeterm/nats-manager/database"
	"gitlab.com/timeterm/timeterm/nats-manager/manager"
	"gitlab.com/timeterm/timeterm/nats-manager/pkg/jwtmigrate"
	"gitlab.com/timeterm/timeterm/nats-manager/pkg/jwtpatch"
)

var jwtMigrations = jwtmigrate.Migrations{
	{
		Name:    "Initial",
		Version: 1,
		CreateAccounts: map[string]jwtmigrate.AccountCreate{
			"BACKEND": {},
		},
		CreateUsers: map[string]jwtmigrate.UserCreate{
			"backend": {
				AccountName: "BACKEND",
				Patches: &jwtpatch.UserClaimsPatches{
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

func RunJWTMigrations(log logr.Logger, dbw *database.Wrapper, mgr *manager.Manager) error {
	return jwtMigrations.Run(log, dbw, mgr)
}
