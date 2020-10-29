#pragma once

#include <QObject>
#include <QNetworkSettingsManager>

class NetworkManager: public QObject
{
    Q_OBJECT

public:
    explicit NetworkManager(QObject *parent = nullptr);

private slots:
    void networkingInterfacesChanged();
    void servicesChanged();

public slots:
    void configLoaded();

private:
    void activateInactiveNetworkingInterfaces();

    QNetworkSettingsManager *m_manager;
    bool m_configLoaded = false;
};
