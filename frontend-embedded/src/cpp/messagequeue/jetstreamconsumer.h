#pragma once

#include <QObject>
#include <QSharedPointer>

#include "enums.h"
#include "natsconnection.h"
#include <nats.h>
#include <src/cpp/messagequeue/messages/disowntokenmessage.h>
#include <src/cpp/messagequeue/messages/retrievenewtokenmessage.h>
#include <timeterm_proto/messages.pb.h>

namespace MessageQueue
{

class JetStreamConsumer: public QObject
{
    Q_OBJECT
    Q_PROPERTY(MessageQueue::NatsConnection *target READ target WRITE setTarget NOTIFY targetChanged)
    Q_PROPERTY(QString subject READ subject WRITE setSubject NOTIFY subjectChanged)
    Q_PROPERTY(MessageQueue::NatsStatus::Enum lastStatus READ lastStatus NOTIFY lastStatusChanged)

public:
    explicit JetStreamConsumer(QObject *parent = nullptr);
    ~JetStreamConsumer() override;

    [[nodiscard]] NatsStatus::Enum lastStatus() const;
    [[nodiscard]] QString subject() const;
    void setSubject(const QString &subject);
    [[nodiscard]] NatsConnection *target() const;
    void setTarget(NatsConnection *target);

    Q_INVOKABLE void subscribe();

signals:
    void targetChanged();
    void subjectChanged();
    void lastStatusChanged();
    void errorOccurred(MessageQueue::NatsStatus::Enum s, const QString &msg);
    void disownTokenMessage(const MessageQueue::DisownTokenMessage &msg);
    void retrieveNewTokenMessage(const MessageQueue::RetrieveNewTokenMessage &msg);

    void updateSubscription(const QSharedPointer<natsSubscription *> &sub, const QSharedPointer<natsConnection *> &spConn, QPrivateSignal);

private slots:
    void setSubscription(const QSharedPointer<natsSubscription *> &sub, const QSharedPointer<natsConnection *> &spConn);

private:
    void updateStatus(NatsStatus::Enum s);
    void handleMessage(natsMsg *msg);
    void handleRetrieveNewTokenProto(const timeterm_proto::messages::RetrieveNewTokenMessage &msg);
    void handleDisownTokenProto(const timeterm_proto::messages::DisownTokenMessage &msg);

    QString m_subject;
    natsSubscription *m_sub = nullptr;
    QSharedPointer<natsConnection *> m_dontDropConn;
    NatsConnection *m_target = nullptr;
    NatsStatus::Enum m_lastStatus = NatsStatus::Enum::Ok;
};

} // namespace MessageQueue

Q_DECLARE_METATYPE(QSharedPointer<natsSubscription *>)