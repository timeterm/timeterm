#include "deviceinfo.h"

DeviceInfo::DeviceInfo(QObject *parent)
    : QObject(parent)
{
}

void DeviceInfo::setId(const QString &id)
{
    if (id != m_id) {
        m_id = id;
        emit idChanged();
    }
}

QString DeviceInfo::id() const
{
    return m_id;
}

void DeviceInfo::setName(const QString &name)
{
    if (name != m_name) {
        m_name = name;
        emit nameChanged();
    }
}

QString DeviceInfo::name() const
{
    return m_name;
}

void DeviceInfo::setToken(const QString &token)
{
    if (token != m_token) {
        m_token = token;
        emit tokenChanged();
    }
}

QString DeviceInfo::token() const
{
    return m_token;
}

void DeviceInfo::write(QJsonObject &json) const
{
    json["id"] = m_id;
    json["name"] = m_name;
    json["token"] = m_token;
}

void DeviceInfo::read(const QJsonObject &json)
{
    if (json.contains("id") && json["id"].isString())
        m_id = json["id"].toString();
    if (json.contains("name") && json["nam"].isString())
        m_name = json["name"].toString();
    if (json.contains("token") && json["token"].isString())
        m_token = json["token"].toString();
}
