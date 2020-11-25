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

type NetworkingServiceType string

const (
	NetworkingServiceTypeEthernet NetworkingServiceType = "Ethernet"
	NetworkingServiceTypeWifi     NetworkingServiceType = "Wifi"
)

type Ipv4ConfigType string

const (
	Ipv4ConfigTypeOff    Ipv4ConfigType = "Off"
	Ipv4ConfigTypeDhcp   Ipv4ConfigType = "Dhcp"
	Ipv4ConfigTypeCustom Ipv4ConfigType = "Custom"
)

type Ipv4ConfigSettings struct {
	Network string `json:"network"`
	Netmask string `json:"netmask"`
	Gateway string `json:"gateway"`
}

type Ipv4Config struct {
	Type     Ipv4ConfigType      `json:"type"`
	Settings *Ipv4ConfigSettings `json:"settings"`
}

type Ipv6ConfigType string

const (
	Ipv6ConfigTypeOff    Ipv6ConfigType = "Off"
	Ipv6ConfigTypeAuto   Ipv6ConfigType = "Auto"
	Ipv6ConfigTypeCustom Ipv6ConfigType = "Custom"
)

type Ipv6ConfigSettings struct {
	Network      string `json:"network"`
	PrefixLength uint64 `json:"prefixLength"`
	Gateway      string `json:"gateway"`
}

type Ipv6Config struct {
	Type     Ipv6ConfigType      `json:"type"`
	Settings *Ipv6ConfigSettings `json:"settings"`
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

type NetworkingService struct {
	ID                       uuid.UUID                `json:"id"`
	OrganizationID           uuid.UUID                `json:"organizationId"`
	Name                     string                   `json:"name"`
	Type                     NetworkingServiceType    `json:"type"`
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

func networkingServiceTypeFrom(cfgType devcfgpb.NetworkingServiceType) NetworkingServiceType {
	switch cfgType {
	case devcfgpb.NetworkingServiceType_NETWORKING_SERVICE_TYPE_ETHERNET:
		return NetworkingServiceTypeEthernet
	case devcfgpb.NetworkingServiceType_NETWORKING_SERVICE_TYPE_WIFI:
		return NetworkingServiceTypeWifi
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

func ipv4ConfigSettingsFrom(ipv4ConfigSettings *devcfgpb.Ipv4ConfigSettings) *Ipv4ConfigSettings {
	if ipv4ConfigSettings == nil {
		return nil
	}

	return &Ipv4ConfigSettings{
		Network: ipv4ConfigSettings.GetNetwork(),
		Netmask: ipv4ConfigSettings.GetNetmask(),
		Gateway: ipv4ConfigSettings.GetGateway(),
	}
}

func ipv4ConfigFrom(ipv4Config *devcfgpb.Ipv4Config) *Ipv4Config {
	if ipv4Config == nil {
		return nil
	}

	return &Ipv4Config{
		Type:     ipv4ConfigTypeFrom(ipv4Config.GetType()),
		Settings: ipv4ConfigSettingsFrom(ipv4Config.GetSettings()),
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

func ipv6ConfigSettingsFrom(ipv6ConfigSettings *devcfgpb.Ipv6ConfigSettings) *Ipv6ConfigSettings {
	if ipv6ConfigSettings == nil {
		return nil
	}

	return &Ipv6ConfigSettings{
		Network:      ipv6ConfigSettings.GetNetwork(),
		PrefixLength: ipv6ConfigSettings.GetPrefixLength(),
		Gateway:      ipv6ConfigSettings.GetGateway(),
	}
}

func ipv6ConfigFrom(ipv6Config *devcfgpb.Ipv6Config) *Ipv6Config {
	if ipv6Config == nil {
		return nil
	}

	return &Ipv6Config{
		Type:     ipv6ConfigTypeFrom(ipv6Config.GetType()),
		Settings: ipv6ConfigSettingsFrom(ipv6Config.GetSettings()),
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

func NetworkingServiceFrom(cfg *devcfgpb.NetworkingService, id uuid.UUID) NetworkingService {
	return NetworkingService{
		ID:                       id,
		Type:                     networkingServiceTypeFrom(cfg.GetType()),
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

func networkingServiceTypeToProto(t NetworkingServiceType) devcfgpb.NetworkingServiceType {
	switch t {
	case NetworkingServiceTypeEthernet:
		return devcfgpb.NetworkingServiceType_NETWORKING_SERVICE_TYPE_ETHERNET
	case NetworkingServiceTypeWifi:
		return devcfgpb.NetworkingServiceType_NETWORKING_SERVICE_TYPE_WIFI
	default:
		return devcfgpb.NetworkingServiceType_NETWORKING_SERVICE_TYPE_UNSPECIFIED
	}
}

func ipv4ConfigTypeToProto(ipv4CfgType Ipv4ConfigType) devcfgpb.Ipv4ConfigType {
	switch ipv4CfgType {
	case Ipv4ConfigTypeOff:
		return devcfgpb.Ipv4ConfigType_IPV4_CONFIG_TYPE_OFF
	case Ipv4ConfigTypeDhcp:
		return devcfgpb.Ipv4ConfigType_IPV4_CONFIG_TYPE_DHCP
	case Ipv4ConfigTypeCustom:
		return devcfgpb.Ipv4ConfigType_IPV4_CONFIG_TYPE_CUSTOM
	default:
		return devcfgpb.Ipv4ConfigType_IPV4_CONFIG_TYPE_UNSPECIFIED
	}
}

func ipv4ConfigSettingsToProto(ipv4CfgSet *Ipv4ConfigSettings) *devcfgpb.Ipv4ConfigSettings {
	if ipv4CfgSet == nil {
		return nil
	}

	return &devcfgpb.Ipv4ConfigSettings{
		Network: ipv4CfgSet.Network,
		Netmask: ipv4CfgSet.Netmask,
		Gateway: ipv4CfgSet.Gateway,
	}
}

func ipv4ConfigToProto(ipv4Cfg *Ipv4Config) *devcfgpb.Ipv4Config {
	if ipv4Cfg == nil {
		return nil
	}

	return &devcfgpb.Ipv4Config{
		Type:     ipv4ConfigTypeToProto(ipv4Cfg.Type),
		Settings: ipv4ConfigSettingsToProto(ipv4Cfg.Settings),
	}
}

func ipv6ConfigTypeToProto(ipv6CfgType Ipv6ConfigType) devcfgpb.Ipv6ConfigType {
	switch ipv6CfgType {
	case Ipv6ConfigTypeOff:
		return devcfgpb.Ipv6ConfigType_IPV6_CONFIG_TYPE_OFF
	case Ipv6ConfigTypeAuto:
		return devcfgpb.Ipv6ConfigType_IPV6_CONFIG_TYPE_AUTO
	case Ipv6ConfigTypeCustom:
		return devcfgpb.Ipv6ConfigType_IPV6_CONFIG_TYPE_CUSTOM
	default:
		return devcfgpb.Ipv6ConfigType_IPV6_CONFIG_TYPE_UNSPECIFIED
	}
}

func ipv6ConfigSettingsToProto(ipv6CfgSet *Ipv6ConfigSettings) *devcfgpb.Ipv6ConfigSettings {
	if ipv6CfgSet == nil {
		return nil
	}

	return &devcfgpb.Ipv6ConfigSettings{
		Network:      ipv6CfgSet.Network,
		PrefixLength: ipv6CfgSet.PrefixLength,
		Gateway:      ipv6CfgSet.Gateway,
	}
}

func ipv6ConfigToProto(ipv6Cfg *Ipv6Config) *devcfgpb.Ipv6Config {
	if ipv6Cfg == nil {
		return nil
	}

	return &devcfgpb.Ipv6Config{
		Type:     ipv6ConfigTypeToProto(ipv6Cfg.Type),
		Settings: ipv6ConfigSettingsToProto(ipv6Cfg.Settings),
	}
}

func ipv6PrivacyToProto(ipv6Privacy Ipv6Privacy) devcfgpb.Ipv6Privacy {
	switch ipv6Privacy {
	case Ipv6PrivacyDisabled:
		return devcfgpb.Ipv6Privacy_IPV6_PRIVACY_DISABLED
	case Ipv6PrivacyEnabled:
		return devcfgpb.Ipv6Privacy_IPV6_PRIVACY_ENABLED
	case Ipv6PrivacyPreferred:
		return devcfgpb.Ipv6Privacy_IPV6_PRIVACY_PREFERRED
	default:
		return devcfgpb.Ipv6Privacy_IPV6_PRIVACY_PREFERRED
	}
}

func securityToProto(security Security) devcfgpb.Security {
	switch security {
	case SecurityPsk:
		return devcfgpb.Security_SECURITY_PSK
	case SecurityIeee8021x:
		return devcfgpb.Security_SECURITY_IEEE8021X
	case SecurityNone:
		return devcfgpb.Security_SECURITY_NONE
	case SecurityWep:
		return devcfgpb.Security_SECURITY_WEP
	default:
		return devcfgpb.Security_SECURITY_UNSPECIFIED
	}
}

func eapToProto(eap Eap) devcfgpb.Eap {
	switch eap {
	case EapTls:
		return devcfgpb.Eap_EAP_TLS
	case EapTtls:
		return devcfgpb.Eap_EAP_TTLS
	case EapPeap:
		return devcfgpb.Eap_EAP_PEAP
	default:
		return devcfgpb.Eap_EAP_UNSPECIFIED
	}
}

func caCertTypeToProto(caCertType CaCertType) devcfgpb.CaCertType {
	switch caCertType {
	case CaCertTypePem:
		return devcfgpb.CaCertType_CA_CERT_TYPE_PEM
	case CaCertTypeDer:
		return devcfgpb.CaCertType_CA_CERT_TYPE_DER
	default:
		return devcfgpb.CaCertType_CA_CERT_TYPE_UNSPECIFIED
	}
}

func privateKeyTypeToProto(pkType PrivateKeyType) devcfgpb.PrivateKeyType {
	switch pkType {
	case PrivateKeyTypePem:
		return devcfgpb.PrivateKeyType_PRIVATE_KEY_TYPE_PEM
	case PrivateKeyTypeDer:
		return devcfgpb.PrivateKeyType_PRIVATE_KEY_TYPE_DER
	case PrivateKeyTypePfx:
		return devcfgpb.PrivateKeyType_PRIVATE_KEY_TYPE_PFX
	default:
		return devcfgpb.PrivateKeyType_PRIVATE_KEY_TYPE_UNSPECIFIED
	}
}

func privateKeyPassphraseTypeToProto(pkPassphraseType PrivateKeyPassphraseType) devcfgpb.PrivateKeyPassphraseType {
	switch pkPassphraseType {
	case PrivateKeyPassphraseTypeFsid:
		return devcfgpb.PrivateKeyPassphraseType_PRIVATE_KEY_PASSPHRASE_TYPE_FSID
	default:
		return devcfgpb.PrivateKeyPassphraseType_PRIVATE_KEY_PASSPHRASE_TYPE_UNSPECIFIED
	}
}

func NetworkingServiceToProto(netServ NetworkingService) *devcfgpb.NetworkingService {
	return &devcfgpb.NetworkingService{
		Type:                     networkingServiceTypeToProto(netServ.Type),
		Ipv4Config:               ipv4ConfigToProto(netServ.Ipv4Config),
		Ipv6Config:               ipv6ConfigToProto(netServ.Ipv6Config),
		Ipv6Privacy:              ipv6PrivacyToProto(netServ.Ipv6Privacy),
		Mac:                      netServ.Mac,
		Nameservers:              netServ.Nameservers,
		SearchDomains:            netServ.SearchDomains,
		Timeservers:              netServ.Timeservers,
		Domain:                   netServ.Domain,
		Name:                     netServ.NetworkName,
		Ssid:                     netServ.SSID,
		Passphrase:               netServ.Passphrase,
		Security:                 securityToProto(netServ.Security),
		IsHidden:                 netServ.IsHidden,
		Eap:                      eapToProto(netServ.Eap),
		CaCert:                   netServ.CaCert,
		CaCertType:               caCertTypeToProto(netServ.CaCertType),
		PrivateKey:               netServ.PrivateKey,
		PrivateKeyType:           privateKeyTypeToProto(netServ.PrivateKeyType),
		PrivateKeyPassphrase:     netServ.PrivateKeyPassphrase,
		PrivateKeyPassphraseType: privateKeyPassphraseTypeToProto(netServ.PrivateKeyPassphraseType),
		Identity:                 netServ.Identity,
		AnonymousIdentity:        netServ.AnonymousIdentity,
		SubjectMatch:             netServ.SubjectMatch,
		AltSubjectMatch:          netServ.AltSubjectMatch,
		DomainSuffixMatch:        netServ.DomainSuffixMatch,
		DomainMatch:              netServ.DomainMatch,
		IsPhase_2EapBased:        netServ.IsPhase2EapBased,
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

type PaginatedNetworkingServices struct {
	Pagination
	Data []NetworkingService `json:"data"`
}

type CreateDeviceResponse struct {
	Device Device    `json:"device"`
	Token  uuid.UUID `json:"token"`
}

type GenerateNATSCredentialsResponse struct {
	Credentials string `json:"credentials"`
}

type RegistrationConfig struct {
	Token uuid.UUID `json:"token"`
	// TODO(Robertvdmaas): NetworkingConfig
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

func PaginatedNetworkingServicesFrom(p database.PaginatedNetworkingServices, data []NetworkingService) PaginatedNetworkingServices {
	return PaginatedNetworkingServices{
		Pagination: PaginationFrom(p.Pagination),
		Data:       data,
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
