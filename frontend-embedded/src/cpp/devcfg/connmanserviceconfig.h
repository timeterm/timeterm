#pragma once

#include <QFile>
#include <QJsonArray>
#include <QJsonObject>
#include <QObject>

class ConnManIpv4Config: public QObject
{
    Q_OBJECT

public:
    class ReadError
    {
    };

    explicit ConnManIpv4Config(QObject *parent = nullptr);

    [[nodiscard]] virtual QString toConnManString() const = 0;

    static ConnManIpv4Config *read(const QJsonObject &cfg);
};

class ConnManIpv4ConfigOff: public ConnManIpv4Config
{
    Q_OBJECT

public:
    explicit ConnManIpv4ConfigOff(QObject *parent = nullptr);

    [[nodiscard]] QString toConnManString() const override;
};

class ConnManIpv4ConfigDhcp: public ConnManIpv4Config
{
    Q_OBJECT

public:
    explicit ConnManIpv4ConfigDhcp(QObject *parent = nullptr);

    [[nodiscard]] QString toConnManString() const override;
};

class ConnManIpv4ConfigCustom: public ConnManIpv4Config
{
    Q_OBJECT
    Q_PROPERTY(QString network WRITE setNetwork READ network NOTIFY networkChanged)
    Q_PROPERTY(QString netmask WRITE setNetmask READ netmask NOTIFY netmaskChanged)
    Q_PROPERTY(QString gateway WRITE setGateway READ gateway NOTIFY gatewayChanged)

public:
    explicit ConnManIpv4ConfigCustom(QObject *parent = nullptr);

    [[nodiscard]] QString toConnManString() const override;

    static ConnManIpv4ConfigCustom *read(const QJsonObject &settings);

    void setNetwork(const QString &network);
    [[nodiscard]] QString network() const;
    void setNetmask(const QString &netmask);
    [[nodiscard]] QString netmask() const;
    void setGateway(const QString &gateway);
    [[nodiscard]] QString gateway() const;

signals:
    void networkChanged();
    void netmaskChanged();
    void gatewayChanged();

private:
    QString m_network;
    QString m_netmask;
    QString m_gateway;
};

class ConnManIpv6Config: public QObject
{
    Q_OBJECT

public:
    class ReadError
    {
    };

    explicit ConnManIpv6Config(QObject *parent = nullptr);

    [[nodiscard]] virtual QString toConnManString() const = 0;

    static ConnManIpv6Config *read(const QJsonObject &cfg);
};

class ConnManIpv6ConfigOff: public ConnManIpv6Config
{
    Q_OBJECT

public:
    explicit ConnManIpv6ConfigOff(QObject *parent = nullptr);

    [[nodiscard]] QString toConnManString() const override;
};

class ConnManIpv6ConfigAuto: public ConnManIpv6Config
{
    Q_OBJECT

public:
    explicit ConnManIpv6ConfigAuto(QObject *parent = nullptr);

    [[nodiscard]] QString toConnManString() const override;
};

class ConnManIpv6ConfigCustom: public ConnManIpv6Config
{
    Q_OBJECT
    Q_PROPERTY(QString network WRITE setNetwork READ network NOTIFY networkChanged)
    Q_PROPERTY(int prefixLength WRITE setPrefixLength READ prefixLength NOTIFY prefixLengthChanged)
    Q_PROPERTY(QString gateway WRITE setGateway READ gateway NOTIFY gatewayChanged)

public:
    explicit ConnManIpv6ConfigCustom(QObject *parent = nullptr);

    [[nodiscard]] QString toConnManString() const override;

    static ConnManIpv6ConfigCustom *read(const QJsonObject &settings);

    void setNetwork(const QString &network);
    [[nodiscard]] QString network() const;
    void setPrefixLength(int prefixLength);
    [[nodiscard]] int prefixLength() const;
    void setGateway(const QString &gateway);
    [[nodiscard]] QString gateway() const;

signals:
    void networkChanged();
    void prefixLengthChanged();
    void gatewayChanged();

private:
    QString m_network;
    int m_prefixLength = 0;
    QString m_gateway;
};

/// See connman-service.config(5) (service-name.config(5))
class ConnManServiceConfig: public QObject
{
    Q_OBJECT
    Q_PROPERTY(QString serviceName READ serviceName WRITE setServiceName NOTIFY serviceNameChanged)
    Q_PROPERTY(ServiceType type READ type WRITE setType NOTIFY typeChanged)
    Q_PROPERTY(ConnManIpv4Config *ipv4Config READ ipv4Config WRITE setIpv4Config NOTIFY ipv4ConfigChanged)
    Q_PROPERTY(ConnManIpv6Config *ipv6Config READ ipv6Config WRITE setIpv6Config NOTIFY ipv6ConfigChanged)
    Q_PROPERTY(Ipv6Privacy ipv6Privacy READ ipv6Privacy WRITE setIpv6Privacy NOTIFY ipv6PrivacyChanged)
    Q_PROPERTY(QString mac READ mac WRITE setMac NOTIFY macChanged)
    Q_PROPERTY(QString deviceName READ deviceName WRITE setDeviceName NOTIFY deviceNameChanged)
    Q_PROPERTY(QStringList nameservers READ nameservers WRITE setNameservers NOTIFY nameserversChanged)
    Q_PROPERTY(QStringList searchDomains READ searchDomains WRITE setSearchDomains NOTIFY searchDomainsChanged)
    Q_PROPERTY(QStringList timeservers READ timeservers WRITE setTimeservers NOTIFY timeserversChanged)
    Q_PROPERTY(QString domain READ domain WRITE setDomain NOTIFY domainChanged)
    Q_PROPERTY(QString name READ name WRITE setName NOTIFY nameChanged)
    Q_PROPERTY(QString ssid READ ssid WRITE setSsid NOTIFY ssidChanged)
    Q_PROPERTY(QString passphrase READ passphrase WRITE setPassphrase NOTIFY passphraseChanged)
    Q_PROPERTY(Security security READ security WRITE setSecurity NOTIFY securityChanged)
    Q_PROPERTY(bool isHidden READ isHidden WRITE setIsHidden NOTIFY isHiddenChanged)
    Q_PROPERTY(EapType eap READ eap WRITE setEap NOTIFY eapChanged)
    Q_PROPERTY(QByteArray caCert READ caCert WRITE setCaCert NOTIFY caCertChanged)
    Q_PROPERTY(CaCertType caCertType READ caCertType WRITE setCaCertType NOTIFY caCertTypeChanged)
    Q_PROPERTY(QByteArray privateKey READ privateKey WRITE setPrivateKey NOTIFY privateKeyChanged)
    Q_PROPERTY(PrivateKeyType privateKeyType READ privateKeyType WRITE setPrivateKeyType NOTIFY privateKeyTypeChanged)
    Q_PROPERTY(QString privateKeyPassphrase READ privateKeyPassphrase WRITE setPrivateKeyPassphrase NOTIFY privateKeyPassphraseChanged)
    Q_PROPERTY(PrivateKeyPassphraseType privateKeyPassphraseType READ privateKeyPassphraseType WRITE setPrivateKeyPassphraseType NOTIFY privateKeyPassphraseTypeChanged)
    Q_PROPERTY(QString identity READ identity WRITE setIdentity NOTIFY identityChanged)
    Q_PROPERTY(QString anonymousIdentity READ anonymousIdentity WRITE setAnonymousIdentity NOTIFY anonymousIdentityChanged)
    Q_PROPERTY(QString subjectMatch READ subjectMatch WRITE setSubjectMatch NOTIFY subjectMatchChanged)
    Q_PROPERTY(QString altSubjectMatch READ altSubjectMatch WRITE setAltSubjectMatch NOTIFY altSubjectMatchChanged)
    Q_PROPERTY(QString domainSuffixMatch READ domainSuffixMatch WRITE setDomainSuffixMatch NOTIFY domainSuffixMatchChanged)
    Q_PROPERTY(QString domainMatch READ domainMatch WRITE setDomainMatch NOTIFY domainMatchChanged)
    Q_PROPERTY(bool isPhase2EapBased READ isPhase2EapBased WRITE setIsPhase2EapBased NOTIFY isPhase2EapBasedChanged)

public:
    enum ServiceType
    {
        ServiceTypeUndefined,
        ServiceTypeEthernet,
        ServiceTypeWifi,
    };
    Q_ENUM(ServiceType)

    enum Ipv6Privacy
    {
        Ipv6PrivacyUndefined,
        Ipv6PrivacyDisabled,
        Ipv6PrivacyEnabled,
        Ipv6PrivacyPreferred,
    };
    Q_ENUM(Ipv6Privacy)

    enum Security
    {
        SecurityUndefined,
        SecurityPsk,       ///< WPA/WPA2 PSK
        SecurityIeee8021x, ///< WPA EAP
        SecurityNone,
        SecurityWep
    };
    Q_ENUM(Security)

    enum EapType
    {
        EapTypeUndefined,
        EapTypeTls,
        EapTypeTtls,
        EapTypePeap,
    };
    Q_ENUM(EapType)

    enum PrivateKeyPassphraseType
    {
        PrivateKeyPassphraseTypeUndefined,
        PrivateKeyPassphraseTypeFsid,
    };
    Q_ENUM(PrivateKeyPassphraseType)

    enum Phase2Type
    {
        Phase2TypeUndefined,
        Phase2TypeMschapV2,
        Phase2TypeGtc,
    };
    Q_ENUM(Phase2Type)

    enum CaCertType
    {
        CaCertTypeUndefined,
        CaCertTypePem,
        CaCertTypeDer,
    };
    Q_ENUM(CaCertType)

    enum PrivateKeyType
    {
        PrivateKeyTypeUndefined,
        PrivateKeyTypePem,
        PrivateKeyTypeDer,
        PrivateKeyTypePfx,
    };
    Q_ENUM(PrivateKeyType)

    explicit ConnManServiceConfig(QObject *parent = nullptr);

    void setServiceName(const QString &serviceName);
    [[nodiscard]] QString serviceName() const;
    void setType(ServiceType type);
    [[nodiscard]] ServiceType type() const;
    void setIpv4Config(ConnManIpv4Config *config);
    [[nodiscard]] ConnManIpv4Config *ipv4Config() const;
    void setIpv6Config(ConnManIpv6Config *config);
    [[nodiscard]] ConnManIpv6Config *ipv6Config() const;
    void setIpv6Privacy(Ipv6Privacy privacy);
    [[nodiscard]] Ipv6Privacy ipv6Privacy() const;
    void setMac(const QString &mac);
    [[nodiscard]] QString mac() const;
    void setDeviceName(const QString &deviceName);
    [[nodiscard]] QString deviceName() const;
    void setNameservers(const QStringList &nameservers);
    [[nodiscard]] QStringList nameservers() const;
    void setSearchDomains(const QStringList &searchDomains);
    [[nodiscard]] QStringList searchDomains() const;
    void setTimeservers(const QStringList &timeservers);
    [[nodiscard]] QStringList timeservers() const;
    void setDomain(const QString &domain);
    [[nodiscard]] QString domain() const;

    void setName(const QString &name);
    [[nodiscard]] QString name() const;
    void setSsid(const QString &ssid);
    [[nodiscard]] QString ssid() const;
    void setPassphrase(const QString &passphrase);
    [[nodiscard]] QString passphrase() const;
    void setSecurity(Security security);
    [[nodiscard]] Security security() const;
    void setIsHidden(bool isHidden);
    [[nodiscard]] bool isHidden() const;

    void setEap(EapType eap);
    [[nodiscard]] EapType eap() const;
    void setCaCert(const QByteArray &caCert);
    [[nodiscard]] QByteArray caCert() const;
    void setCaCertType(CaCertType caCertType);
    [[nodiscard]] CaCertType caCertType() const;
    void setPrivateKey(const QByteArray &privateKey);
    [[nodiscard]] QByteArray privateKey() const;
    void setPrivateKeyType(PrivateKeyType privateKeyType);
    [[nodiscard]] PrivateKeyType privateKeyType() const;
    void setPrivateKeyPassphrase(const QString &privateKeyPassphrase);
    [[nodiscard]] QString privateKeyPassphrase() const;
    void setPrivateKeyPassphraseType(PrivateKeyPassphraseType privateKeyPassphraseType);
    [[nodiscard]] PrivateKeyPassphraseType privateKeyPassphraseType() const;
    void setIdentity(const QString &identity);
    [[nodiscard]] QString identity() const;
    void setAnonymousIdentity(const QString &anonymousIdentity);
    [[nodiscard]] QString anonymousIdentity() const;
    void setSubjectMatch(const QString &subjectMatch);
    [[nodiscard]] QString subjectMatch() const;
    void setAltSubjectMatch(const QString &altSubjectMatch);
    [[nodiscard]] QString altSubjectMatch() const;
    void setDomainSuffixMatch(const QString &domainSuffixMatch);
    [[nodiscard]] QString domainSuffixMatch() const;
    void setDomainMatch(const QString &domainMatch);
    [[nodiscard]] QString domainMatch() const;
    void setPhase2Type(Phase2Type phase2Type);
    [[nodiscard]] Phase2Type phase2Type() const;
    void setIsPhase2EapBased(bool isPhase2EapBased);
    [[nodiscard]] bool isPhase2EapBased() const;

    enum ReadError
    {
        ReadErrorNoError = 0,
    };

    void read(const QJsonObject &obj, ReadError *error = nullptr);
    void saveCerts(QFile::FileError *error = nullptr);
    void saveConnManConf(QFile::FileError *error = nullptr);

    static ServiceType readServiceType(const QString &t);
    static Ipv6Privacy readIpv6Privacy(const QString &p);
    static Security readSecurity(const QString &s);
    static EapType readEapType(const QString &t);
    static CaCertType readCaCertType(const QString &t);
    static PrivateKeyType readPrivateKeyType(const QString &t);
    static PrivateKeyPassphraseType readPrivateKeyPassphraseType(const QString &t);
    static Phase2Type readPhase2Type(const QString &t);

    static QString serviceTypeToConnManString(ServiceType t);
    static QString ipv6PrivacyToConnManString(Ipv6Privacy p);
    static QString securityToConnManString(Security s);
    static QString eapTypeToConnManString(EapType t);
    static QString privateKeyPassphraseTypeToConnManString(PrivateKeyPassphraseType t);
    static QString phase2TypeToConnManString(Phase2Type t);

signals:
    void serviceNameChanged();
    void typeChanged();
    void ipv4ConfigChanged();
    void ipv6ConfigChanged();
    void ipv6PrivacyChanged();
    void macChanged();
    void deviceNameChanged();
    void nameserversChanged();
    void searchDomainsChanged();
    void timeserversChanged();
    void domainChanged();

    void nameChanged();
    void ssidChanged();
    void passphraseChanged();
    void securityChanged();
    void isHiddenChanged();

    void eapChanged();
    void caCertChanged();
    void caCertTypeChanged();
    void privateKeyChanged();
    void privateKeyTypeChanged();
    void privateKeyPassphraseChanged();
    void privateKeyPassphraseTypeChanged();
    void identityChanged();
    void anonymousIdentityChanged();
    void subjectMatchChanged();
    void altSubjectMatchChanged();
    void domainSuffixMatchChanged();
    void domainMatchChanged();
    void phase2TypeChanged();
    void isPhase2EapBasedChanged();

private:
    /// Mandatory. Interpolated in the [service_*] config section name.
    /// Named 'timeterm' by default.
    QString m_serviceName = "timeterm";

    /// Mandatory. Other types than Ethernet or Wifi are not supported.
    ServiceType m_type = ServiceTypeUndefined;
    /// IPv4 settings for the service. If set to off, IPv4 won't be used.
    /// If set to Dhcp, DHCP will be used to obtain the network settings.
    /// netmask can be specified as length of the mask rather than the mask itself.
    /// The gateway can be omitted when using a static IP.
    ConnManIpv4Config *m_ipv4Config = nullptr;
    /// IPv6 settings for the service. If set to Off, IPv6 won't be used.
    /// If set to Auto, settings will be obtained from the network.
    ConnManIpv6Config *m_ipv6Config = nullptr;
    /// IPv6 privacy settings as per RFC3041.
    Ipv6Privacy m_ipv6Privacy = Ipv6PrivacyUndefined;
    /// MAC address of the interface to be used.
    /// If not specified, the first found interface is used.
    /// Must be in format ab:cd:ef:01:23:45.
    QString m_mac;
    /// Device name the interface to be used, e.g. eth0.
    /// mac takes preference over deviceName.
    QString m_deviceName;
    /// Comma separated list of nameservers.
    QStringList m_nameservers;
    /// Comma separated list of DNS search domains.
    QStringList m_searchDomains;
    /// Comma separated list of timeservers.
    QStringList m_timeservers;
    /// Domain name to be used.
    QString m_domain;

    /// A string representation of an network SSID.
    /// If the ssid field is present, the name field is ignored.
    /// If the ssid field is not present, this field is mandatory.
    QString m_name;
    /// SSID: A hexadecimal representation of an 802.11 SSID.
    /// Use this format to encode special characters including starting or ending spaces.
    QString m_ssid;
    /// RSN/WPA/WPA2 Passphrase.
    QString m_passphrase;
    /// The security type of the network.
    /// Possible values are Psk (WPA/WPA2 PSK), Ieee8021x (WPA EAP), None and Wep.
    /// When not set, the default value is Ieee8021x if an EAP type is configured,
    /// Psk if a passphrase is present and None otherwise.
    Security m_security = SecurityUndefined;
    /// If set to true, then this AP is hidden.
    /// If missing or set to false, then AP is not hidden.
    bool m_isHidden = false;

    /// EAP type to use.
    /// Only Tls, Ttls and Peap are supported.
    EapType m_eap = EapTypeUndefined;
    /// CA certificate.
    /// Only PEM and DER formats are supported.
    QByteArray m_caCert;
    CaCertType m_caCertType = CaCertTypeUndefined;
    /// Private key.
    /// Only PEM, DER and PFX formats are supported.
    QByteArray m_privateKey;
    PrivateKeyType m_privateKeyType = PrivateKeyTypeUndefined;
    /// Passphrase of the private key.
    QString m_privateKeyPassphrase;
    /// If specified, use the private key's FSID as the passphrase, and ignore the privateKeyPassphrase field.
    PrivateKeyPassphraseType m_privateKeyPassphraseType = PrivateKeyPassphraseTypeUndefined;
    /// Identity string for EAP.
    QString m_identity;
    /// Anonymous identity string for EAP.
    QString m_anonymousIdentity;
    /// Substring to be matched against the subject of the authentication server certificate for EAP.
    QString m_subjectMatch;
    /// Semicolon separated string of entries to be matched against
    /// the alternative subject name of the authentication server certificate for EAP.
    QString m_altSubjectMatch;
    /// Constraint for server domain name.
    /// If set, this FQDN is used as a suffix match requirement for the authentication server certificate for EAP.
    QString m_domainSuffixMatch;
    /// This FQDN is used as a full match requirement for the authentication server certificate for EAP.
    QString m_domainMatch;
    /// Inner authentication type with for eap = Tls or eap = Ttls.
    /// Set phase2EapBased to true to indicate usage of EAP-based authentication method (should only be used with eap = Ttls).
    Phase2Type m_phase2Type = Phase2TypeUndefined;
    bool m_isPhase2EapBased = false;
};

Q_DECLARE_METATYPE(ConnManServiceConfig *)
