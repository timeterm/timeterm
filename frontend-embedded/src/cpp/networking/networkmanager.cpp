#include "networkmanager.h"

#include <QDebug>
#include <QNetworkSettingsInterface>
#include <QNetworkSettingsInterfaceModel>
#include <QNetworkSettingsManager>
#include <QNetworkSettingsService>
#include <QNetworkSettingsServiceModel>
#include <QNetworkSettingsType>

NetworkManager::NetworkManager(QObject *parent)
    : QObject(parent)
    , m_manager(new QNetworkSettingsManager(this))
{
    QObject::connect(m_manager, &QNetworkSettingsManager::interfacesChanged, this, &NetworkManager::networkingInterfacesChanged);
    QObject::connect(m_manager, &QNetworkSettingsManager::servicesChanged, this, &NetworkManager::servicesChanged);
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
    auto interfaces = m_manager->interfaces()->getModel();
    qDebug() << "TtNetworkManager: found" << interfaces.size() << "interfaces";

    int i = 0;
    for (auto &iface : interfaces) {
        if (iface->type() == QNetworkSettingsType::Wifi)
            qDebug() << "TtNetworkManager: interface" << i << "," << iface->name() << "is a wireless network";
        else {
            qDebug() << "TtNetworkManager: interface" << i << "," << iface->name() << "is not a wireless network";
            continue;
        }

        if (!iface->powered()) {
            qDebug() << "TtNetworkManager: interface" << i << "," << iface->name() << "is not yet powered, powering it on";
            iface->setPowered(true);
        } else {
            qDebug() << "TtNetworkManager: interface" << i << "," << iface->name() << "is already powered, scanning";
            iface->scanServices();
        }

        i++;
    }
}

void NetworkManager::servicesChanged()
{
    QList<QNetworkSettingsService *> services = qobject_cast<QNetworkSettingsServiceModel *>(m_manager->services()->sourceModel())->getModel();
    qDebug() << "TtNetworkManager: found" << services.size() << "services";

    int i = 0;
    for (const auto &service : services) {
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

        i++;
    }
}
