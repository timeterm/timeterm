#pragma once

#include <QObject>

#include "enums.h"
#include "natsoptions.h"

namespace MessageQueue
{

class NatsConnection: public QObject
{
    Q_OBJECT
    Q_PROPERTY(MessageQueue::NatsStatus::Enum lastStatus READ lastStatus)
    Q_PROPERTY(MessageQueue::NatsOptions *options READ options WRITE setOptions NOTIFY optionsChanged)

public:
    explicit NatsConnection(QObject *parent = nullptr);
    ~NatsConnection() override;

    [[nodiscard]] NatsStatus::Enum lastStatus();
    void setOptions(NatsOptions *options);
    [[nodiscard]] NatsOptions *options() const;

    Q_INVOKABLE void connect();
    NatsStatus::Enum subscribe(const QString &topic, natsSubscription **ppNatsSub, QSharedPointer<natsConnection *>& spConn);

signals:
    void errorOccurred(MessageQueue::NatsStatus::Enum s, const QString &message);
    void optionsChanged();
    void connected();
    // void connectionLost();
    void setConnectionPrivate(const QSharedPointer<natsConnection *>&conn, QPrivateSignal);
    void lastStatusChanged();

private slots:
    void setConnection(const QSharedPointer<natsConnection *> &conn);

private:
    void updateStatus(NatsStatus::Enum s);

    NatsStatus::Enum m_lastStatus = NatsStatus::Enum::Ok;
    NatsOptions *m_options = nullptr;
    QSharedPointer<natsConnection *> m_natsConnection;
};

} // namespace MessageQueue

Q_DECLARE_METATYPE(QSharedPointer<natsConnection *>)