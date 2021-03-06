#pragma once

#include <QObject>
#include <QThread>

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
    Q_PROPERTY(QString ip MEMBER ip)

public:
    bool isConnected = false;
    bool isOnline = false;
    bool isWired = false;
    int signalStrength = 0;
    QString ip;
};

bool operator==(const NetworkState &a, const NetworkState &b);
bool operator!=(const NetworkState &a, const NetworkState &b);

class NetworkManagerWorker: public QObject
{
    Q_OBJECT

public:
    explicit NetworkManagerWorker(QObject *parent = nullptr);

signals:
    void networkStateRetrieved(NetworkState);

public slots:
    void start();
    void configLoaded();
    void retrieveNewNetworkState();
    void activateInactiveNetworkingInterfaces();
    void networkingInterfacesChanged();
    void servicesChanged();

protected:
    void timerEvent(QTimerEvent *event) override;

private:
    bool m_configLoaded = false;
    int m_checkNetworkStateTimerId = 0;

#ifdef TIMETERMOS
    QNetworkSettingsManager *m_manager;
#endif
};

class NetworkManager: public QObject
{
    Q_OBJECT

public:
    explicit NetworkManager(QObject *parent = nullptr);
    ~NetworkManager() override;

signals:
    void stateChanged(NetworkState);
    void onlineChanged(bool online);
    void retrieveNewNetworkState();
    void activateInactiveNetworkingInterfaces();
    void configLoaded();

private slots:
    void networkStateRetrieved(NetworkState state);

private:
    std::optional<NetworkState> m_lastState = std::nullopt;

    QThread *m_workerThread;
    NetworkManagerWorker *m_worker = nullptr;
};

Q_DECLARE_METATYPE(NetworkState)
