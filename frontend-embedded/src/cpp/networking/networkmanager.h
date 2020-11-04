#pragma once

#include <QObject>

#ifdef TIMETERMOS
#include <QNetworkSettingsManager>
#endif

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

    bool m_configLoaded = false;

#ifdef TIMETERMOS
    QNetworkSettingsManager *m_manager;
#endif
};
