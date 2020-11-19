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
	Settings *Ipv4Settings  `json:"settings"`
}

type Ipv6ConfigType string

const (
	Ipv6ConfigTypeOff    Ipv6ConfigType = "Off"
	Ipv6ConfigTypeAuto   Ipv6ConfigType = "Auto"
	Ipv6ConfigTypeCustom Ipv6ConfigType = "Custom"
)

type Ipv6Settings struct {
	Network      string `json:"network"`
	PrefixLength uint64 `json:"prefixLength"`
	Gateway      string `json:"gateway"`
}

type Ipv6Config struct {
	Type     Ipv6ConfigType `json:"type"`
	Settings *Ipv6Settings  `json:"settings"`
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
	SecurityIeee8021x Security = "Ieee8021x"
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
	OrganizationID           uuid.UUID                `json:"organizationId"`
	Name                     string                   `json:"name"`
	Type                     EthernetServiceType      `json:"type"`
	Ipv4Config               *Ipv4Config              `json:"ipv4Config"`
	Ipv6Config               *Ipv6Config              `json:"ipv6Convig"`
	Ipv6Privacy              Ipv6Privacy              `json:"ipv6Privacy`
	Mac                      string                   `json:"mac"`
	Nameservers              []string                 `json:"nameservers"`
	SearchDomains            []string                 `json:"searchDomains"`
	Timeservers              []string                 `json:"timeservers"`
	Domain                   string                   `json:"domain"`
	NetworkName              string                   `json:"name"`
	SSID                     string                   `json:"ssid"`
	Passphrase               string                   `json:"passphrase"`
	Security                 Security                 `json:"security"`
	IsHidden                 bool                     `json:"isHidden"`
	Eap                      Eap                      `json:"eap"`
	CaCert                   []byte                   `json:"caCert"`
	CaCertType               CaCertType               `json:"caCertType"`
	PrivateKey               []byte                   `json:"privateKey"`
	PrivateKeyType           PrivateKeyType           `json:"privateKeyType"`
	PrivateKeyPassphrase     string                   `json:"privateKeyPassphrase"`
	PrivateKeyPassphraseType PrivateKeyPassphraseType `json:"privateKeyPassphraseType"`
	Identity                 string                   `json:"identity"`
	AnonymousIdentity        string                   `json:"anonymousIdentify"`
	SubjectMatch             string                   `json:"subjectMatch"`
	AltSubjectMatch          string                   `json:"altSubjectMatch"`
	DomainSuffixMatch        string                   `json:"domainSuffixMatch"`
	DomainMatch              string                   `json:"domainMatch"`
	IsPhase2EapBased         bool                     `json:"isPhase2EapBased"`
}

func ethernetServiceTypeFrom(cfgType devcfgpb.EthernetServiceType) EthernetServiceType {
	switch cfgType {
	case devcfgpb.EthernetServiceType_ETHERNET_SERVICE_TYPE_ETHERNET:
		return EthernetServiceTypeEthernet
	case devcfgpb.EthernetServiceType_ETHERNET_SERVICE_TYPE_WIFI:
		return EthernetServiceTypeWifi
	default:
		return ""
	}
}

func ipv4ConfigTypeFrom(ipv4ConfigType devcfgpb.Ipv4ConfigType) Ipv4ConfigType {
	switch ipv4ConfigType {
	case devcfgpb.Ipv4ConfigType_IPV4_CONFIG_TYPE_OFF:
		return Ipv4ConfigTypeOff
	case devcfgpb.Ipv4ConfigType_IPV4_CONFIG_TYPE_DHCP:
		return Ipv4ConfigTypeDhcp
	case devcfgpb.Ipv4ConfigType_IPV4_CONFIG_TYPE_CUSTOM:
		return Ipv4ConfigTypeCustom
	default:
		return ""
	}
}

func ipv4SettingsFrom(ipv4Settings *devcfgpb.Ipv4ConfigSettings) *Ipv4Settings {
	if ipv4Settings == nil {
		return nil
	}

	return &Ipv4Settings{
		Network: ipv4Settings.GetNetwork(),
		Netmask: ipv4Settings.GetNetmask(),
		Gateway: ipv4Settings.GetGateway(),
	}
}

func ipv4ConfigFrom(ipv4Config *devcfgpb.Ipv4Config) *Ipv4Config {
	if ipv4Config == nil {
		return nil
	}

	return &Ipv4Config{
		Type:     ipv4ConfigTypeFrom(ipv4Config.GetType()),
		Settings: ipv4SettingsFrom(ipv4Config.GetSettings()),
	}
}

func ipv6ConfigTypeFrom(ipv6ConfigType devcfgpb.Ipv6ConfigType) Ipv6ConfigType {
	switch ipv6ConfigType {
	case devcfgpb.Ipv6ConfigType_IPV6_CONFIG_TYPE_OFF:
		return Ipv6ConfigTypeOff
	case devcfgpb.Ipv6ConfigType_IPV6_CONFIG_TYPE_AUTO:
		return Ipv6ConfigTypeAuto
	case devcfgpb.Ipv6ConfigType_IPV6_CONFIG_TYPE_CUSTOM:
		return Ipv6ConfigTypeCustom
	default:
		return ""
	}
}

func ipv6SettingsFrom(ipv6Settings *devcfgpb.Ipv6ConfigSettings) *Ipv6Settings {
	if ipv6Settings == nil {
		return nil
	}

	return &Ipv6Settings{
		Network:      ipv6Settings.GetNetwork(),
		PrefixLength: ipv6Settings.GetPrefixLength(),
		Gateway:      ipv6Settings.GetGateway(),
	}
}

func ipv6ConfigFrom(ipv6Config *devcfgpb.Ipv6Config) *Ipv6Config {
	if ipv6Config == nil {
		return nil
	}

	return &Ipv6Config{
		Type:     ipv6ConfigTypeFrom(ipv6Config.GetType()),
		Settings: ipv6SettingsFrom(ipv6Config.GetSettings()),
	}
}

func ipv6PrivacyFrom(ipv6Privacy devcfgpb.Ipv6Privacy) Ipv6Privacy {
	switch ipv6Privacy {
	case devcfgpb.Ipv6Privacy_IPV6_PRIVACY_DISABLED:
		return Ipv6PrivacyDisabled
	case devcfgpb.Ipv6Privacy_IPV6_PRIVACY_ENABLED:
		return Ipv6PrivacyEnabled
	case devcfgpb.Ipv6Privacy_IPV6_PRIVACY_PREFERRED:
		return Ipv6PrivacyPreferred
	default:
		return ""
	}
}

func securityFrom(secr devcfgpb.Security) Security {
	switch secr {
	case devcfgpb.Security_SECURITY_PSK:
		return SecurityPsk
	case devcfgpb.Security_SECURITY_IEEE8021X:
		return SecurityIeee8021x
	case devcfgpb.Security_SECURITY_NONE:
		return SecurityNone
	case devcfgpb.Security_SECURITY_WEP:
		return SecurityWep
	default:
		return ""
	}
}

func eapFrom(eap devcfgpb.Eap) Eap {
	switch eap {
	case devcfgpb.Eap_EAP_TLS:
		return EapTls
	case devcfgpb.Eap_EAP_TTLS:
		return EapTtls
	case devcfgpb.Eap_EAP_PEAP:
		return EapPeap
	default:
		return ""
	}
}

func caCertTypeFrom(ccType devcfgpb.CaCertType) CaCertType {
	switch ccType {
	case devcfgpb.CaCertType_CA_CERT_TYPE_PEM:
		return CaCertTypePem
	case devcfgpb.CaCertType_CA_CERT_TYPE_DER:
		return CaCertTypeDer
	default:
		return ""
	}
}

func privateKeyTypeFrom(pkType devcfgpb.PrivateKeyType) PrivateKeyType {
	switch pkType {
	case devcfgpb.PrivateKeyType_PRIVATE_KEY_TYPE_PEM:
		return PrivateKeyTypePem
	case devcfgpb.PrivateKeyType_PRIVATE_KEY_TYPE_DER:
		return PrivateKeyTypeDer
	case devcfgpb.PrivateKeyType_PRIVATE_KEY_TYPE_PFX:
		return PrivateKeyTypePfx
	default:
		return ""
	}
}

func privateKeyPassphraseTypeFrom(pkPassphraseType devcfgpb.PrivateKeyPassphraseType) PrivateKeyPassphraseType {
	switch pkPassphraseType {
	case devcfgpb.PrivateKeyPassphraseType_PRIVATE_KEY_PASSPHRASE_TYPE_FSID:
		return PrivateKeyPassphraseTypeFsid
	default:
		return ""
	}
}
func EthernetConfigFrom(cfg *devcfgpb.EthernetService, id uuid.UUID) EthernetService {
	return EthernetService{
		ID:                       id,
		Type:                     ethernetServiceTypeFrom(cfg.GetType()),
		Ipv4Config:               ipv4ConfigFrom(cfg.GetIpv4Config()),
		Ipv6Config:               ipv6ConfigFrom(cfg.GetIpv6Config()),
		Ipv6Privacy:              ipv6PrivacyFrom(cfg.GetIpv6Privacy()),
		Mac:                      cfg.GetMac(),
		Nameservers:              cfg.GetNameservers(),
		SearchDomains:            cfg.GetSearchDomains(),
		Timeservers:              cfg.GetTimeservers(),
		Domain:                   cfg.GetDomain(),
		NetworkName:              cfg.GetName(),
		SSID:                     cfg.GetSsid(),
		Passphrase:               cfg.GetPassphrase(),
		Security:                 securityFrom(cfg.GetSecurity()),
		IsHidden:                 cfg.GetIsHidden(),
		Eap:                      eapFrom(cfg.GetEap()),
		CaCert:                   cfg.GetCaCert(),
		CaCertType:               caCertTypeFrom(cfg.GetCaCertType()),
		PrivateKey:               cfg.GetPrivateKey(),
		PrivateKeyType:           privateKeyTypeFrom(cfg.GetPrivateKeyType()),
		PrivateKeyPassphrase:     cfg.GetPrivateKeyPassphrase(),
		PrivateKeyPassphraseType: privateKeyPassphraseTypeFrom(cfg.GetPrivateKeyPassphraseType()),
		Identity:                 cfg.GetIdentity(),
		AnonymousIdentity:        cfg.GetAnonymousIdentity(),
		SubjectMatch:             cfg.GetSubjectMatch(),
		AltSubjectMatch:          cfg.GetAltSubjectMatch(),
		DomainSuffixMatch:        cfg.GetDomainSuffixMatch(),
		DomainMatch:              cfg.GetDomainMatch(),
		IsPhase2EapBased:         cfg.GetIsPhase_2EapBased(),
	}
}

func NetworkingServiceToProto(ethServ EthernetService) *devcfgpb.EthernetService {
	return &devcfgpb.EthernetService{
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
