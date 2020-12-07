#pragma once

#include <QObject>

#ifdef TIMETERMOS
#include <QNetworkSettingsManager>
#endif

class NetworkState
{
    Q_GADGET
    Q_PROPERTY(bool isConnected MEMBER isConnected)
    Q_PROPERTY(bool isOnline MEMBER isOnline)
    Q_PROPERTY(bool isWired MEMBER isWired)
    Q_PROPERTY(int signalStrength MEMBER signalStrength)

public:
    bool isConnected = false;
    bool isOnline = false;
    bool isWired = false;
    int signalStrength = 0;
};

bool operator==(const NetworkState &a, const NetworkState &b);
bool operator!=(const NetworkState &a, const NetworkState &b);

class NetworkManager: public QObject
{
    Q_OBJECT

public:
    explicit NetworkManager(QObject *parent = nullptr);

    Q_INVOKABLE NetworkState getNetworkState();

signals:
    void stateChanged(NetworkState);
    void onlineChanged(bool online);

private slots:
    void networkingInterfacesChanged();
    void servicesChanged();

public slots:
    void configLoaded();
    void checkNetworkState();

protected:
    void timerEvent(QTimerEvent *event) override;

private:
    void activateInactiveNetworkingInterfaces();

    bool m_configLoaded = false;
    int m_checkNetworkStateTimerId = 0;
    std::optional<NetworkState> m_lastState = std::nullopt;

#ifdef TIMETERMOS
    QNetworkSettingsManager *m_manager;
#endif
};
