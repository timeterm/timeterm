#include "configmanager.h"
#include "usbmount.h"

#ifdef TIMETERMOS
#include "ttsystemd.h"
#endif

#include <QJsonDocument>

#include <QDir>
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
    if (obj.contains("networkingServices") && obj["networkingServices"].isArray()) {
        auto arr = obj["networkingServices"].toArray();
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

        setNetworkingServices(services);
    }
}

void Config::setNetworkingServices(const QList<ConnManServiceConfig *> &networkingServices)
{
    if (networkingServices != m_networkingServices) {
        m_networkingServices = networkingServices;
        emit networkingServicesChanged();
    }
}

QList<ConnManServiceConfig *> Config::networkingServices()
{
    return m_networkingServices;
}

void Config::setToken(const QString &token)
{
    if (token != m_token) {
        m_token = token;
        emit tokenChanged();
    }
}

QString Config::token() const
{
    return m_token;
}

QString createSignupTokenPath() {
    const QString filename = QStringLiteral("signup-token");
    auto relative = QStringLiteral("tokens/");

#if TIMETERMOS
    return "/opt/frontend-embedded/" + relative;
#else
    const QString &dir = relative;
#endif

    QDir(dir).mkpath(dir);

    return dir + filename;
}

QString createDeviceTokenPath() {
    const QString filename = QStringLiteral("device-token");
    auto relative = QStringLiteral("tokens/");

#if TIMETERMOS
    return "/opt/frontend-embedded/" + relative;
#else
    const QString &dir = relative;
#endif

    QDir(dir).mkpath(dir);

    return dir + filename;
}

QString createDeviceIdPath() {
    QString filename = QStringLiteral("device-id");

#if TIMETERMOS
    return "/opt/frontend-embedded/" + filename;
#endif
    return filename;
}

void Config::saveSignupToken()
{
    auto path = createSignupTokenPath();
    auto f = QFile(path);
    if (!f.open(QIODevice::WriteOnly | QIODevice::Truncate)) {
        qCritical() << "Could not open signup token file";
        return;
    }

    auto tokenBytes = m_token.toLocal8Bit();
    f.write(tokenBytes);
    f.close();
}

ConfigManager::ConfigManager(QObject *parent)
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

void ConfigManager::reloadSystem()
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

void ConfigManager::loadConfig()
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

        config->saveSignupToken();

        for (auto &svc : config->networkingServices()) {
            qDebug() << "Configuring ethernet service" << svc->name();
            svc->saveCerts();
            svc->saveConnManConf();
            qDebug() << "Ethernet service" << svc->name() << "configured";
        }
    } else {
        qDebug() << "Mounting config volume failed";
    }

    qDebug() << "Reloading system...";
    reloadSystem();
}

void ConfigManager::setDeviceInfo(const QString &token)
{

}
