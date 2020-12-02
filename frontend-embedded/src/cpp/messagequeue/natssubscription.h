#pragma once

#include <QObject>
#include <QSharedPointer>
#include <QString>
#include <QThread>

#include <nats.h>

#include "messages/decoders.h"

namespace MessageQueue
{

class NatsSubscription: public QObject
{
    Q_OBJECT
    Q_PROPERTY(MessageQueue::NatsConnection *connection READ connection WRITE setConnection NOTIFY connectionChanged)
    Q_PROPERTY(QString subject READ subject WRITE setSubject NOTIFY subjectChanged)

public:
    explicit NatsSubscription(QObject *parent = nullptr);
    ~NatsSubscription() override;

    Q_INVOKABLE void connectDecoder(Decoder *decoder) const;

    [[nodiscard]] QString subject() const;
    void setSubject(const QString &subject);
    [[nodiscard]] NatsConnection *connection() const;
    void setConnection(NatsConnection *connection);

public slots:
    void start();
    void stop();

signals:
    void messageReceived(const QSharedPointer<natsMsg *> &msg);
    void subjectChanged();
    void connectionChanged();

private:
    static void handleMessageReceived(natsConnection *nc, natsSubscription *sub, natsMsg *msg, void *closure);
    void handleMessageReceived(natsMsg *msg);

    NatsConnection *m_conn;
    QSharedPointer<natsConnection *> m_nc;
    natsSubscription *m_sub = nullptr;
    QString m_subject;
};

} // namespace MessageQueue
