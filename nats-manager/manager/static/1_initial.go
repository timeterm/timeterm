package static

import (
	"github.com/nats-io/jwt/v2"

	"gitlab.com/timeterm/timeterm/nats-manager/jwtmigrate"
	"gitlab.com/timeterm/timeterm/nats-manager/pkg/jwtpatch"
	nmsdk "gitlab.com/timeterm/timeterm/nats-manager/pkg/sdk"
)

func migration1Initial() jwtmigrate.Migration {
	return jwtmigrate.Migration{
		Name:    "initial",
		Version: 1,
		CreateAccounts: jwtmigrate.AccountCreates{
			backendAccount(): {},
			emdevsAccount(): {
				Patches: &jwtpatch.AccountClaimsPatches{
					AccountPatches: jwtpatch.AccountPatches{
						Exports: jwtpatch.ExportsPatches{
							Add: []*jwt.Export{
								{
									Name:     "EMDEV.>",
									Type:     jwt.Stream,
									Subject:  "EMDEV.>",
									TokenReq: true,
								},
							},
						},
					},
				},
			},
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
	}
}
