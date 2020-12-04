#include "deviceconfig.h"

#include <QCryptographicHash>

DeviceConfig::DeviceConfig(QObject *parent)
    : QObject(parent)
{
}

void DeviceConfig::setId(const QString &id)
{
    if (id != m_id) {
        m_id = id;
        emit idChanged();
    }
}

QString DeviceConfig::id() const
{
    return m_id;
}

void DeviceConfig::setName(const QString &name)
{
    if (name != m_name) {
        m_name = name;
        emit nameChanged();
    }
}

QString DeviceConfig::name() const
{
    return m_name;
}

void DeviceConfig::setSetupToken(const QString &token)
{
    if (token != m_setupToken) {
        m_setupToken = token;
        emit setupTokenChanged();
    }
}

QString DeviceConfig::setupToken() const
{
    return m_setupToken;
}

void DeviceConfig::setDeviceToken(const QString &token)
{
    if (token != m_deviceToken) {
        m_deviceToken = token;
        emit deviceTokenChanged();
    }
}

QString DeviceConfig::deviceToken() const
{
    return m_deviceToken;
}

void DeviceConfig::write(QJsonObject &json) const
{
    json["id"] = m_id;
    json["name"] = m_name;
    json["setupToken"] = m_setupToken;
    json["deviceToken"] = m_deviceToken;
    json["deviceTokenSetupTokenHash"] = m_deviceTokenSetupTokenHash;
}

void DeviceConfig::read(const QJsonObject &json)
{
    if (json.contains("id") && json["id"].isString())
        setId(json["id"].toString());
    if (json.contains("name") && json["name"].isString())
        setName(json["name"].toString());
    if (json.contains("setupToken") && json["setupToken"].isString())
        setSetupToken(json["setupToken"].toString());
    if (json.contains("deviceToken") && json["deviceToken"].isString())
        setDeviceToken(json["deviceToken"].toString());
    if (json.contains("deviceTokenSetupTokenHash") && json["deviceTokenSetupTokenHash"].isString())
        setDeviceTokenSetupTokenHash(json["deviceTokenSetupTokenHash"].toString());
}

void DeviceConfig::setDeviceTokenSetupToken(const QString &token)
{
    setDeviceTokenSetupTokenHash(hashToken(token));
}

void DeviceConfig::setDeviceTokenSetupTokenHash(const QString &hash)
{
    if (hash != m_deviceTokenSetupTokenHash) {
        m_deviceTokenSetupTokenHash = hash;
        emit deviceTokenSetupTokenHashChanged();
    }
}

QString DeviceConfig::deviceTokenSetupTokenHash() const
{
    return m_deviceTokenSetupTokenHash;
}

QString DeviceConfig::hashToken(const QString &token)
{
    auto bytes = token.toUtf8();
    auto hash = QCryptographicHash(QCryptographicHash::Sha3_256);
    hash.addData(bytes);

    auto hashBytes = hash.result();
    auto hashString = QString(hashBytes.toHex());
    return hashString;
}

bool DeviceConfig::needsRegistration()
{
    return m_setupToken != "" && hashToken(m_setupToken) != m_deviceTokenSetupTokenHash;
}
