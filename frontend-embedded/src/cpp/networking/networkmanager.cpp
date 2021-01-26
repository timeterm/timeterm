#include "networkmanager.h"

#include <QDebug>
#include <QTimerEvent>

#ifdef TIMETERMOS
#include <QNetworkSettingsInterface>
#include <QNetworkSettingsInterfaceModel>
#include <QNetworkSettingsManager>
#include <QNetworkSettingsService>
#include <QNetworkSettingsServiceModel>
#include <QNetworkSettingsType>
#endif

NetworkManager::NetworkManager(QObject *parent)
    : QObject(parent)
    , m_workerThread(new QThread(this))
{
    m_worker = new NetworkManagerWorker();
    m_worker->moveToThread(m_workerThread);
    connect(m_workerThread, &QThread::started, m_worker, &NetworkManagerWorker::start);
    connect(m_workerThread, &QThread::finished, m_worker, &NetworkManagerWorker::deleteLater);
    connect(m_worker, &NetworkManagerWorker::networkStateRetrieved, this, &NetworkManager::networkStateRetrieved);
    connect(this, &NetworkManager::retrieveNewNetworkState, m_worker, &NetworkManagerWorker::retrieveNewNetworkState);
    connect(this, &NetworkManager::configLoaded, m_worker, &NetworkManagerWorker::configLoaded);
    connect(this, &NetworkManager::activateInactiveNetworkingInterfaces, m_worker, &NetworkManagerWorker::activateInactiveNetworkingInterfaces);
    m_workerThread->start();
}

NetworkManager::~NetworkManager()
{
    m_workerThread->quit();
    m_workerThread->wait();
}

void NetworkManager::networkStateRetrieved(NetworkState state)
{
    if (!m_lastState.has_value() || *m_lastState != state) {
        if (!m_lastState.has_value() || (*m_lastState).isOnline != state.isOnline) {
            emit onlineChanged(state.isOnline);
        }

        m_lastState = state;
        emit stateChanged(state);
    }
}

bool operator==(const NetworkState &a, const NetworkState &b)
{
    return a.isConnected == b.isConnected
        && a.isOnline == b.isOnline
        && a.isWired == b.isWired
        && a.signalStrength == b.signalStrength;
}

bool operator!=(const NetworkState &a, const NetworkState &b)
{
    return !(a == b);
}

NetworkManagerWorker::NetworkManagerWorker(QObject *parent)
#ifdef TIMETERMOS
    : m_manager(new QNetworkSettingsManager(this))
#endif
{
#ifdef TIMETERMOS
    QObject::connect(m_manager, &QNetworkSettingsManager::interfacesChanged, this, &NetworkManagerWorker::networkingInterfacesChanged);
    QObject::connect(m_manager, &QNetworkSettingsManager::servicesChanged, this, &NetworkManagerWorker::servicesChanged);
#endif
}

void NetworkManagerWorker::start()
{
    qDebug() << "NetworkManagerWorker: starting timers...";
    m_checkNetworkStateTimerId = startTimer(5000);
}

void NetworkManagerWorker::timerEvent(QTimerEvent *event)
{
    if (event->timerId() == m_checkNetworkStateTimerId)
        retrieveNewNetworkState();
}

void NetworkManagerWorker::retrieveNewNetworkState()
{
    auto state = NetworkState();
#ifdef TIMETERMOS
    if (!m_manager) return;

    auto *svc = m_manager->currentWifiConnection();
    if (svc != nullptr) {
        state.isOnline = svc->state() == QNetworkSettingsState::Online;
        state.isConnected = state.isOnline || (svc->state() == QNetworkSettingsState::Ready);
        state.signalStrength = svc->wirelessConfig()->signalStrength();

        if (svc->ipv4() != nullptr && svc->ipv4()->address() != "")
            state.ip = svc->ipv4()->address();
        else if (svc->ipv6() != nullptr && svc->ipv6()->address() != "")
            state.ip = svc->ipv6()->address();
    }

    svc = m_manager->currentWiredConnection();
    if (svc != nullptr && !state.isOnline) {
        state.isOnline = svc->state() == QNetworkSettingsState::Online;
        if (!state.isConnected) {
            state.isConnected = state.isOnline || (svc->state() == QNetworkSettingsState::Ready);
        }

        if (svc->ipv4() != nullptr && svc->ipv4()->address() != "")
            state.ip = svc->ipv4()->address();
        else if (svc->ipv6() != nullptr && svc->ipv6()->address() != "")
            state.ip = svc->ipv6()->address();
    }
#else
    // Sleep here for a while so we can recognize when threading has been borked
    // (in that case the UI would freeze for 3 seconds about every 5 seconds).
    QThread::sleep(3);
    if (m_configLoaded) {
        state.isOnline = true;
        state.isWired = false;
        state.signalStrength = 50;
        state.isConnected = true;
    }
#endif
    emit networkStateRetrieved(state);
}

void NetworkManagerWorker::activateInactiveNetworkingInterfaces()
{
#ifdef TIMETERMOS
    auto interfaces = m_manager->interfaces()->getModel();

    int i = 0;
    for (auto &iface : interfaces) {
        i++;

        if (iface->type() != QNetworkSettingsType::Wifi) continue;

        if (!iface->powered()) iface->setPowered(true);
        else iface->scanServices();
    }
#endif
}

void NetworkManagerWorker::servicesChanged()
{
#ifdef TIMETERMOS
    QList<QNetworkSettingsService *> services = qobject_cast<QNetworkSettingsServiceModel *>(m_manager->services()->sourceModel())->getModel();

    int i = 0;
    for (const auto &service : services) {
        i++;

        if (service->type() != QNetworkSettingsType::Wifi) continue;

        QString stateString = "";
        switch (service->state()) {
        case QNetworkSettingsState::Idle:
            stateString = "Idle";
            break;
        case QNetworkSettingsState::Failure:
            stateString = "Failure";
            break;
        case QNetworkSettingsState::Association:
            stateString = "Association";
            break;
        case QNetworkSettingsState::Configuration:
            stateString = "Configuration";
            break;
        case QNetworkSettingsState::Ready:
            stateString = "Ready";
            break;
        case QNetworkSettingsState::Disconnect:
            stateString = "Disconnect";
            break;
        case QNetworkSettingsState::Online:
            stateString = "Online";
            break;
        case QNetworkSettingsState::Undefined:
            stateString = "Undefined";
            break;
        }
        qDebug() << "TtNetworkManager: service" << i << "," << service->name() << "currently has state" << stateString;
    }
#endif
}

void NetworkManagerWorker::networkingInterfacesChanged()
{
    if (m_configLoaded)
        emit activateInactiveNetworkingInterfaces();
}

void NetworkManagerWorker::configLoaded()
{
    m_configLoaded = true;
    activateInactiveNetworkingInterfaces();
    retrieveNewNetworkState();
}
