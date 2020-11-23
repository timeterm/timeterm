package static

import (
	"context"

	"github.com/go-logr/logr"
	"github.com/nats-io/jwt/v2"

	"gitlab.com/timeterm/timeterm/nats-manager/manager"
)

func AllowPub(subj ...string) manager.UserClaimsEditor {
	return func(c *jwt.UserClaims) {
		c.Pub.Allow.Add(subj...)
	}
}

func AllowSub(subj ...string) manager.UserClaimsEditor {
	return func(c *jwt.UserClaims) {
		c.Sub.Allow.Add(subj...)
	}
}

func ConfigureUsers(ctx context.Context, log logr.Logger, mgr *manager.Manager) error {
	log.Info("configuring static users")

	exists, err := mgr.AccountExists(ctx, "BACKEND")
	if err != nil {
		return err
	}
	if !exists {
		log.V(1).Info("account 'BACKEND' doesn't exist yet, creating it")
		if _, err = mgr.NewAccount(ctx, "BACKEND"); err != nil {
			return err
		}
		log.V(1).Info("account 'BACKEND' created")

		log.V(1).Info("user 'backend' doesn't exist yet (because account 'BACKEND' didn't exist yet)")
		if _, err = mgr.NewUser(
			ctx,
			"backend",
			"BACKEND",
			AllowPub(">"),
			AllowSub(">"),
		); err != nil {
			return err
		}
		log.V(1).Info("user 'backend' created")
	} else {
		log.V(1).Info("account 'BACKEND' already exists, not configuring it")
	}

	log.Info("static users configured")
	return nil
}
