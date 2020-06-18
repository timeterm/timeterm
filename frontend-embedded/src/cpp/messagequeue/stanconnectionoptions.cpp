#include "stanconnectionoptions.h"
#include "stancallbackhandlersingleton.h"
#include "strings.h"

namespace MessageQueue
{

StanConnectionOptions::StanConnectionOptions(QObject *parent)
    : QObject(parent)
    , m_connOptions(nullptr, stanConnOptions_Destroy)
{
    stanConnOptions *connOptions;
    auto s = stanConnOptions_Create(&connOptions);
    if (s == NATS_OK) {
        s = stanConnOptions_SetConnectionLostHandler(connOptions, StanCallbackHandlerSingleton::onConnLost, nullptr);
        if (s == NATS_OK) {
            m_connOptions.reset(connOptions);
        }
    }

    updateStatus(NatsStatus::fromC(s));

    if (s != NATS_OK)
        stanConnOptions_Destroy(connOptions);
}

NatsStatus::Enum StanConnectionOptions::setConnectionWait(qint64 wait)
{
    auto s = NatsStatus::fromC(stanConnOptions_SetConnectionWait(m_connOptions.get(), wait));
    updateStatus(s);
    return s;
}

NatsStatus::Enum StanConnectionOptions::setDiscoveryPrefix(const QString &prefix)
{
    auto discoveryPrefixCstr = asUtf8CString(prefix);

    auto s = NatsStatus::fromC(stanConnOptions_SetDiscoveryPrefix(m_connOptions.get(), discoveryPrefixCstr.get()));
    updateStatus(s);
    return s;
}

NatsStatus::Enum StanConnectionOptions::setMaxPubAcksInflight(int maxPubAcksInflight, float percentage)
{
    auto s = NatsStatus::fromC(stanConnOptions_SetMaxPubAcksInflight(m_connOptions.get(), maxPubAcksInflight, percentage));
    updateStatus(s);
    return s;
}

NatsStatus::Enum StanConnectionOptions::setPings(int interval, int maxOut)
{
    auto s = NatsStatus::fromC(stanConnOptions_SetPings(m_connOptions.get(), interval, maxOut));
    updateStatus(s);
    return s;
}

NatsStatus::Enum StanConnectionOptions::setPubAckWait(qint64 ms)
{
    auto s = NatsStatus::fromC(stanConnOptions_SetPubAckWait(m_connOptions.get(), ms));
    updateStatus(s);
    return s;
}

NatsStatus::Enum StanConnectionOptions::setUrl(const QString &url)
{
    auto urlCstr = asUtf8CString(url);

    auto s = NatsStatus::fromC(stanConnOptions_SetURL(m_connOptions.get(), urlCstr.get()));
    updateStatus(s);
    return s;
}

QSharedPointer<stanConnOptions> StanConnectionOptions::connectionOptions()
{
    return m_connOptions;
}

NatsStatus::Enum StanConnectionOptions::lastStatus() const
{
    return m_lastStatus;
}

void StanConnectionOptions::updateStatus(NatsStatus::Enum s)
{
    m_lastStatus = s;
    if (s == NatsStatus::Enum::Ok)
        return;

    const char *text = natsStatus_GetText(NatsStatus::asC(s));
    auto statusStr = QString::fromLocal8Bit(text);
    emit errorOccurred(s, statusStr);
}

NatsStatus::Enum StanConnectionOptions::setNatsOptions(NatsOptions *opts)
{
    opts->setParent(this);

    auto s = NatsStatus::fromC(stanConnOptions_SetNATSOptions(m_connOptions.get(), opts->options().get()));
    updateStatus(s);
    return s;
}

} // namespace MessageQueue
