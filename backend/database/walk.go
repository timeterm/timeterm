package database

import (
	"context"
	"fmt"

	"github.com/google/uuid"
)

func (w *Wrapper) WalkDevices(ctx context.Context, organizationID uuid.UUID, f func(dev *Device) bool) error {
	var offset uint64
	for {
		devs, err := w.GetDevices(ctx, GetDevicesOpts{
			OrganizationID: organizationID,
			Limit:          PointerToUint64(50),
			Offset:         &offset,
		})
		if err != nil {
			return fmt.Errorf("could not retrieve devices with offset %d: %w", offset, err)
		}
		if len(devs.Devices) ==0 {
			break
		}
		offset += devs.Limit

		for _, dev := range devs.Devices {
			if !f(dev) {
				return nil
			}
		}
	}

	return nil
}
