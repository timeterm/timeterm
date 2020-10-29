#pragma once

#include <QMutex>
#include <QObject>
#include <iostream>
#include <util/scopeguard.h>

class TtLogManager: public QObject
{
    Q_OBJECT
    Q_PROPERTY(QStringList messages READ messages WRITE setMessages NOTIFY messagesChanged)

public:
    static TtLogManager *singleton();

    void setMessages(const QStringList &messages);

    [[nodiscard]] QStringList messages();

private:
    explicit TtLogManager(QObject *parent = nullptr);

    void _handleMessage(QtMsgType type, const QMessageLogContext &context,
                        const QString &buf);

    QMutex m_mut;
    QStringList m_messages;

signals:
    void messagesChanged();

public:
    TtLogManager(TtLogManager const &) = delete;
    void operator=(TtLogManager const &) = delete;

    static void handleMessage(QtMsgType type, const QMessageLogContext &context,
                              const QString &buf);
};