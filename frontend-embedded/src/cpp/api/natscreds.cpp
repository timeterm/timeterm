#include "natscreds.h"

void NatsCredsResponse::read(const QJsonObject &json)
{
    if (json.contains("credentials") && json["credentials"].isString())
        credentials = json["credentials"].toString();
}
