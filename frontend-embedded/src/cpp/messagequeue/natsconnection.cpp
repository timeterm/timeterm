#include "natsconnection.h"
#include "strings.h"

#include <QtConcurrent/QtConcurrentRun>

namespace MessageQueue
{

NatsConnection::NatsConnection(QObject *parent)
    : QObject(parent)
{
    QObject::connect(this, &MessageQueue::NatsConnection::setHolderPrivate, this, &MessageQueue::NatsConnection::setHolder);
}

void NatsConnection::setOptions(NatsOptions *options)
{
    if (options != nullptr && options != m_options) {
        if (m_options != nullptr)
            m_options->deleteLater();
        options->setParent(this);
        m_options = options;
        emit optionsChanged();
    }
}

NatsOptions *NatsConnection::options() const
{
    return m_options;
}

void NatsConnection::connect()
{
    natsOptions *pOpts = nullptr;

    auto optsStatus = m_options->build(&pOpts);
    updateStatus(optsStatus);
    if (optsStatus != NatsStatus::Enum::Ok) {
        qCritical() << "Could not create NATS options";
        return;
    }

    optsStatus = NatsStatus::fromC(natsOptions_SetDisconnectedCB(pOpts, NatsConnectionHolder::connectionLostCB, this));
    updateStatus(optsStatus);
    if (optsStatus != NatsStatus::Enum::Ok) {
        qCritical() << "Could not set connection lost callback handler";
        return;
    }

    QSharedPointer<natsOptions *> opts(
        new natsOptions *(pOpts),
        [](natsOptions **ppOpts) {
            if (ppOpts != nullptr) {
                if (*ppOpts != nullptr) {
                    natsOptions_Destroy(*ppOpts);
                }
                delete ppOpts;
            }
        });

    QtConcurrent::run(
        [this, opts]() {
            natsConnection *conn = nullptr;
            auto connectionStatus = natsConnection_Connect(&conn, *opts);
            updateStatus(NatsStatus::fromC(connectionStatus));
            if (connectionStatus != NATS_OK)
                return;
            qDebug() << "Connected";

            QSharedPointer<NatsConnectionHolder> holder(new NatsConnectionHolder(conn));

            emit setHolderPrivate(holder, QPrivateSignal());
            emit connected();
        });
}

NatsStatus::Enum NatsConnection::lastStatus()
{
    return m_lastStatus;
}

void NatsConnection::updateStatus(NatsStatus::Enum s)
{
    if (s != m_lastStatus) {
        m_lastStatus = s;
        emit lastStatusChanged();
    }

    if (s == NatsStatus::Enum::Ok)
        return;

    const char *text = natsStatus_GetText(NatsStatus::asC(s));
    auto statusStr = QString::fromLocal8Bit(text);
    emit errorOccurred(s, statusStr);
}

QSharedPointer<NatsConnectionHolder> NatsConnection::getHolder() const
{
    return m_holder;
}

void NatsConnection::setHolder(const QSharedPointer<NatsConnectionHolder> &holder)
{
    if (holder != m_holder) {
        m_holder = holder;
        QObject::connect(holder.get(), &NatsConnectionHolder::connectionLost, this, &NatsConnection::connectionLost);
        emit holderChanged();
    }
}

NatsConnectionHolder::NatsConnectionHolder(natsConnection *conn, QObject *parent)
    : QObject(parent)
    , m_nc(conn)
{
}

natsConnection *NatsConnectionHolder::getConnection() const
{
    return m_nc;
}

void NatsConnectionHolder::updateStatus(NatsStatus::Enum s)
{
    if (s != m_lastStatus) {
        m_lastStatus = s;
        emit lastStatusChanged();
    }

    if (s == NatsStatus::Enum::Ok)
        return;

    const char *text = natsStatus_GetText(NatsStatus::asC(s));
    auto statusStr = QString::fromLocal8Bit(text);
    emit errorOccurred(s, statusStr);
}

void NatsConnectionHolder::connectionLostCB(natsConnection *conn, void *closure)
{
    if (closure != nullptr)
        static_cast<NatsConnectionHolder *>(closure)->connectionLostCB();
}

void NatsConnectionHolder::connectionLostCB()
{
    emit connectionLost();
}

} // namespace MessageQueue
