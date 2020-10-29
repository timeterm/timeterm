#include "ttlogmanager.h"

TtLogManager::TtLogManager(QObject *parent)
    : QObject(parent)
{
}

TtLogManager *TtLogManager::singleton()
{
    static TtLogManager instance;
    return &instance;
}

void TtLogManager::setMessages(const QStringList &messages)
{
    QMutexLocker locker(&m_mut);

    if (messages != m_messages) {
        m_messages = messages;
        emit messagesChanged();
    }
}

QStringList TtLogManager::messages()
{
    QMutexLocker locker(&m_mut);
    QStringList copy = m_messages;

    return copy;
}

void TtLogManager::_handleMessage(QtMsgType type, const QMessageLogContext &context, const QString &buf)
{
    QMutexLocker locker(&m_mut);

    auto msg = qFormatLogMessage(type, context, buf);
    std::cerr << msg.toStdString() << std::endl;
    m_messages.append(msg);

    locker.unlock();

    emit messagesChanged();
}

void TtLogManager::handleMessage(QtMsgType type, const QMessageLogContext &context, const QString &buf)
{
    singleton()->_handleMessage(type, context, buf);
}
