package manager

import (
	"fmt"

	"github.com/google/uuid"
	"github.com/nats-io/jwt/v2"
)

func emdevEditor(devID uuid.UUID) UserClaimsEditor {
	return func(c *jwt.UserClaims) {
		c.Sub.Allow.Add(
			fmt.Sprintf("EMDEV.%s.>", devID),
		)
		c.Pub.Allow.Add(
			fmt.Sprintf("$JS.API.CONSUMER.MSG.NEXT.EMDEV-RETRIEVE-NEW-NETWORKING-CONFIG.EMDEV-%s-EMDEV-RETRIEVE-NEW-NETWORKING-CONFIG", devID),
		)
	}
}
