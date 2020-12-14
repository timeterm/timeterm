#include "servicesresponse.h"

void NetworkingServicesResponse::read(const QJsonArray &json)
{
    for (const auto &it : json) {
        if (it.isObject()) {
            auto service = new ConnManServiceConfig();
            ConnManServiceConfig::ReadError err;
            service->read(it.toObject(), &err);
            if (!err) {
                auto ptr = QSharedPointer<ConnManServiceConfig>(service);
                m_services.append(ptr);
            }
        }
    }
}

void NetworkingServicesResponse::save()
{
    ConnManServiceConfig::deleteCurrentConnManConfigs();
    for (const auto& svc : m_services) {
        svc->saveCerts();
        svc->saveConnManConf();
    }
}

void NetworkingServicesResponse::append(const QList<QSharedPointer<ConnManServiceConfig>> &services)
{
    m_services.append(services);
}

void NetworkingServicesResponse::append(const QSharedPointer<ConnManServiceConfig> &service)
{
    m_services.append(service);
}
