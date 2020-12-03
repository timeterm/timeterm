#pragma once

#include <QObject>

#include "connmanserviceconfig.h"
#include "deviceinfo.h"

class Config: public QObject
{
    Q_OBJECT
    Q_PROPERTY(QString token READ token WRITE setToken NOTIFY tokenChanged)
    Q_PROPERTY(QList<ConnManServiceConfig *> ethernetServices READ networkingServices WRITE setNetworkingServices NOTIFY networkingServicesChanged)

public:
    explicit Config(QObject *parent = nullptr);

    void read(const QJsonDocument &doc);
    void read(const QJsonObject &root);

    void setNetworkingServices(const QList<ConnManServiceConfig *> &networkingServices);
    QList<ConnManServiceConfig *> networkingServices();
    void setToken(const QString &token);
    [[nodiscard]] QString token() const;

    void saveSignupToken();

signals:
    void networkingServicesChanged();
    void tokenChanged();

private:
    QList<ConnManServiceConfig *> m_networkingServices;
    QString m_token;
};

class ConfigManager: public QObject
{
    Q_OBJECT
    Q_PROPERTY(DeviceInfo *deviceInfo READ deviceInfo)

public:
    explicit ConfigManager(QObject *parent = nullptr);

    [[nodiscard]] DeviceInfo *deviceInfo() const;

public slots:
    void loadConfig();

signals:
    void configLoaded();

private:
    void reloadSystem();
    void loadDeviceConfig();

    DeviceInfo *m_deviceInfo;
};

Q_DECLARE_METATYPE(Config *)
