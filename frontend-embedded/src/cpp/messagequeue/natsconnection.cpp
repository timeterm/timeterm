#include "natsconnection.h"
#include "strings.h"

#include <QtConcurrent/QtConcurrentRun>

namespace MessageQueue
{

NatsConnection::NatsConnection(QObject *parent)
    : QObject(parent)
{
    QObject::connect(this, &MessageQueue::NatsConnection::setConnectionPrivate, this, &MessageQueue::NatsConnection::setConnection);
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

    optsStatus = NatsStatus::fromC(natsOptions_SetDisconnectedCB(pOpts, NatsConnection::connectionLostCB, this));
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
            QSharedPointer<natsConnection *> natsConnPtr(
                new natsConnection *(nullptr),
                [](natsConnection **ppConn) {
                    if (ppConn != nullptr) {
                        if (*ppConn != nullptr) {
                            natsConnection_Destroy(*ppConn);
                        }
                        delete ppConn;
                    }
                });

            auto connectionStatus = natsConnection_Connect(natsConnPtr.get(), *opts);
            updateStatus(NatsStatus::fromC(connectionStatus));
            if (connectionStatus != NATS_OK)
                return;
            qDebug() << "Connected";

            emit setConnectionPrivate(natsConnPtr, QPrivateSignal());
            emit connected();
        });
}

void NatsConnection::setConnection(const QSharedPointer<natsConnection *> &conn)
{
    m_natsConnection = conn;
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

QSharedPointer<natsConnection *> NatsConnection::getConnection() const
{
    return m_natsConnection;
}

void NatsConnection::connectionLostCB(natsConnection *nc, void *closure)
{
    if (closure != nullptr)
        static_cast<NatsConnection*>(closure)->connectionLostCB();
}

void NatsConnection::connectionLostCB() {
    emit connectionLost();
}

} // namespace MessageQueue
