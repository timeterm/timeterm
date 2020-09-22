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
	OrganizationID uuid.UUID `json:"organization_id"`
}

type Device struct {
	ID             uuid.UUID `json:"id"`
	OrganizationID uuid.UUID `json:"organization_id"`
	Name           string
	Status         database.DeviceStatus `json:"device_status"`
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

func StudentFrom(student database.Student) Student {
	return Student{
		ID:             student.ID,
		OrganizationID: student.OrganizationID,
	}
}

func DeviceFrom(device database.Device) Device {
	return Device{
		ID:             device.ID,
		OrganizationID: device.OrganizationID,
		Name:           device.Name,
		Status:         device.Status,
	}
}
