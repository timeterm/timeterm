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
    json["deviceTokenOrganizationId"] = m_deviceTokenOrganizationId;
    json["setupTokenOrganizationId"] = m_setupTokenOrganizationId;
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
    if (json.contains("deviceTokenOrganizationId") && json["deviceTokenOrganizationId"].isString())
        setDeviceTokenOrganizationId(json["deviceTokenOrganizationId"].toString());
    if (json.contains("setupTokenOrganizationId") && json["setupTokenOrganizationId"].isString())
        setSetupTokenOrganizationId(json["setupTokenOrganizationId"].toString());
}

void DeviceConfig::setDeviceTokenOrganizationId(const QString &id)
{
    if (id != m_deviceTokenOrganizationId) {
        m_deviceTokenOrganizationId = id;
        emit deviceTokenOrganizationIdChanged();
    }
}

QString DeviceConfig::deviceTokenOrganizationId() const
{
    return m_deviceTokenOrganizationId;
}

bool DeviceConfig::needsRegistration()
{
    return m_setupToken != "" && m_setupTokenOrganizationId != m_deviceTokenOrganizationId;
}

void DeviceConfig::setSetupTokenOrganizationId(const QString &id)
{
    if (id != m_setupTokenOrganizationId) {
        m_setupTokenOrganizationId = id;
        emit setupTokenOrganizationIdChanged();
    }
}

QString DeviceConfig::setupTokenOrganizationId() const
{
    return m_deviceTokenOrganizationId;
}
