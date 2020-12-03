#include "createdevice.h"

void CreateDeviceRequest::write(QJsonObject &json) const
{
    json["name"] = name;
}

void CreateDeviceResponse::read(const QJsonObject &json)
{
    if (json.contains("device") && json["device"].isObject())
        device.read(json["device"].toObject());
    if (json.contains("setupToken") && json["setupToken"].isString())
        token = json["setupToken"].toString();
}
