#include "stanconnection.h"
#include "enums.h"
#include "stancallbackhandlersingleton.h"
#include "strings.h"

#include <QtConcurrent/QtConcurrentRun>

namespace MessageQueue
{

StanConnection::StanConnection(QObject *parent)
    : QObject(parent)
    , m_options(nullptr)
{
    QObject::connect(this, &MessageQueue::StanConnection::setConnectionPrivate, this, &MessageQueue::StanConnection::setConnection);
}

void StanConnection::connect()
{
    // We don't want connecting to the NATS Streaming server to block the user interface.
    // For that reason, we're not directly setting any properties, but using signals so the
    // event loop can send it to the thread that the object is actually running on.
    QtConcurrent::run(
        [this](const QString &cluster, const QString &clientId) {
            auto clusterCstr = asUtf8CString(cluster);
            auto clientIdCstr = asUtf8CString(clientId);

            QSharedPointer<stanConnection *> stanConnPtr(new stanConnection *(nullptr));

            auto connectionStatus = stanConnection_Connect(
                stanConnPtr.get(),
                clusterCstr.get(),
                clientIdCstr.get(),
                m_options->connectionOptions().get());

            updateStatus(NatsStatus::fromC(connectionStatus));
            if (connectionStatus != NATS_OK)
                return;

            emit setConnectionPrivate(stanConnPtr, QPrivateSignal());

            StanCallbackHandlerSingleton::singleton().setConnectionLostHandler(
                *stanConnPtr.get(),
                [this](const char *msg) {
                    emit connectionLost();
                });

            emit connected();
        },
        m_cluster, m_clientId);
}

NatsStatus::Enum StanConnection::lastStatus() const
{
    return m_lastStatus;
}

void StanConnection::updateStatus(NatsStatus::Enum s)
{
    m_lastStatus = s;
    if (s == NatsStatus::Enum::Ok)
        return;

    const char *text = natsStatus_GetText(NatsStatus::asC(s));
    auto statusStr = QString::fromLocal8Bit(text);
    emit errorOccurred(s, statusStr);
}

StanSubscription *StanConnection::subscribe(const QString &channel, StanSubOptions *opts)
{
    auto subOptions = opts->subOptions();
    auto channelCstr = asUtf8CString(channel);

    stanSubscription *subDest = nullptr;
    stanConnection_Subscribe(
        &subDest,                            // subscription (output parameter)
        m_stanConnection.get(),              // connection
        channelCstr.get(),                   // channel
        StanCallbackHandlerSingleton::onMsg, // message handler
        nullptr,                             // message handler closure (not needed)
        subOptions.get());                   // subscription options

    auto subWrapper = new StanSubscription(this);
    subWrapper->setSubscription(subDest);

    return subWrapper;
}

void StanConnection::setCluster(const QString &cluster)
{
    if (cluster != m_cluster) {
        m_cluster = cluster;
        emit clusterChanged();
    }
}

QString StanConnection::cluster() const
{
    return m_cluster;
}

void StanConnection::setClientId(const QString &clientId)
{
    if (clientId != m_clientId) {
        m_clientId = clientId;
        emit clientIdChanged();
    }
}

QString StanConnection::clientId() const
{
    return m_clientId;
}

void StanConnection::setConnectionOptions(StanConnectionOptions *options)
{
    if (options != m_options) {
        options->setParent(this);
        m_options = options;
        emit connectionOptionsChanged();
    }
}

StanConnectionOptions *StanConnection::connectionOptions() const
{
    return m_options;
}

void StanConnection::setConnection(const QSharedPointer<stanConnection *> &conn)
{
    m_stanConnection.reset(*conn.get());
}

} // namespace MessageQueue
