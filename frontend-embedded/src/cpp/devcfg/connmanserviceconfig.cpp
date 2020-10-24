#include "connmanserviceconfig.h"

ConnManIpv4Config::ConnManIpv4Config(QObject *parent)
    : QObject(parent)
{}

ConnManIpv4Config *ConnManIpv4Config::read(const QJsonObject &cfg)
{
    if (!cfg.contains("type") || !cfg["type"].isString())
        return nullptr; // TODO: set error (invalid type)
    auto type = cfg["type"].toString();

    if (type == "Off")
        return new ConnManIpv4ConfigOff();
    if (type == "Dhcp")
        return new ConnManIpv4ConfigDhcp();
    if (type == "Custom") {
        if (!cfg.contains("settings") || !cfg["settings"].isObject())
            return nullptr; // TODO: set error (missing settings)
        return ConnManIpv4ConfigCustom::read(cfg["settings"].toObject());
    }

    return nullptr; // TODO: set error (unknown type)
}

ConnManIpv4ConfigOff::ConnManIpv4ConfigOff(QObject *parent)
    : ConnManIpv4Config(parent)
{}

QString ConnManIpv4ConfigOff::toString() const
{
    return QStringLiteral("off");
}

ConnManIpv4ConfigDhcp::ConnManIpv4ConfigDhcp(QObject *parent)
    : ConnManIpv4Config(parent)
{}

QString ConnManIpv4ConfigDhcp::toString() const
{
    return QStringLiteral("dhcp");
}

ConnManIpv4ConfigCustom::ConnManIpv4ConfigCustom(QObject *parent)
    : ConnManIpv4Config(parent)
{}

ConnManIpv4ConfigCustom *ConnManIpv4ConfigCustom::read(const QJsonObject &settings)
{
    auto config = new ConnManIpv4ConfigCustom();

    if (settings.contains("network") && settings["network"].isString())
        config->setNetwork(settings["network"].toString());
    if (settings.contains("netmask") && settings["netmask"].isString())
        config->setNetmask(settings["netmask"].toString());
    if (settings.contains("gateway") && settings["gateway"].isString())
        config->setGateway(settings["gateway"].toString());

    return config;
}

QString ConnManIpv4ConfigCustom::toString() const
{
    return QStringLiteral("%1/%2/%3").arg(m_network, m_netmask, m_gateway);
}

void ConnManIpv4ConfigCustom::setNetwork(const QString &network)
{
    if (network != m_network) {
        m_network = network;
        emit networkChanged();
    }
}

QString ConnManIpv4ConfigCustom::network() const
{
    return m_network;
}

void ConnManIpv4ConfigCustom::setNetmask(const QString &netmask)
{
    if (netmask != m_netmask) {
        m_netmask = netmask;
        emit netmaskChanged();
    }
}

QString ConnManIpv4ConfigCustom::netmask() const
{
    return m_netmask;
}

void ConnManIpv4ConfigCustom::setGateway(const QString &gateway)
{
    if (gateway != m_gateway) {
        m_gateway = gateway;
        emit gatewayChanged();
    }
}

QString ConnManIpv4ConfigCustom::gateway() const
{
    return m_gateway;
}

ConnManIpv6Config::ConnManIpv6Config(QObject *parent)
    : QObject(parent)
{}

ConnManIpv6Config *ConnManIpv6Config::read(const QJsonObject &cfg)
{
    if (!cfg.contains("type") || !cfg["type"].isString())
        return nullptr; // TODO: set error (invalid type)
    auto type = cfg["type"].toString();

    if (type == "Off")
        return new ConnManIpv6ConfigOff();
    if (type == "Auto")
        return new ConnManIpv6ConfigAuto();
    if (type == "Custom") {
        if (!cfg.contains("settings") || !cfg["settings"].isObject())
            return nullptr; // TODO: set error (missing settings)
        return ConnManIpv6ConfigCustom::read(cfg);
    }

    return nullptr; // TODO: set error (unknown type)
}

ConnManIpv6ConfigOff::ConnManIpv6ConfigOff(QObject *parent)
    : ConnManIpv6Config(parent)
{}

QString ConnManIpv6ConfigOff::toString() const
{
    return QStringLiteral("off");
}

ConnManIpv6ConfigAuto::ConnManIpv6ConfigAuto(QObject *parent)
    : ConnManIpv6Config(parent)
{}

QString ConnManIpv6ConfigAuto::toString() const
{
    return QStringLiteral("auto");
}

ConnManIpv6ConfigCustom::ConnManIpv6ConfigCustom(QObject *parent)
    : ConnManIpv6Config(parent)
{}

QString ConnManIpv6ConfigCustom::toString() const
{
    return QStringLiteral("%1/%2/%3").arg(m_network).arg(m_prefixLength).arg(m_gateway);
}

ConnManIpv6ConfigCustom *ConnManIpv6ConfigCustom::read(const QJsonObject &settings)
{
    auto config = new ConnManIpv6ConfigCustom();

    if (settings.contains("network") && settings["network"].isString())
        config->setNetwork(settings["network"].toString());
    if (settings.contains("prefixLength") && settings["prefixLength"].isDouble())
        config->setPrefixLength(settings["prefixLength"].toInt());
    if (settings.contains("gateway") && settings["gateway"].isString())
        config->setGateway(settings["gateway"].toString());

    return config;
}

void ConnManIpv6ConfigCustom::setNetwork(const QString &network)
{
    if (network != m_network) {
        m_network = network;
        emit networkChanged();
    }
}

QString ConnManIpv6ConfigCustom::network() const
{
    return m_network;
}

void ConnManIpv6ConfigCustom::setPrefixLength(int prefixLength)
{
    if (prefixLength != m_prefixLength) {
        m_prefixLength = prefixLength;
        emit prefixLengthChanged();
    }
}

int ConnManIpv6ConfigCustom::prefixLength() const
{
    return m_prefixLength;
}

void ConnManIpv6ConfigCustom::setGateway(const QString &gateway)
{
    if (gateway != m_gateway) {
        m_gateway = gateway;
        emit gatewayChanged();
    }
}

QString ConnManIpv6ConfigCustom::gateway() const
{
    return m_gateway;
}

ConnManServiceConfig::ConnManServiceConfig(QObject *parent)
    : QObject(parent)
{}

QStringList jsonArrayToQStringList(const QJsonArray &a)
{
    auto l = QStringList();
    l.reserve(a.size());

    for (const auto &it : a) {
        if (it.isString())
            l.append(it.toString());
    }

    return l;
}

void ConnManServiceConfig::read(QJsonObject &obj, ConnManServiceConfig::ReadError *error)
{
    if (obj.contains("type") && obj["type"].isString())
        setType(readServiceType(obj["type"].toString()));
    if (obj.contains("ipv4Config") && obj["ipv4Config"].isObject())
        setIpv4Config(ConnManIpv4Config::read(obj["ipv4Config"].toObject()));
    if (obj.contains("ipv6Config") && obj["ipv6Config"].isObject())
        setIpv6Config(ConnManIpv6Config::read(obj["ipv6Config"].toObject()));
    if (obj.contains("ipv6Privacy") && obj["ipv6Privacy"].isString())
        setIpv6Privacy(readIpv6Privacy(obj["ipv6Privacy"].toString()));
    if (obj.contains("mac") && obj["mac"].isString())
        setMac(obj["mac"].toString());
    if (obj.contains("deviceName") && obj["deviceName"].isString())
        setDeviceName(obj["deviceName"].toString());
    if (obj.contains("searchDomains") && obj["searchDomains"].isArray())
        setSearchDomains(jsonArrayToQStringList(obj["searchDomains"].toArray()));
    if (obj.contains("timeServers") && obj["timeServers"].isArray())
        setTimeServers(jsonArrayToQStringList(obj["timeServers"].toArray()));
    if (obj.contains("domain") && obj["domain"].isString())
        setDomain(obj["domain"].toString());

    if (obj.contains("name") && obj["name"].isString())
        setName(obj["name"].toString());
    if (obj.contains("ssid") && obj["ssid"].isString())
        setSsid(obj["ssid"].toString());
    if (obj.contains("passphrase") && obj["passphrase"].isString())
        setPassphrase(obj["passphrase"].toString());
    if (obj.contains("security") && obj["security"].isString())
        setSecurity(readSecurity(obj["security"].toString()));
    if (obj.contains("isHidden") && obj["isHidden"].isBool())
        setIsHidden(obj["isHidden"].toBool());

    if (obj.contains("eap") && obj["eap"].isString())
        setEap(readEapType(obj["eap"].toString()));
    if (obj.contains("caCert") && obj["caCert"].isString())
        setCaCert(QByteArray::fromBase64(obj["caCert"].toString().toLocal8Bit()));
    if (obj.contains("caCertType") && obj["caCertType"].isString())
        setCaCertType(readCaCertType(obj["caCertType"].toString()));
    if (obj.contains("privateKey") && obj["privateKey"].isString())
        setPrivateKey(QByteArray::fromBase64(obj["privateKey"].toString().toLocal8Bit()));
    if (obj.contains("privateKeyType") && obj["privateKeyType"].isString())
        setPrivateKeyType(readPrivateKeyType(obj["privateKeyType"].toString()));
    if (obj.contains("privateKeyPassphrase") && obj["privateKeyPassphrase"].isString())
        setPrivateKeyPassphrase(obj["privateKeyPassphrase"].toString());
    if (obj.contains("identity") && obj["identity"].isString())
        setIdentity(obj["identity"].toString());
    if (obj.contains("anonymousIdentity") && obj["anonymousIdentity"].isString())
        setAnonymousIdentity(obj["anonymousIdentity"].toString());
    if (obj.contains("subjectMatch") && obj["subjectMatch"].isString())
        setSubjectMatch(obj["subjectMatch"].toString());
    if (obj.contains("altSubjectMatch") && obj["altSubjectMatch"].isString())
        setAltSubjectMatch(obj["altSubjectMatch"].toString());
    if (obj.contains("domainSubjectMatch") && obj["domainSubjectMatch"].isString())
        setDomainSubjectMatch(obj["domainSubjectMatch"].toString());
    if (obj.contains("domainMatch") && obj["domainMatch"].isString())
        setDomainMatch(obj["domainMatch"].toString());
    if (obj.contains("phase2Type") && obj["phase2Type"].isString())
        setPhase2Type(readPhase2Type(obj["phase2Type"].toString()));
    if (obj.contains("isPhase2EapBased") && obj["isPhase2EapBased"].isString())
        setIsPhase2EapBased(obj["isPhase2EapBased"].toBool());
}

ConnManServiceConfig::ServiceType ConnManServiceConfig::readServiceType(const QString &t)
{
    if (t == "Ethernet")
        return ServiceTypeEthernet;
    if (t == "Wifi")
        return ServiceTypeWifi;
    return ServiceTypeUndefined;
}

ConnManServiceConfig::Ipv6Privacy ConnManServiceConfig::readIpv6Privacy(const QString &p)
{
    if (p == "Disabled")
        return Ipv6PrivacyDisabled;
    if (p == "Enabled")
        return Ipv6PrivacyEnabled;
    if (p == "Preferred")
        return Ipv6PrivacyPreferred;
    return Ipv6PrivacyUndefined;
}

ConnManServiceConfig::Security ConnManServiceConfig::readSecurity(const QString &s)
{
    if (s == "Psk")
        return SecurityPsk;
    if (s == "Ieee8021x")
        return SecurityIeee8021x;
    if (s == "None")
        return SecurityNone;
    if (s == "Wep")
        return SecurityWep;
    return SecurityUndefined;
}

ConnManServiceConfig::EapType ConnManServiceConfig::readEapType(const QString &t)
{
    if (t == "Tls")
        return EapTypeTls;
    if (t == "Ttls")
        return EapTypeTtls;
    if (t == "Peap")
        return EapTypePeap;
    return EapTypeUndefined;
}

ConnManServiceConfig::CaCertType ConnManServiceConfig::readCaCertType(const QString &t)
{
    if (t == "Pem")
        return CaCertTypePem;
    if (t == "Der")
        return CaCertTypeDer;
    return CaCertTypeUndefined;
}

ConnManServiceConfig::PrivateKeyType ConnManServiceConfig::readPrivateKeyType(const QString &t)
{
    if (t == "Pem")
        return PrivateKeyTypePem;
    if (t == "Der")
        return PrivateKeyTypeDer;
    if (t == "Pfx")
        return PrivateKeyTypePfx;
    return PrivateKeyTypeUndefined;
}

ConnManServiceConfig::PrivateKeyPassphraseType ConnManServiceConfig::readPrivateKeyPassphraseType(const QString &t)
{
    if (t == "Fsid")
        return PrivateKeyPassphraseTypeFsid;
    return PrivateKeyPassphraseTypeUndefined;
}

ConnManServiceConfig::Phase2Type ConnManServiceConfig::readPhase2Type(const QString &t)
{
    if (t == "MschapV2")
        return Phase2TypeMschapV2;
    if (t == "Gtc")
        return Phase2TypeGtc;
    return Phase2TypeUndefined;
}

void ConnManServiceConfig::setType(ConnManServiceConfig::ServiceType type)
{
    if (type != m_type) {
        m_type = type;
        emit typeChanged();
    }
}

ConnManServiceConfig::ServiceType ConnManServiceConfig::type() const
{
    return m_type;
}

void ConnManServiceConfig::setIpv4Config(ConnManIpv4Config *config)
{
    if (config != m_ipv4Config) {
        if (config != nullptr)
            config->setParent(this);
        m_ipv4Config = config;
        emit ipv4ConfigChanged();
    }
}

ConnManIpv4Config *ConnManServiceConfig::ipv4Config() const
{
    return m_ipv4Config;
}

void ConnManServiceConfig::setIpv6Config(ConnManIpv6Config *config)
{
    if (config != m_ipv6Config) {
        if (config != nullptr)
            config->setParent(this);
        m_ipv6Config = config;
        emit ipv6ConfigChanged();
    }
}

ConnManIpv6Config *ConnManServiceConfig::ipv6Config() const
{
    return m_ipv6Config;
}

void ConnManServiceConfig::setIpv6Privacy(Ipv6Privacy privacy)
{
    if (privacy != m_ipv6Privacy) {
        m_ipv6Privacy = privacy;
        emit ipv6PrivacyChanged();
    }
}

ConnManServiceConfig::Ipv6Privacy ConnManServiceConfig::ipv6Privacy() const
{
    return m_ipv6Privacy;
}

void ConnManServiceConfig::setMac(const QString &mac)
{
    if (mac != m_mac) {
        m_mac = mac;
        emit macChanged();
    }
}

QString ConnManServiceConfig::mac() const
{
    return m_mac;
}

void ConnManServiceConfig::setDeviceName(const QString &deviceName)
{
    if (deviceName != m_deviceName) {
        m_deviceName = deviceName;
        emit deviceNameChanged();
    }
}

QString ConnManServiceConfig::deviceName() const
{
    return m_deviceName;
}

void ConnManServiceConfig::setNameservers(const QStringList &nameservers)
{
    if (nameservers != m_nameservers) {
        m_nameservers = nameservers;
        emit nameserversChanged();
    }
}

QStringList ConnManServiceConfig::nameservers() const
{
    return m_nameservers;
}

void ConnManServiceConfig::setSearchDomains(const QStringList &searchDomains)
{
    if (searchDomains != m_searchDomains) {
        m_searchDomains = searchDomains;
        emit searchDomainsChanged();
    }
}

QStringList ConnManServiceConfig::searchDomains() const
{
    return m_searchDomains;
}

void ConnManServiceConfig::setTimeServers(const QStringList &timeServers)
{
    if (timeServers != m_timeServers) {
        m_timeServers = timeServers;
        emit timeServersChanged();
    }
}

QStringList ConnManServiceConfig::timeServers() const
{
    return m_timeServers;
}

void ConnManServiceConfig::setDomain(const QString &domain)
{
    if (domain != m_domain) {
        m_domain = domain;
        emit domainChanged();
    }
}

QString ConnManServiceConfig::domain() const
{
    return m_domain;
}

void ConnManServiceConfig::setName(const QString &name)
{
    if (name != m_name) {
        m_name = name;
        emit nameChanged();
    }
}

QString ConnManServiceConfig::name() const
{
    return m_name;
}

void ConnManServiceConfig::setSsid(const QString &ssid)
{
    if (ssid != m_ssid) {
        m_ssid = ssid;
        emit ssidChanged();
    }
}

QString ConnManServiceConfig::ssid() const
{
    return m_ssid;
}

void ConnManServiceConfig::setPassphrase(const QString &passphrase)
{
    if (passphrase != m_passphrase) {
        m_passphrase = passphrase;
        emit passphraseChanged();
    }
}

QString ConnManServiceConfig::passphrase() const
{
    return m_passphrase;
}

void ConnManServiceConfig::setSecurity(ConnManServiceConfig::Security security)
{
    if (security != m_security) {
        m_security = security;
        emit securityChanged();
    }
}

ConnManServiceConfig::Security ConnManServiceConfig::security() const
{
    return m_security;
}

void ConnManServiceConfig::setIsHidden(bool isHidden)
{
    if (isHidden != m_isHidden) {
        m_isHidden = isHidden;
        emit isHiddenChanged();
    }
}

bool ConnManServiceConfig::isHidden() const
{
    return m_isHidden;
}

void ConnManServiceConfig::setEap(ConnManServiceConfig::EapType eap)
{
    if (eap != m_eap) {
        m_eap = eap;
        emit eapChanged();
    }
}

ConnManServiceConfig::EapType ConnManServiceConfig::eap() const
{
    return m_eap;
}

void ConnManServiceConfig::setCaCert(const QByteArray &caCert)
{
    if (caCert != m_caCert) {
        m_caCert = caCert;
        emit caCertChanged();
    }
}

QByteArray ConnManServiceConfig::caCert() const
{
    return m_caCert;
}

void ConnManServiceConfig::setCaCertType(ConnManServiceConfig::CaCertType caCertType)
{
    if (caCertType != m_caCertType) {
        m_caCertType = caCertType;
        emit caCertTypeChanged();
    }
}

ConnManServiceConfig::CaCertType ConnManServiceConfig::caCertType() const
{
    return m_caCertType;
}

void ConnManServiceConfig::setPrivateKey(const QByteArray &privateKey)
{
    if (privateKey != m_privateKey) {
        m_privateKey = privateKey;
        emit privateKeyChanged();
    }
}

QByteArray ConnManServiceConfig::privateKey() const
{
    return m_privateKey;
}

void ConnManServiceConfig::setPrivateKeyType(PrivateKeyType privateKeyType)
{
    if (privateKeyType != m_privateKeyType) {
        m_privateKeyType = privateKeyType;
        emit privateKeyTypeChanged();
    }
}

ConnManServiceConfig::PrivateKeyType ConnManServiceConfig::privateKeyType() const
{
    return m_privateKeyType;
}

void ConnManServiceConfig::setPrivateKeyPassphrase(const QString &privateKeyPassphrase)
{
    if (privateKeyPassphrase != m_privateKeyPassphrase) {
        m_privateKeyPassphrase = privateKeyPassphrase;
        emit privateKeyPassphraseChanged();
    }
}

QString ConnManServiceConfig::privateKeyPassphrase() const
{
    return m_privateKeyPassphrase;
}

void ConnManServiceConfig::setPrivateKeyPassphraseType(ConnManServiceConfig::PrivateKeyPassphraseType privateKeyPassphraseType)
{
    if (privateKeyPassphraseType != m_privateKeyPassphraseType) {
        m_privateKeyPassphraseType = privateKeyPassphraseType;
        emit privateKeyPassphraseTypeChanged();
    }
}

ConnManServiceConfig::PrivateKeyPassphraseType ConnManServiceConfig::privateKeyPassphraseType() const
{
    return m_privateKeyPassphraseType;
}

void ConnManServiceConfig::setIdentity(const QString &identity)
{
    if (identity != m_identity) {
        m_identity = identity;
        emit identityChanged();
    }
}

QString ConnManServiceConfig::identity() const
{
    return m_identity;
}

void ConnManServiceConfig::setAnonymousIdentity(const QString &anonymousIdentity)
{
    if (anonymousIdentity != m_anonymousIdentity) {
        m_anonymousIdentity = anonymousIdentity;
        emit anonymousIdentityChanged();
    }
}

QString ConnManServiceConfig::anonymousIdentity() const
{
    return m_anonymousIdentity;
}

void ConnManServiceConfig::setSubjectMatch(const QString &subjectMatch)
{
    if (subjectMatch != m_subjectMatch) {
        m_subjectMatch = subjectMatch;
        emit subjectMatchChanged();
    }
}

QString ConnManServiceConfig::subjectMatch() const
{
    return m_subjectMatch;
}

void ConnManServiceConfig::setAltSubjectMatch(const QString &altSubjectMatch)
{
    if (altSubjectMatch != m_altSubjectMatch) {
        m_altSubjectMatch = altSubjectMatch;
        emit altSubjectMatchChanged();
    }
}

QString ConnManServiceConfig::altSubjectMatch() const
{
    return m_altSubjectMatch;
}

void ConnManServiceConfig::setDomainSubjectMatch(const QString &domainSubjectMatch)
{
    if (domainSubjectMatch != m_domainSubjectMatch) {
        m_domainSubjectMatch = domainSubjectMatch;
        emit domainSubjectMatchChanged();
    }
}

QString ConnManServiceConfig::domainSubjectMatch() const
{
    return m_domainSubjectMatch;
}

void ConnManServiceConfig::setDomainMatch(const QString &domainMatch)
{
    if (domainMatch != m_domainMatch) {
        m_domainMatch = domainMatch;
        emit domainMatchChanged();
    }
}

QString ConnManServiceConfig::domainMatch() const
{
    return m_domainMatch;
}

void ConnManServiceConfig::setPhase2Type(ConnManServiceConfig::Phase2Type phase2Type)
{
    if (phase2Type != m_phase2Type) {
        m_phase2Type = phase2Type;
        emit phase2TypeChanged();
    }
}

ConnManServiceConfig::Phase2Type ConnManServiceConfig::phase2Type() const
{
    return m_phase2Type;
}

void ConnManServiceConfig::setIsPhase2EapBased(bool isPhase2EapBased)
{
    if (isPhase2EapBased != m_isPhase2EapBased) {
        m_isPhase2EapBased = isPhase2EapBased;
        emit isPhase2EapBasedChanged();
    }
}

bool ConnManServiceConfig::isPhase2EapBased() const
{
    return m_isPhase2EapBased;
}
