package api

import (
	"database/sql"
	"time"

	"github.com/google/uuid"

	"gitlab.com/timeterm/timeterm/backend/database"
	devcfgpb "gitlab.com/timeterm/timeterm/proto/go/devcfg"
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

type EthernetServiceType string

const (
	EthernetServiceTypeEthernet EthernetServiceType = "Ethernet"
	EthernetServiceTypeWifi     EthernetServiceType = "Wifi"
)

type Ipv4ConfigType string

const (
	Ipv4ConfigTypeOff    Ipv4ConfigType = "Off"
	Ipv4ConfigTypeDhcp   Ipv4ConfigType = "Dhcp"
	Ipv4ConfigTypeCustom Ipv4ConfigType = "Custom"
)

type Ipv4Settings struct {
	Network string `json:"network"`
	Netmask string `json:"netmask"`
	Gateway string `json:"gateway"`
}

type Ipv4Config struct {
	Type     Ipv4ConfigType `json:"type"`
	Settings Ipv4Settings   `json:"settings"`
}

type Ipv6ConfigType string

const (
	Ipv6ConfigTypeOff    Ipv6ConfigType = "Off"
	Ipv6ConfigTypeAuto   Ipv6ConfigType = "Auto"
	Ipv6ConfigTypeCustom Ipv6ConfigType = "Custom"
)

type Ipv6Settings struct {
	Network      string `json:"network"`
	PrefixLength int    `json:"prefixLength"`
	Gateway      string `json:"gateway"`
}

type Ipv6Config struct {
	Type     Ipv6ConfigType `json:"type"`
	Settings Ipv6Settings   `json:"settings"`
}

type Ipv6Privacy string

const (
	Ipv6PrivacyDisabled  Ipv6Privacy = "Disabled"
	Ipv6PrivacyEnabled   Ipv6Privacy = "Enabled"
	Ipv6PrivacyPreferred Ipv6Privacy = "Preferred"
)

type Security string

const (
	SecurityPsk       Security = "Psk"
	SecurityLeee8021x Security = "Leee8021x"
	SecurityNone      Security = "None"
	SecurityWep       Security = "Wep"
)

type Eap string

const (
	EapTls  Eap = "Tls"
	EapTtls Eap = "Ttls"
	EapPeap Eap = "Peap"
)

type CaCertType string

const (
	CaCertTypePem CaCertType = "Pem"
	CaCertTypeDer CaCertType = "Der"
)

type PrivateKeyType string

const (
	PrivateKeyTypePem PrivateKeyType = "Pem"
	PrivateKeyTypeDer PrivateKeyType = "Der"
	PrivateKeyTypePfx PrivateKeyType = "Pfx"
)

type PrivateKeyPassphraseType string

const PrivateKeyPassphraseTypeFsid PrivateKeyPassphraseType = "Fsid"

type EthernetService struct {
	ID                       uuid.UUID                `json:"id"`
	Type                     EthernetServiceType      `json:"type"`
	Ipv4Config               Ipv4Config               `json:"ipv4Config"`
	Ipv6Config               Ipv6Config               `json:"ipv6Convig"`
	Ipv6Privacy              Ipv6Privacy              `json:"ipv6Privacy`
	Mac                      string                   `json:"mac"`
	Nameservers              []string                 `json:"nameservers"`
	SearchDomains            []string                 `json:"searchDomains"`
	Timeservers              []string                 `json:"timeservers"`
	Domain                   string                   `json:"domain"`
	Name                     string                   `json:"name"`
	SSID                     string                   `json:"ssid"`
	Passphrase               string                   `json:"passphrase"`
	Security                 Security                 `json:"security"`
	IsHidden                 bool                     `json:"isHidden"`
	Eap                      Eap                      `json:"eap"`
	CaCert                   byte                     `json:"caCert"`
	caCertType               CaCertType               `json:"caCertType"`
	PrivateKey               byte                     `json:"privateKey"`
	PrivateKeyType           PrivateKeyType           `json:"privateKeyType"`
	PrivateKeyPassphrase     string                   `json:"privateKeyPassphrase"`
	PrivateKeyPassphraseType PrivateKeyPassphraseType `json:"privateKeyPassphraseType"`
	Identity                 string                   `json:"identity"`
	AnonymousIdentify        string                   `json:"anonymousIdentify"`
	SubjectMatch             string                   `json:"subjectMatch"`
	AltSubjectMatch          string                   `json:"altSubjectMatch"`
	DomainSuffixMatch        string                   `json:"domainSuffixMatch"`
	DomainMatch              string                   `json:"domainMatch"`
	IsPhase2EapBased         bool                     `json:"isPhase2EapBased"`
}

func EthernetServiceTypeFrom(cfgType devcfgpb.EthernetServiceType) EthernetServiceType {
	switch cfgType {
	case 1:
		return EthernetServiceTypeEthernet
	case 2:
		return EthernetServiceTypeWifi
	default:
		return ""
	}
}

func EthernetConfigFrom(cfg *devcfgpb.EthernetService, id uuid.UUID) EthernetService {
	return EthernetService{
		ID:   id,
		Type: EthernetServiceTypeFrom(cfg.GetType()),
		// Go on here
	}
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

func UserToDB(user User) database.User {
	return database.User{
		ID:             user.ID,
		OrganizationID: user.OrganizationID,
		Email:          user.Email,
		Name:           user.Name,
	}
}

func StudentToDB(s Student) database.Student {
	return database.Student{
		ID:             s.ID,
		OrganizationID: s.OrganizationID,
		ZermeloUser:    s.Zermelo.User,
	}
}

type paginationParams struct {
	Offset    *uint64 `query:"offset"`
	MaxAmount *uint64 `query:"maxAmount"`
}
