#include "natsconnection.h"
#include "natscallbackhandlersingleton.h"
#include "strings.h"
#include <QtConcurrent/QtConcurrentRun>

namespace MessageQueue
{

NatsConnection::NatsConnection(QObject *parent)
    : QObject(parent)
{
    QObject::connect(this, &MessageQueue::NatsConnection::setConnectionPrivate, this, &MessageQueue::NatsConnection::setConnection);
}

NatsConnection::~NatsConnection()
{
    if (!m_natsConnection.isNull()) {
    }
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
    QtConcurrent::run(
        [this]() {
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

            natsOptions *opts = nullptr;
            // TODO(rutgerbrf): should m_options be passed to the closure as a parameter?
            auto optsStatus = m_options->build(&opts);
            updateStatus(optsStatus);
            if (optsStatus != NatsStatus::Enum::Ok)
                return;

            auto connectionStatus = natsConnection_Connect(natsConnPtr.get(), opts);
            updateStatus(NatsStatus::fromC(connectionStatus));
            if (connectionStatus != NATS_OK)
                return;
            qDebug() << "Connected";

            emit setConnectionPrivate(natsConnPtr, QPrivateSignal());

            // TODO(rutgerbrf): check if this is the right approach for a plain NATS connection.
//            NatsCallbackHandlerSingleton::singleton().setConnectionLostHandler(
//                *natsConnPtr,
//                [this](const char*msg) {
//                    emit connectionLost();
//                });

            emit connected();
        });
}

NatsStatus::Enum NatsConnection::subscribe(const QString &topic, natsSubscription **ppNatsSub, QSharedPointer<natsConnection *> &spConn)
{
    if (m_natsConnection.isNull())
        return NatsStatus::Enum::NotYetConnected;

    auto topicCstr = asUtf8CString(topic);
    auto status = natsConnection_Subscribe(
        ppNatsSub,
        *m_natsConnection,
        topicCstr.get(),
        NatsCallbackHandlerSingleton::onMsg,
        nullptr);
    if (status != NATS_OK)
        return NatsStatus::fromC(status);
    spConn.swap(m_natsConnection);

    qDebug() << "Subscribed to topic" << topic;

    return NatsStatus::fromC(status);
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

} // namespace MessageQueue
