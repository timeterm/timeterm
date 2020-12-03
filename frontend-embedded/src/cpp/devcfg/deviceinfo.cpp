#include "deviceinfo.h"

void DeviceInfo::write(QJsonObject &json) {
    json["id"] = m_id;
    json["name"] = m_name;
    json["token"] = m_token;
}

void DeviceInfo::read(const QJsonObject &json) {
    if (json.contains("id") && id.isString())
        m_id = json["id"].toString();
    if (json.contains("name") && name.isString())
        m_name = json["name"].toString();
    if (json.contains("token") && token.isString())
        m_token = json["token"].toString();
}
