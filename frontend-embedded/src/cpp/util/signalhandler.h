#pragma once

#include <QGuiApplication>

#ifdef Q_OS_UNIX
#include "unixsignalhandler.h"
#endif

template<typename T>
T teardownAppOnSignal(const std::function<T()> &run)
{
#ifdef Q_OS_UNIX
    UnixSignalHandler::setup();
    auto ush = UnixSignalHandler();
    ush.setDoOnTermination(QGuiApplication::quit);
#endif

    return run();
}