#pragma once

#ifdef Q_OS_UNIX
#include "unixsignalhandler.h"
#endif

#include <QGuiApplication>

template<typename T>
T tearDownAppOnSignal(const std::function<T()> &run)
{
#ifdef Q_OS_UNIX
    UnixSignalHandler::setup();
    auto ush = UnixSignalHandler();
    ush.setDoOnTermination(QGuiApplication::quit);
#endif

    return run();
}