#pragma once

#include <QObject>
#include <devcfg/connmanserviceconfig.h>

class NetworkingServicesResponse
{
    Q_GADGET

public:
    void append(const QSharedPointer<ConnManServiceConfig> &service);
    void append(const QList<QSharedPointer<ConnManServiceConfig>> &services);

    void read(const QJsonArray &json);
    Q_INVOKABLE void save();

private:
    QList<QSharedPointer<ConnManServiceConfig>> m_services;
};

Q_DECLARE_METATYPE(NetworkingServicesResponse)
