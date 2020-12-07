#include "systemd.h"

Systemd::Systemd(QObject *parent)
    : QObject(parent)
{
}

void Systemd::rebootDevice() {
#ifdef TIMETERMOS
    auto manager = org::freedesktop::systemd1::Manager("org.freedesktop.systemd1", "/org/freedesktop/systemd1", QDBusConnection::systemBus());
    manager.Reboot();
#endif
}
