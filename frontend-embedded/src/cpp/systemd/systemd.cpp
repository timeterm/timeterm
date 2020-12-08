#include "systemd.h"

Systemd::Systemd(QObject *parent)
    : QObject(parent)
{
}

void Systemd::rebootDevice() {
#ifdef TIMETERMOS
    auto manager = org::freedesktop::systemd1::Manager("org.freedesktop.systemd1", "/org/freedesktop/systemd1", QDBusConnection::systemBus());
    auto reply = manager.Reboot();

    reply.waitForFinished();

    qDebug() << "Rebooting...";
    if (reply.isError())
        qCritical() << "Could not reboot:" << reply.error().message();
    else
        qDebug() << "Successfully issued reboot command to systemd";
#endif
}
