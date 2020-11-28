package static

import (
	"github.com/go-logr/logr"

	"gitlab.com/timeterm/timeterm/nats-manager/database"
	"gitlab.com/timeterm/timeterm/nats-manager/jwtmigrate"
	"gitlab.com/timeterm/timeterm/nats-manager/manager"
	"gitlab.com/timeterm/timeterm/nats-manager/pkg/jwtpatch"
	nmsdk "gitlab.com/timeterm/timeterm/nats-manager/pkg/sdk"
)

func defaultOperator() jwtmigrate.OperatorRef {
	return jwtmigrate.OperatorRef{
		Name: "TIMETERM",
	}
}

func backendAccount() jwtmigrate.AccountRef {
	return defaultOperator().Account("BACKEND")
}

func emdevsAccount() jwtmigrate.AccountRef {
	return defaultOperator().Account("EMDEVS")
}

func jwtMigrations() jwtmigrate.Migrations {
	return jwtmigrate.Migrations{
		{
			Name:    "initial",
			Version: 1,
			CreateAccounts: jwtmigrate.AccountCreates{
				backendAccount(): {},
				emdevsAccount():  {},
			},
			CreateUsers: jwtmigrate.UserCreates{
				backendAccount().User("superuser"): {
					Patches: superuserPatches(),
				},
				backendAccount().User("backend"): {
					Patches: &jwtpatch.UserClaimsPatches{
						UserPatches: jwtpatch.UserPatches{
							PermissionsPatches: jwtpatch.PermissionsPatches{
								Pub: &jwtpatch.PermissionPatches{
									Allow: jwtpatch.StringListPatches{
										Add: []string{
											nmsdk.SubjectGenerateDeviceCredentials,
											nmsdk.SubjectProvisionNewDevice,
										},
									},
								},
							},
						},
					},
				},
				backendAccount().User("nats-manager"): {
					Patches: &jwtpatch.UserClaimsPatches{
						UserPatches: jwtpatch.UserPatches{
							PermissionsPatches: jwtpatch.PermissionsPatches{
								Pub: &jwtpatch.PermissionPatches{
									Allow: jwtpatch.StringListPatches{
										Add: []string{"NATS-MANAGER.>", "_INBOX.>"},
									},
								},
								Sub: &jwtpatch.PermissionPatches{
									Allow: jwtpatch.StringListPatches{
										Add: []string{"NATS-MANAGER.>"},
									},
								},
							},
						},
					},
				},
				emdevsAccount().User("superuser"): {
					Patches: superuserPatches(),
				},
			},
		},
	}
}

func superuserPatches() *jwtpatch.UserClaimsPatches {
	return &jwtpatch.UserClaimsPatches{
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
	}
}

func RunJWTMigrations(log logr.Logger, dbw *database.Wrapper, mgr *manager.Manager) error {
	return jwtMigrations().Run(log, dbw, mgr)
}
