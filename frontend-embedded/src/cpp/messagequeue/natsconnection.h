#pragma once

#include "enums.h"
#include "natsoptions.h"

#include <QObject>

namespace MessageQueue
{

class NatsConnectionHolder: public QObject
{
    Q_OBJECT
    Q_PROPERTY(MessageQueue::NatsStatus::Enum lastStatus READ lastStatus)

public:
    explicit NatsConnectionHolder(natsConnection *conn, QObject *parent = nullptr);
    ~NatsConnectionHolder() override;

    [[nodiscard]] NatsStatus::Enum lastStatus();
    [[nodiscard]] natsConnection *getConnection() const;

signals:
    void errorOccurred(MessageQueue::NatsStatus::Enum s, const QString &message);
    void optionsChanged();
    void setConnectionPrivate(const QSharedPointer<natsConnection *> &conn, QPrivateSignal);
    void lastStatusChanged();
    void connectionLost();

private:
    friend class NatsConnection;

    void updateStatus(NatsStatus::Enum s);

    static void connectionLostCB(natsConnection *conn, void *closure);
    void connectionLostCB();

    NatsStatus::Enum m_lastStatus = NatsStatus::Enum::Ok;
    natsConnection *m_nc;
};

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
    [[nodiscard]] QSharedPointer<NatsConnectionHolder> getHolder() const;

    Q_INVOKABLE void connect();

signals:
    void errorOccurred(MessageQueue::NatsStatus::Enum s, const QString &message);
    void optionsChanged();
    void connected();
    void setHolderPrivate(const QSharedPointer<MessageQueue::NatsConnectionHolder> &holder, QPrivateSignal);
    void lastStatusChanged();
    void connectionLost();
    void holderChanged();

private slots:
    void setHolder(const QSharedPointer<NatsConnectionHolder> &holder);

private:
    void updateStatus(NatsStatus::Enum s);

    NatsStatus::Enum m_lastStatus = NatsStatus::Enum::Ok;
    NatsOptions *m_options = nullptr;
    QSharedPointer<NatsConnectionHolder> m_holder;
};

} // namespace MessageQueue

Q_DECLARE_METATYPE(QSharedPointer<MessageQueue::NatsConnectionHolder>)