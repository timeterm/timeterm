package api

import (
	"database/sql"
	"time"

	"github.com/google/uuid"

	"gitlab.com/timeterm/timeterm/backend/database"
)

type Organization struct {
	ID      uuid.UUID               `json:"id"`
	Name    string                  `json:"name"`
	Zermelo OrganizationZermeloInfo `json:"zermelo"`
}

type OrganizationZermeloInfo struct {
	Institution string `json:"institution"`
}

type Student struct {
	ID             uuid.UUID          `json:"id"`
	OrganizationID uuid.UUID          `json:"organizationId"`
	Zermelo        StudentZermeloInfo `json:"zermelo"`
}

type StudentZermeloInfo struct {
	User string `json:"user"`
}

type PrimaryDeviceStatus string
type SecondaryDeviceStatus string

const (
	PrimaryDeviceStatusOnline  = "Online"
	PrimaryDeviceStatusOffline = "Offline"

	SecondaryDeviceStatusNotActivated = "NotActivated"
	SecondaryDeviceStatusOK           = "Ok"
)

type Device struct {
	ID              uuid.UUID             `json:"id"`
	OrganizationID  uuid.UUID             `json:"organizationId"`
	Name            string                `json:"name"`
	PrimaryStatus   PrimaryDeviceStatus   `json:"primaryStatus"`
	SecondaryStatus SecondaryDeviceStatus `json:"secondaryStatus"`
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

type PaginatedStudents struct {
	Pagination
	Data []Student `json:"data"`
}

type CreateDeviceResponse struct {
	Device Device    `json:"device"`
	Token  uuid.UUID `json:"token"`
}

func OrganizationFrom(org database.Organization) Organization {
	return Organization{
		ID:   org.ID,
		Name: org.Name,
		Zermelo: OrganizationZermeloInfo{
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
		Zermelo: StudentZermeloInfo{
			User: student.ZermeloUser,
		},
	}
}

func SecondaryDeviceStatusFrom(s database.DeviceStatus) SecondaryDeviceStatus {
	switch s {
	case database.DeviceStatusNotActivated:
		return SecondaryDeviceStatusNotActivated
	case database.DeviceStatusOK:
		fallthrough
	default:
		return SecondaryDeviceStatusOK
	}
}

func DeviceStatusToDB(s SecondaryDeviceStatus) database.DeviceStatus {
	switch s {
	case SecondaryDeviceStatusNotActivated:
		return database.DeviceStatusNotActivated
	case SecondaryDeviceStatusOK:
		fallthrough
	default:
		return database.DeviceStatusOK
	}
}

func lastHeartbeatToPrimaryDeviceStatus(t sql.NullTime) PrimaryDeviceStatus {
	if t.Valid && t.Time.After(time.Now().Add(-1*time.Minute)) {
		return PrimaryDeviceStatusOnline
	}
	return PrimaryDeviceStatusOffline
}

func DeviceFrom(device database.Device) Device {
	return Device{
		ID:              device.ID,
		OrganizationID:  device.OrganizationID,
		Name:            device.Name,
		PrimaryStatus:   lastHeartbeatToPrimaryDeviceStatus(device.LastHeartbeat),
		SecondaryStatus: SecondaryDeviceStatusFrom(device.Status),
	}
}

func CreateDeviceResponseFrom(device database.Device, token uuid.UUID) CreateDeviceResponse {
	return CreateDeviceResponse{
		Device: DeviceFrom(device),
		Token:  token,
	}
}

func DeviceToDB(device Device) database.Device {
	return database.Device{
		ID:             device.ID,
		OrganizationID: device.OrganizationID,
		Name:           device.Name,
		Status:         DeviceStatusToDB(device.SecondaryStatus),
	}
}

func DevicesFrom(dbDevices []database.Device) []Device {
	apiDevices := make([]Device, len(dbDevices))

	for i, dev := range dbDevices {
		apiDevices[i] = DeviceFrom(dev)
	}

	return apiDevices
}

func StudentsFrom(dbStudents []database.Student) []Student {
	apiStudents := make([]Student, len(dbStudents))

	for i, std := range dbStudents {
		apiStudents[i] = StudentFrom(std)
	}

	return apiStudents
}

func PaginatedDevicesFrom(p database.PaginatedDevices) PaginatedDevices {
	return PaginatedDevices{
		Pagination: PaginationFrom(p.Pagination),
		Data:       DevicesFrom(p.Devices),
	}
}

func PaginatedStudentsFrom(p database.PaginatedStudents) PaginatedStudents {
	return PaginatedStudents{
		Pagination: PaginationFrom(p.Pagination),
		Data:       StudentsFrom(p.Students),
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

func StudentToDB(s Student) database.Student {
	return database.Student{
		ID:             s.ID,
		OrganizationID: s.OrganizationID,
		ZermeloUser:    s.Zermelo.User,
	}
}
