#pragma once

#include "qsystemdetection.h"

#ifdef Q_OS_UNIX

#include <QObject>
#include <QSocketNotifier>

class UnixSignalHandler: public QObject {
    Q_OBJECT

public:
    explicit UnixSignalHandler(QObject *parent = nullptr);

    void setDoOnTermination(const std::function<void()> &fn);

    static void setup();
    static void intSignalHandler(int _);
    static void termSignalHandler(int _);

public slots:
    void handleSigInt();
    void handleSigTerm();

private:
    static int sigintFd[2];
    static int sigtermFd[2];

    QSocketNotifier *m_snInt;
    QSocketNotifier *m_snTerm;

    std::function<void()> m_doOnTermination;
};

#endif
