#include "device.h"

void Device::read(const QJsonObject &json)
{
    if (json.contains("id") && json["id"].isString())
        id = json["id"].toString();
    if (json.contains("organizationId") && json["organizationId"].isString())
        organizationId = json["organizationId"].toString();
    if (json.contains("name") && json["name"].isString())
        name = json["name"].toString();
}

bool operator==(const Device &a, const Device &b)
{
    return a.id == b.id
        && a.name == b.name
        && a.organizationId == b.organizationId;
}

bool operator!=(const Device &a, const Device &b)
{
    return !(a == b);
}
