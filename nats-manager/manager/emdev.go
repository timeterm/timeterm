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
			fmt.Sprintf("$JS.API.CONSUMER.MSG.NEXT.EMDEV-DISOWN-TOKEN.EMDEV-%s-EMDEV-DISOWN-TOKEN", devID),
			fmt.Sprintf("$JS.ACK.EMDEV-DISOWN-TOKEN.EMDEV-%s-EMDEV-DISOWN-TOKEN.>", devID),
		)
	}
}
