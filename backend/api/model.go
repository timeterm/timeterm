package api

import (
	"github.com/google/uuid"

	"gitlab.com/timeterm/timeterm/backend/database"
)

type Organization struct {
	ID      uuid.UUID   `json:"id"`
	Name    string      `json:"name"`
	Zermelo ZermeloInfo `json:"zermelo"`
}

type ZermeloInfo struct {
	Institution string `json:"institution"`
}

type Student struct {
	ID             uuid.UUID `json:"id"`
	OrganizationID uuid.UUID `json:"organizationId"`
}

type DeviceStatus string

const (
	DeviceStatusOnline  = "Online"
	DeviceStatusOffline = "Offline"
)

type Device struct {
	ID             uuid.UUID    `json:"id"`
	OrganizationID uuid.UUID    `json:"organizationId"`
	Name           string       `json:"name"`
	Status         DeviceStatus `json:"status"`
}

type User struct {
	ID             uuid.UUID `json:"id"`
	OrganizationID uuid.UUID `json:"organizationId"`
	Name           string    `json:"name"`
	Email          string    `json:"email"`
}

type Pagination struct {
	Offset    uint64 `json:"offset"`
	MaxAmount uint64 `json:"maxAmount"`
	Total     uint64 `json:"total"`
}

func PaginationFrom(p database.Pagination) Pagination {
	return Pagination{
		Offset:    p.Offset,
		MaxAmount: p.Limit,
		Total:     p.Total,
	}
}

type PaginatedDevices struct {
	Pagination
	Data []Device `json:"data"`
}

func OrganizationFrom(org database.Organization) Organization {
	return Organization{
		ID:   org.ID,
		Name: org.Name,
		Zermelo: ZermeloInfo{
			Institution: org.ZermeloInstitution,
		},
	}
}

func OrganisationToDB(org Organization) database.Organization {
	return database.Organization{
		ID:                 org.ID,
		Name:               org.Name,
		ZermeloInstitution: org.Zermelo.Institution,
	}
}

func StudentFrom(student database.Student) Student {
	return Student{
		ID:             student.ID,
		OrganizationID: student.OrganizationID,
	}
}

func DeviceStatusFrom(s database.DeviceStatus) DeviceStatus {
	switch s {
	case database.DeviceStatusOnline:
		return DeviceStatusOnline
	case database.DeviceStatusOffline:
		fallthrough
	default:
		return DeviceStatusOffline
	}
}

func DeviceStatusToDB(s DeviceStatus) database.DeviceStatus {
	switch s {
	case DeviceStatusOnline:
		return database.DeviceStatusOnline
	case DeviceStatusOffline:
		fallthrough
	default:
		return database.DeviceStatusOffline
	}
}

func DeviceFrom(device database.Device) Device {
	return Device{
		ID:             device.ID,
		OrganizationID: device.OrganizationID,
		Name:           device.Name,
		Status:         DeviceStatusFrom(device.Status),
	}
}

func DeviceToDB(device Device) database.Device {
	return database.Device{
		ID:             device.ID,
		OrganizationID: device.OrganizationID,
		Name:           device.Name,
		Status:         DeviceStatusToDB(device.Status),
	}
}

func DevicesFrom(dbDevices []database.Device) []Device {
	apiDevices := make([]Device, len(dbDevices))

	for i, dev := range dbDevices {
		apiDevices[i] = DeviceFrom(dev)
	}

	return apiDevices
}

func PaginatedDevicesFrom(p database.PaginatedDevices) PaginatedDevices {
	return PaginatedDevices{
		Pagination: PaginationFrom(p.Pagination),
		Data:       DevicesFrom(p.Devices),
	}
}

func UserFrom(user database.User) User {
	return User{
		ID:             user.ID,
		OrganizationID: user.OrganizationID,
		Name:           user.Name,
		Email:          user.Email,
	}
}
