#pragma once

#include <QObject>
#include <QJsonObject>

class NatsCredsResponse
{
    Q_GADGET
    Q_PROPERTY(QString credentials MEMBER credentials)

public:
    void read(const QJsonObject &json);

    QString credentials;
};

Q_DECLARE_METATYPE(NatsCredsResponse)
