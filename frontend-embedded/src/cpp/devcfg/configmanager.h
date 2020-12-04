#pragma once

#include <QObject>

#include "connmanserviceconfig.h"
#include "deviceconfig.h"

class SetupConfig: public QObject
{
    Q_OBJECT
    Q_PROPERTY(QString token READ token WRITE setToken NOTIFY tokenChanged)
    Q_PROPERTY(QList<ConnManServiceConfig *> ethernetServices READ networkingServices WRITE setNetworkingServices NOTIFY networkingServicesChanged)

public:
    explicit SetupConfig(QObject *parent = nullptr);

    void read(const QJsonDocument &doc);
    void read(const QJsonObject &root);

    void setNetworkingServices(const QList<ConnManServiceConfig *> &networkingServices);
    QList<ConnManServiceConfig *> networkingServices();
    void setToken(const QString &token);
    [[nodiscard]] QString token() const;

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
    Q_PROPERTY(DeviceConfig *deviceConfig READ deviceConfig)

public:
    explicit ConfigManager(QObject *parent = nullptr);

    [[nodiscard]] DeviceConfig *deviceConfig() const;

public slots:
    void loadConfig();
    void saveDeviceConfig();

signals:
    void configLoaded();

private:
    void reloadSystem();
    void loadDeviceConfig();

    DeviceConfig *m_deviceConfig;
};

Q_DECLARE_METATYPE(SetupConfig *)
