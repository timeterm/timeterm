package static

import (
	"context"

	"gitlab.com/timeterm/timeterm/nats-manager/secrets"
)

func ConfigureStaticUsers(ctx context.Context, mgr *secrets.Manager) error {
	if _, err := mgr.NewAccount(ctx, "BACKEND"); err != nil {
		return err
	}

	if _, err := mgr.NewUser(ctx, "backend", "BACKEND"); err != nil {
		return err
	}

	return nil
}
