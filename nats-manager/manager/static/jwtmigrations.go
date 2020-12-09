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
)

func defaultOperator(operatorName string) jwtmigrate.OperatorRef {
	return jwtmigrate.OperatorRef{
		Name: operatorName,
	}
}

func backendAccount(operatorName string) jwtmigrate.AccountRef {
	return defaultOperator(operatorName).Account("BACKEND")
}

func emdevsAccount(operatorName string) jwtmigrate.AccountRef {
	return defaultOperator(operatorName).Account("EMDEVS")
}

func jwtMigrations() jwtmigrate.Migrations {
	return jwtmigrate.Migrations{
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
