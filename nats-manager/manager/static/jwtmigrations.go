package static

import (
	"fmt"
	"strings"

	"github.com/go-logr/logr"
	"github.com/google/uuid"
	"github.com/nats-io/jwt/v2"

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
		{
			Name:    "Emdevs",
			Version: 2,
			CreateUsers: map[jwtmigrate.UserRef]jwtmigrate.UserCreate{
				emdevsAccount().User("backend"): {
					Patches: &jwtpatch.UserClaimsPatches{
						UserPatches: jwtpatch.UserPatches{
							PermissionsPatches: jwtpatch.PermissionsPatches{
								Pub: &jwtpatch.PermissionPatches{
									Allow: jwtpatch.StringListPatches{
										Add: []string{"EMDEV.>", "_INBOX.>"},
									},
								},
								Sub: &jwtpatch.PermissionPatches{
									Allow: jwtpatch.StringListPatches{
										Add: []string{"$JS.>"},
									},
								},
							},
						},
					},
				},
			},
			UsersUp: []*jwtmigrate.UserMigration{
				{
					NameRegex:        `^emdev-.*$`,
					AccountNameRegex: `^EMDEV$`,
					Patch: func(log logr.Logger, r jwtmigrate.UserRef, c *jwt.UserClaims) {
						id := strings.TrimPrefix(r.Name, "emdev-")
						uid, err := uuid.Parse(id)
						if err != nil {
							log.Error(err, "could not parse device ID in migration",
								"id", id, "userName", r.Name,
							)
							return
						}

						c.Sub.Allow.Add(
							fmt.Sprintf("EMDEV.%s.>", uid),
						)
						c.Pub.Allow.Add(
							fmt.Sprintf("$JS.API.CONSUMER.MSG.NEXT.EMDEV-DISOWN-TOKEN.EMDEV-%s-EMDEV-DISOWN-TOKEN", uid),
							fmt.Sprintf("$JS.ACK.EMDEV-DISOWN-TOKEN.EMDEV-%s-EMDEV-DISOWN-TOKEN.>", uid),
						)
					},
				},
			},
		},
		{
			Name:    "RETRIEVE-NEW-NETWORKING-CONFIG",
			Version: 3,
			UsersUp: []*jwtmigrate.UserMigration{
				{
					NameRegex:        `^emdev-.*$`,
					AccountNameRegex: `^EMDEV$`,
					Patch: func(log logr.Logger, r jwtmigrate.UserRef, c *jwt.UserClaims) {
						id := strings.TrimPrefix(r.Name, "emdev-")
						uid, err := uuid.Parse(id)
						if err != nil {
							log.Error(err, "could not parse device ID in migration",
								"id", id, "userName", r.Name,
							)
							return
						}

						c.Pub.Allow.Add(
							fmt.Sprintf("$JS.API.CONSUMER.MSG.NEXT.EMDEV-RETRIEVE-NEW-NETWORKING-CONFIG.EMDEV-%s-RETRIEVE-NEW-NETWORKING-CONFIG", uid),
						)
					},
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
