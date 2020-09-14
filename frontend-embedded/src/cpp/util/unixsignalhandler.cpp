#include <qsystemdetection.h>

#ifdef Q_OS_UNIX

#include "unixsignalhandler.h"

#include <csignal>
#include <sys/socket.h>
#include <unistd.h>

#include <QDebug>

int UnixSignalHandler::sigintFd[2] = {0, 0};
int UnixSignalHandler::sigtermFd[2] = {0, 0};

UnixSignalHandler::UnixSignalHandler(QObject *parent)
    : QObject(parent)
{
    if (::socketpair(AF_UNIX, SOCK_STREAM, 0, sigintFd))
        qFatal("Couldn't create INT socketpair");
    if (::socketpair(AF_UNIX, SOCK_STREAM, 0, sigtermFd))
        qFatal("Couldn't create TERM socketpair");

    m_snInt = new QSocketNotifier(sigintFd[1], QSocketNotifier::Read, this);
    connect(m_snInt, &QSocketNotifier::activated, this, &UnixSignalHandler::handleSigInt);
    m_snTerm = new QSocketNotifier(sigtermFd[1], QSocketNotifier::Read, this);
    connect(m_snTerm, &QSocketNotifier::activated, this, &UnixSignalHandler::handleSigTerm);
}

void UnixSignalHandler::intSignalHandler(int _)
{
    char a = 1;
    ::write(sigintFd[0], &a, sizeof(a));
}

void UnixSignalHandler::termSignalHandler(int _) {
    char a = 1;
    ::write(sigtermFd[0], &a, sizeof(a));
}

void UnixSignalHandler::handleSigInt()
{
    qInfo() << "Caught SIGINT";

    m_snInt->setEnabled(false);
    char tmp;
    ::read(sigintFd[1], &tmp, sizeof(tmp));

    if (m_doOnTermination) m_doOnTermination();

    m_snInt->setEnabled(false);
}

void UnixSignalHandler::handleSigTerm()
{
    qInfo() << "Caught SIGTERM";

    m_snTerm->setEnabled(false);
    char tmp;
    ::read(sigtermFd[1], &tmp, sizeof(tmp));

    if (m_doOnTermination) m_doOnTermination();

    m_snTerm->setEnabled(false);
}

void UnixSignalHandler::setDoOnTermination(const std::function<void()> &fn)
{
    m_doOnTermination = fn;
}

void UnixSignalHandler::setup()
{
    struct sigaction sigint {
    };

    sigint.sa_handler = UnixSignalHandler::intSignalHandler;
    sigemptyset(&sigint.sa_mask);
    sigint.sa_flags = SA_RESTART;

    if (sigaction(SIGINT, &sigint, nullptr))
        qCritical() << "Could not set up SIGINT handler";

    struct sigaction sigterm {
    };

    sigint.sa_handler = UnixSignalHandler::termSignalHandler;
    sigemptyset(&sigint.sa_mask);
    sigint.sa_flags = SA_RESTART;

    if (sigaction(SIGTERM, &sigterm, nullptr))
        qCritical() << "Could not set up SIGINT handler";
}

#endif