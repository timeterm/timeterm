#include "configloader.h"
#include "usbmount.h"

#include <QJsonDocument>
#include <util/scopeguard.h>

Config::Config(QObject *parent)
    : QObject(parent)
{
}

void Config::read(const QJsonDocument &doc)
{
    if (!doc.isObject())
        return; // TODO: return error

    read(doc.object());
}

void Config::read(const QJsonObject &obj)
{
    if (obj.contains("ethernetServices") && obj["ethernetServices"].isArray()) {
        auto arr = obj["ethernetServices"].toArray();
        auto services = QList<ConnManServiceConfig *>();
        services.reserve(arr.size());

        for (auto svcItem : arr) {
            if (svcItem.isObject()) {
                auto svc = new ConnManServiceConfig(this);
                auto err = ConnManServiceConfig::ReadErrorNoError;
                svc->read(obj, &err);
                if (err == ConnManServiceConfig::ReadErrorNoError) {
                    services.append(svc);
                }
            }
        }

        setEthernetServices(services);
    }
}

void Config::setEthernetServices(const QList<ConnManServiceConfig *> &ethernetServices)
{
    if (ethernetServices != m_ethernetServices) {
        m_ethernetServices = ethernetServices;
        emit ethernetServicesChanged();
    }
}

QList<ConnManServiceConfig *> Config::ethernetServices()
{
    return m_ethernetServices;
}

ConfigLoader::ConfigLoader(QObject *parent)
    : QObject(parent)
{
}

QString configLocation() {
#ifdef TIMETERMOS
    return "/mnt/config/timeterm-config.json";
#else
    return "timeterm-config.json";
#endif
}

void ConfigLoader::loadConfig()
{
    auto _loadedGuard = onScopeExit([this]() {
      emit configLoaded();
    });

    if (tryMountConfig() == std::nullopt) {
        auto _unmountGuard = onScopeExit([]() {
            tryUnmountConfig();
        });

        auto doc = QFile(configLocation());
        if (!doc.exists())
            return;

        if (!doc.open(QIODevice::ReadOnly))
            return;
        auto bytes = doc.readAll();

        auto parseError = QJsonParseError();
        auto jsonDoc = QJsonDocument::fromJson(bytes, &parseError);
        if (parseError.error)
            return;

        auto config = new Config(this);
        config->read(jsonDoc);

        for (auto &svc : config->ethernetServices()) {
            svc->saveCerts();
            svc->saveConnManConf();
        }
    }
}
