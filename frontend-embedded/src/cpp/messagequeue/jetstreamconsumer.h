#pragma once

#include "enums.h"
#include "natsconnection.h"
#include "natssubscription.h"

#include <QObject>
#include <QSharedPointer>
#include <QThread>
#include <QtCore/QMutex>
#include <QtCore/QTimer>

#include <nats.h>
#include <timeterm_proto/mq/mq.pb.h>
#include "messages/decoders.h"

namespace MessageQueue
{

class JetStreamPullConsumerWorker: public QObject
{
    Q_OBJECT

public:
    JetStreamPullConsumerWorker(
        const QSharedPointer<NatsConnectionHolder> &holder,
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
    QSharedPointer<NatsConnectionHolder> m_connHolder;
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

    Q_INVOKABLE void connectDecoder(MessageQueue::Decoder *decoder) const;

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
    void messageReceived(const QSharedPointer<natsMsg *> &msg);

private slots:
    void handleMessage(const QSharedPointer<natsMsg *> &msg);

private:
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
