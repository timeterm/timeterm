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
#ifdef TIMETERMOS
    , m_manager(new QNetworkSettingsManager(this))
#endif
{
#ifdef TIMETERMOS
    QObject::connect(m_manager, &QNetworkSettingsManager::interfacesChanged, this, &NetworkManager::networkingInterfacesChanged);
    QObject::connect(m_manager, &QNetworkSettingsManager::servicesChanged, this, &NetworkManager::servicesChanged);
#endif

    m_checkNetworkStateTimerId = startTimer(5000);
}

void NetworkManager::configLoaded()
{
    m_configLoaded = true;
    activateInactiveNetworkingInterfaces();
}

void NetworkManager::networkingInterfacesChanged()
{
    if (m_configLoaded)
        activateInactiveNetworkingInterfaces();
}

void NetworkManager::activateInactiveNetworkingInterfaces()
{
#ifdef TIMETERMOS
    auto interfaces = m_manager->interfaces()->getModel();
    qDebug() << "TtNetworkManager: found" << interfaces.size() << "interfaces";

    int i = 0;
    for (auto &iface : interfaces) {
        i++;

        if (iface->type() == QNetworkSettingsType::Wifi)
            qDebug() << "TtNetworkManager: interface" << i << "," << iface->name() << "is a wireless interface";
        else {
            qDebug() << "TtNetworkManager: interface" << i << "," << iface->name() << "is not a wireless interface";
            continue;
        }

        if (!iface->powered()) {
            qDebug() << "TtNetworkManager: interface" << i << "," << iface->name() << "is not yet powered, powering it on";
            iface->setPowered(true);
        } else {
            qDebug() << "TtNetworkManager: interface" << i << "," << iface->name() << "is already powered, scanning";
            iface->scanServices();
        }
    }
#endif
}

void NetworkManager::servicesChanged()
{
#ifdef TIMETERMOS
    QList<QNetworkSettingsService *> services = qobject_cast<QNetworkSettingsServiceModel *>(m_manager->services()->sourceModel())->getModel();
    qDebug() << "TtNetworkManager: found" << services.size() << "services";

    int i = 0;
    for (const auto &service : services) {
        i++;

        if (service->type() == QNetworkSettingsType::Wifi)
            qDebug() << "TtNetworkManager: service" << i << "," << service->name() << "is a wireless network";
        else {
            qDebug() << "TtNetworkManager: service" << i << "," << service->name() << "is not a wireless network";
            continue;
        }

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

NetworkState NetworkManager::getNetworkState()
{
    auto state = NetworkState();
#ifdef TIMETERMOS
    auto *svc = m_manager->currentWifiConnection();
    if (svc != nullptr) {
        state.isOnline = svc->state() == QNetworkSettingsState::Online;
        state.isConnected = state.isOnline || (svc->state() == QNetworkSettingsState::Ready);
        state.signalStrength = svc->wirelessConfig()->signalStrength();
    }

    svc = m_manager->currentWiredConnection();
    if (svc != nullptr && !state.isOnline) {
        state.isOnline = svc->state() == QNetworkSettingsState::Online;
        if (!state.isConnected) {
            state.isConnected = state.isOnline || (svc->state() == QNetworkSettingsState::Ready);
        }
    }
#else
    state.isOnline = true;
    state.isWired = false;
    state.signalStrength = 50;
    state.isConnected = true;
#endif
    return state;
}

void NetworkManager::timerEvent(QTimerEvent *event)
{
    if (event->timerId() == m_checkNetworkStateTimerId) {
        auto state = getNetworkState();
        if (!m_lastState.has_value() || *m_lastState != state) {
            if (!m_lastState.has_value() || (*m_lastState).isOnline != state.isOnline) {
                emit onlineChanged(state.isOnline);
            }

            m_lastState = state;
            emit stateChanged(state);
        }
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
