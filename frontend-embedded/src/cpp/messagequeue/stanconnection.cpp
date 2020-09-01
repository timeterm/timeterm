#include "stanconnection.h"
#include "enums.h"
#include "stancallbackhandlersingleton.h"
#include "stansubscription.h"
#include "strings.h"

#include <QDebug>
#include <QtConcurrent/QtConcurrentRun>

namespace MessageQueue
{

StanConnection::StanConnection(QObject *parent)
    : QObject(parent)
    , m_options(nullptr)
{
    QObject::connect(this, &MessageQueue::StanConnection::setConnectionPrivate, this, &MessageQueue::StanConnection::setConnection);
}

StanConnection::~StanConnection()
{
    if (!m_stanConnection.isNull()) {
        StanCallbackHandlerSingleton::singleton()
            .removeConnectionLostHandler(*m_stanConnection);
    }
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

            QSharedPointer<stanConnection *> stanConnPtr(
                new stanConnection *(nullptr),
                [](stanConnection **ppConn) {
                    if (*ppConn != nullptr) {
                        stanConnection_Destroy(*ppConn);
                    }
                });

            auto connectionStatus = stanConnection_Connect(
                stanConnPtr.get(),
                clusterCstr.get(),
                clientIdCstr.get(),
                m_options->connectionOptions().get());

            updateStatus(NatsStatus::fromC(connectionStatus));
            if (connectionStatus != NATS_OK)
                return;
            qDebug() << "Connected";

            emit setConnectionPrivate(stanConnPtr, QPrivateSignal());

            StanCallbackHandlerSingleton::singleton().setConnectionLostHandler(
                *stanConnPtr,
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

NatsStatus::Enum StanConnection::subscribe(StanSubOptions *opts, stanSubscription **ppStanSub, QSharedPointer<stanConnection *> &spConn)
{
    stanSubOptions *pSubOptions = nullptr;
    auto buildStatus = opts->build(&pSubOptions);
    if (buildStatus != NatsStatus::Enum::Ok)
        return buildStatus;
    StanSubOptionsScopedPointer subOptions(pSubOptions);

    if (m_stanConnection.isNull())
        return NatsStatus::Enum::NotYetConnected;

    auto channelCstr = asUtf8CString(opts->channel());
    auto status = stanConnection_Subscribe(
        ppStanSub,                           // subscription (output parameter)
        *m_stanConnection,                   // connection
        channelCstr.get(),                   // channel
        StanCallbackHandlerSingleton::onMsg, // message handler
        nullptr,                             // message handler closure (not needed)
        subOptions.get());                   // subscription options
    if (status != NATS_OK)
        return NatsStatus::fromC(status);
    spConn.swap(m_stanConnection);

    qDebug() << "Subscribed to channel" << opts->channel();

    return NatsStatus::fromC(status);
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
    m_stanConnection = conn;
}

} // namespace MessageQueue
