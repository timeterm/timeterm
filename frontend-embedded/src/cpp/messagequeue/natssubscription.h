#pragma once

#include <QObject>
#include <QSharedPointer>
#include <QString>
#include <QThread>

#include <nats.h>

#include "messages/decoders.h"
#include "natsconnection.h"

namespace MessageQueue
{

class NatsSubscription: public QObject
{
    Q_OBJECT
    Q_PROPERTY(MessageQueue::NatsStatus::Enum lastStatus READ lastStatus)
    Q_PROPERTY(QString subject READ subject WRITE setSubject NOTIFY subjectChanged)

public:
    explicit NatsSubscription(QObject *parent = nullptr);
    ~NatsSubscription() override;

    Q_INVOKABLE void connectDecoder(MessageQueue::Decoder *decoder) const;
    Q_INVOKABLE void useConnection(MessageQueue::NatsConnection *connection);

    [[nodiscard]] NatsStatus::Enum lastStatus();
    [[nodiscard]] QString subject() const;
    void setSubject(const QString &subject);

public slots:
    void start();
    void stop();

signals:
    void errorOccurred(MessageQueue::NatsStatus::Enum s, const QString &message);
    void messageReceived(const QSharedPointer<natsMsg *> &msg);
    void subjectChanged();
    void lastStatusChanged();

private:
    void updateStatus(NatsStatus::Enum s);

    static void handleMessageReceived(natsConnection *nc, natsSubscription *sub, natsMsg *msg, void *closure);
    void handleMessageReceived(natsMsg *msg);

    NatsStatus::Enum m_lastStatus = NatsStatus::Enum::Ok;
    QSharedPointer<NatsConnectionHolder> m_connHolder;
    natsSubscription *m_sub = nullptr;
    QString m_subject;
};

} // namespace MessageQueue
