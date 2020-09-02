#include "signalhandler.h"
#include <QGuiApplication>

#ifdef Q_OS_UNIX
#include "unixsignalhandler.h"
#endif

void teardownAppOnSignal() {
#ifdef Q_OS_UNIX
    auto ush = UnixSignalHandler();
    UnixSignalHandler::setup();
    ush.setDoOnTermination(QGuiApplication::quit);
#endif
}