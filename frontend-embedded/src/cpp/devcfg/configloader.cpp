#include "configloader.h"
#include "usbmount.h"

#ifdef TIMETERMOS
#include "ttsystemd.h"
#endif

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
                svc->read(svcItem.toObject(), &err);
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

QString configLocation()
{
#ifdef TIMETERMOS
    return "/mnt/config/timeterm-config.json";
#else
    return "timeterm-config.json";
#endif
}

void ConfigLoader::reloadSystem()
{
#ifdef TIMETERMOS
    auto manager = org::freedesktop::systemd1::Manager("org.freedesktop.systemd1", "/org/freedesktop/systemd1", QDBusConnection::systemBus(), this);
    auto reply = manager.RestartUnit("connman.service", "replace");
    reply.waitForFinished();

    qDebug() << "Restarting ConnMan...";
    if (reply.isError())
        qCritical() << "Could not restart ConnMan:" << reply.error().message();
    else
        qDebug() << "ConnMan restarted";

    reply = manager.RestartUnit("wpa_supplicant.service", "replace");
    reply.waitForFinished();

    qDebug() << "Restarting wpa_supplicant...";
    if (reply.isError())
        qCritical() << "Could not restart wpa_supplicant:" << reply.error().message();
    else
        qDebug() << "wpa_supplicant restarted";
#endif
}

void ConfigLoader::loadConfig()
{
    qDebug() << "Loading configuration";
    auto _loadedGuard = onScopeExit([this]() {
        emit configLoaded();
    });

    qDebug() << "Trying to mount config volume...";
    if (tryMountConfig() == std::nullopt) {
        qDebug() << "Config volume mounted";
        auto _unmountGuard = onScopeExit([]() {
            if (tryUnmountConfig() == std::nullopt) {
                qDebug() << "Config volume unmounted";
            } else {
                qCritical() << "Unmounting config volume failed";
            }
        });

        qDebug() << "Trying to load config...";
        auto doc = QFile(configLocation());
        if (!doc.exists()) {
            qDebug() << "Config file not present";
            return;
        }
        qDebug() << "Config file present, opening...";

        if (!doc.open(QIODevice::ReadOnly)) {
            qDebug() << "Could not open config file";
            return;
        }
        auto bytes = doc.readAll();
        qDebug() << "Config file opened";

        qDebug() << "Parsing config file...";
        auto parseError = QJsonParseError();
        auto jsonDoc = QJsonDocument::fromJson(bytes, &parseError);
        if (parseError.error) {
            qDebug() << "Invalid config file";
            return;
        }
        qDebug() << "Parsed config file";

        qDebug() << "Reading config file...";
        auto config = new Config(this);
        config->read(jsonDoc);
        qDebug() << "Config file read";

        for (auto &svc : config->ethernetServices()) {
            qDebug() << "Configuring ethernet service" << svc->serviceName();
            svc->saveCerts();
            svc->saveConnManConf();
            qDebug() << "Ethernet service" << svc->serviceName() << "configured";
        }
    } else {
        qDebug() << "Mounting config volume failed";
    }

    qDebug() << "Reloading system...";
    reloadSystem();
}
