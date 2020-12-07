#pragma once

#include <QJsonObject>
#include <QObject>

class NatsCredsResponse
{
    Q_GADGET
    Q_PROPERTY(QString credentials MEMBER credentials)

public slots:
    Q_INVOKABLE void writeToFile() const;

public:
    void read(const QJsonObject &json);

    QString credentials;
};

Q_DECLARE_METATYPE(NatsCredsResponse)
