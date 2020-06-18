#include "stanconnection.h"
#include "enums.h"
#include "stancallbackhandlersingleton.h"
#include "strings.h"

namespace MessageQueue
{

NatsStatus::Enum newNatsOptions(NatsOptionsScopedPointer &ptr)
{
    natsOptions *natsOpts = nullptr;
    auto s = natsOptions_Create(&natsOpts);
    if (s == NATS_OK)
        ptr.reset(natsOpts);
    return NatsStatus::fromC(s);
}

NatsStatus::Enum newStanConnOptions(StanConnOptionsScopedPointer &ptr, NatsOptionsScopedPointer &natsOpts)
{
    stanConnOptions *stanConnOpts = nullptr;
    auto s = stanConnOptions_Create(&stanConnOpts);
    if (s == NATS_OK) {
        s = stanConnOptions_SetNATSOptions(stanConnOpts, natsOpts.get());
        if (s == NATS_OK) {
            ptr.reset(stanConnOpts);
        }
    }
    return NatsStatus::fromC(s);
}

StanConnection::StanConnection(QObject *parent)
    : QObject(parent)
{
    updateStatus(newNatsOptions(m_natsOpts));
    if (m_lastStatus == NatsStatus::Enum::Ok)
        newStanConnOptions(m_connOpts, m_natsOpts);
}

void StanConnection::connect()
{
    auto cluster = asUtf8CString(m_cluster);
    auto clientId = asUtf8CString(m_clientId);
    stanConnection *stanConn = nullptr;

    auto s = NatsStatus::fromC(stanConnection_Connect(&stanConn, cluster.get(), clientId.get(), m_connOpts.get()));
    m_stanConnection.reset(stanConn);

    updateStatus(s);
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

void StanConnection::setUrl(const QString &url)
{
    if (url != m_url) {
        m_url = url;
        emit urlChanged();
    }
}

QString StanConnection::url() const
{
    return m_url;
}

} // namespace MessageQueue
