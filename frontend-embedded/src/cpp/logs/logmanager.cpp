#include "logmanager.h"

LogManager::LogManager(QObject *parent)
    : QObject(parent)
{
}

LogManager *LogManager::singleton()
{
    static LogManager instance;
    return &instance;
}

void LogManager::setMessages(const QStringList &messages)
{
    QMutexLocker locker(&m_mut);

    if (messages != m_messages) {
        m_messages = messages;
        emit messagesChanged();
    }
}

QStringList LogManager::messages()
{
    QMutexLocker locker(&m_mut);
    QStringList copy = m_messages;

    return copy;
}

void LogManager::_handleMessage(QtMsgType type, const QMessageLogContext &context, const QString &buf)
{
    QMutexLocker locker(&m_mut);

    auto msg = qFormatLogMessage(type, context, buf);
    std::cerr << msg.toStdString() << std::endl;
    if (m_messages.length() == KEEP_MAX_AMOUNT_OF_MESSAGES) {
        m_messages.removeFirst();
    }
    m_messages.append(msg);

    // Unlock the locker because otherwise messagesChanged() being handled synchronously
    // could cause a deadlock (due to the handler wanting access to m_mut too).
    locker.unlock();

    emit messagesChanged();
}

void LogManager::handleMessage(QtMsgType type, const QMessageLogContext &context, const QString &buf)
{
    singleton()->_handleMessage(type, context, buf);
}
