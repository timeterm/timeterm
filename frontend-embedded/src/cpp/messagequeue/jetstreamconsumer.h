#pragma once

#include "enums.h"
#include "messagequeue/messages/disowntokenmessage.h"
#include "messagequeue/messages/retrievenewtokenmessage.h"
#include "natsconnection.h"

#include <QObject>
#include <QSharedPointer>
#include <QThread>
#include <QtCore/QMutex>
#include <QtCore/QTimer>

#include <nats.h>
#include <timeterm_proto/messages.pb.h>

namespace MessageQueue
{

class JetStreamPullConsumerWorker: public QObject
{
    Q_OBJECT

public:
    JetStreamPullConsumerWorker(
        const QSharedPointer<natsConnection *> &conn,
        QString stream,
        QString consumerId,
        QObject *parent = nullptr);
    ~JetStreamPullConsumerWorker() override;

signals:
    void messageReceived(const QSharedPointer<natsMsg *> &msg);

public slots:
    void start();
    void stop();
    void getNextMessage();

private:
    QTimer m_timer;
    QString m_stream;
    QString m_consumerId;
    QSharedPointer<natsConnection *> m_conn;
};

class JetStreamConsumer: public QObject
{
    Q_OBJECT
    Q_PROPERTY(MessageQueue::NatsConnection *connection READ connection WRITE setConnection NOTIFY connectionChanged)
    Q_PROPERTY(QString subject READ subject WRITE setSubject NOTIFY subjectChanged)
    Q_PROPERTY(QString stream READ stream WRITE setStream NOTIFY streamChanged)
    Q_PROPERTY(QString consumerId READ consumerId WRITE setConsumerId NOTIFY consumerIdChanged)
    Q_PROPERTY(MessageQueue::JetStreamConsumerType::Enum type READ type WRITE setType NOTIFY typeChanged)

public:
    explicit JetStreamConsumer(QObject *parent = nullptr);
    ~JetStreamConsumer() override;

    [[nodiscard]] QString subject() const;
    void setSubject(const QString &subject);
    [[nodiscard]] QString stream() const;
    void setStream(const QString &stream);
    [[nodiscard]] QString consumerId() const;
    void setConsumerId(const QString &consumerId);
    [[nodiscard]] NatsConnection *connection() const;
    void setConnection(NatsConnection *connection);
    [[nodiscard]] JetStreamConsumerType::Enum type() const;
    void setType(JetStreamConsumerType::Enum consumerType);

    Q_INVOKABLE void start();
    Q_INVOKABLE void stop();

signals:
    void connectionChanged();
    void subjectChanged();
    void streamChanged();
    void consumerIdChanged();
    void typeChanged();
    void disownTokenMessage(const MessageQueue::DisownTokenMessage &msg);
    void retrieveNewTokenMessage(const MessageQueue::RetrieveNewTokenMessage &msg);

private slots:
    void handleMessageSP(const QSharedPointer<natsMsg *> &msg);

private:
    void handleMessage(natsMsg *msg);
    void handleRetrieveNewTokenProto(const timeterm_proto::messages::RetrieveNewTokenMessage &msg);
    void handleDisownTokenProto(const timeterm_proto::messages::DisownTokenMessage &msg);

    QString m_subject;
    QString m_stream;
    QString m_consumerId;
    JetStreamConsumerType::Enum m_type = JetStreamConsumerType::Pull;
    QThread m_workerThread;
    NatsConnection *m_connection = nullptr;
};

} // namespace MessageQueue

Q_DECLARE_METATYPE(QSharedPointer<natsSubscription *>)
Q_DECLARE_METATYPE(QSharedPointer<natsMsg *>)