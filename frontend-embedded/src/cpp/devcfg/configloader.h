#pragma once

#include <QObject>

#include "connmanserviceconfig.h"

class Config: public QObject
{
    Q_OBJECT
    Q_PROPERTY(QList<ConnManServiceConfig *> ethernetServices READ ethernetServices WRITE setEthernetServices NOTIFY ethernetServicesChanged)

public:
    explicit Config(QObject *parent = nullptr);

    void read(const QJsonDocument &doc);
    void read(const QJsonObject &root);

    void setEthernetServices(const QList<ConnManServiceConfig *> &ethernetServices);
    QList<ConnManServiceConfig *> ethernetServices();

signals:
    void ethernetServicesChanged();

private:
    QList<ConnManServiceConfig *> m_ethernetServices;
};

class ConfigLoader: public QObject
{
    Q_OBJECT

public:
    explicit ConfigLoader(QObject *parent = nullptr);

public slots:
    void loadConfig();

signals:
    void configLoaded();

private:
    void restartConnMan();
};

Q_DECLARE_METATYPE(Config *)
