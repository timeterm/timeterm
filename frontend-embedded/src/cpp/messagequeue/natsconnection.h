#pragma once

#include "enums.h"
#include "natsoptions.h"

#include <QObject>

namespace MessageQueue
{

class NatsConnection: public QObject
{
    Q_OBJECT
    Q_PROPERTY(MessageQueue::NatsStatus::Enum lastStatus READ lastStatus)
    Q_PROPERTY(MessageQueue::NatsOptions *options READ options WRITE setOptions NOTIFY optionsChanged)

public:
    explicit NatsConnection(QObject *parent = nullptr);

    [[nodiscard]] NatsStatus::Enum lastStatus();
    void setOptions(NatsOptions *options);
    [[nodiscard]] NatsOptions *options() const;
    [[nodiscard]] QSharedPointer<natsConnection *> getConnection() const;

    Q_INVOKABLE void connect();
    NatsStatus::Enum subscribe(const QString &topic, natsSubscription **ppNatsSub, QSharedPointer<natsConnection *> &spConn);

signals:
    void errorOccurred(MessageQueue::NatsStatus::Enum s, const QString &message);
    void optionsChanged();
    void connected();
    void setConnectionPrivate(const QSharedPointer<natsConnection *> &conn, QPrivateSignal);
    void lastStatusChanged();
    void connectionLost();

private slots:
    void setConnection(const QSharedPointer<natsConnection *> &conn);

private:
    void updateStatus(NatsStatus::Enum s);

    static void connectionLostCB(natsConnection *conn, void *closure);
    void connectionLostCB();

    NatsStatus::Enum m_lastStatus = NatsStatus::Enum::Ok;
    NatsOptions *m_options = nullptr;
    QSharedPointer<natsConnection *> m_natsConnection;
};

} // namespace MessageQueue

Q_DECLARE_METATYPE(QSharedPointer<natsConnection *>)