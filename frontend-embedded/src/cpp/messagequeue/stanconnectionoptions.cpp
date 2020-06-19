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
            m_connOptions.reset(connOptions, stanConnOptions_Destroy);
        }
    }

    updateStatus(NatsStatus::fromC(s));

    if (s != NATS_OK)
        stanConnOptions_Destroy(connOptions);
}

void StanConnectionOptions::setConnectionWait(int wait)
{
    auto s = NatsStatus::fromC(stanConnOptions_SetConnectionWait(m_connOptions.get(), wait));
    updateStatus(s);
    if (s == NatsStatus::Enum::Ok) {
        if (wait != m_connectionWait) {
            m_connectionWait = wait;
            emit connectionWaitChanged();
        }
    }
}

void StanConnectionOptions::setDiscoveryPrefix(const QString &prefix)
{
    auto discoveryPrefixCstr = asUtf8CString(prefix);

    auto s = NatsStatus::fromC(stanConnOptions_SetDiscoveryPrefix(m_connOptions.get(), discoveryPrefixCstr.get()));
    updateStatus(s);
    if (s == NatsStatus::Enum::Ok) {
        if (prefix != m_discoveryPrefix) {
            m_discoveryPrefix = prefix;
            emit discoveryPrefixChanged();
        }
    }
}

void StanConnectionOptions::setMaxPubAcksInflight(int maxPubAcksInflight)
{
    m_maxPubAcksInflight = maxPubAcksInflight;

    if (updateMaxPubAcksInflight() == NatsStatus::Enum::Ok) {
        if (maxPubAcksInflight != m_maxPubAcksInflight) {
            m_maxPubAcksInflight = maxPubAcksInflight;
            emit maxPubAcksInflightChanged();
        }
    }
}

void StanConnectionOptions::setMaxPubAcksInflightPercentage(float percentage)
{
    m_maxPubAcksInflightPercentage = percentage;

    if (updateMaxPubAcksInflight() == NatsStatus::Enum::Ok) {
        if (percentage != m_maxPubAcksInflightPercentage) {
            m_maxPubAcksInflightPercentage = percentage;
            emit maxPubAcksInflightPercentageChanged();
        }
    }
}

NatsStatus::Enum StanConnectionOptions::updateMaxPubAcksInflight() {
    auto s = NatsStatus::fromC(stanConnOptions_SetMaxPubAcksInflight(m_connOptions.get(), m_maxPubAcksInflight, m_maxPubAcksInflightPercentage));
    updateStatus(s);
    return s;
}

void StanConnectionOptions::setPingsInterval(int interval)
{
    m_pingsInterval = interval;

    if (updatePings() == NatsStatus::Enum::Ok) {
        if (interval != m_pingsInterval) {
            m_pingsInterval = interval;
            emit pingsIntervalChanged();
        }
    }
}

void StanConnectionOptions::setPingsMaxOut(int maxOut)
{
    m_pingsMaxOut = maxOut;

    if (updatePings() == NatsStatus::Enum::Ok) {
        if (maxOut != m_pingsMaxOut) {
            m_pingsMaxOut = maxOut;
            emit pingsMaxOutChanged();
        }
    }
}

NatsStatus::Enum StanConnectionOptions::updatePings()
{
    auto s = NatsStatus::fromC(stanConnOptions_SetPings(m_connOptions.get(), m_pingsInterval, m_pingsMaxOut));
    updateStatus(s);
    return s;
}

void StanConnectionOptions::setPubAckWait(int ms)
{
    auto s = NatsStatus::fromC(stanConnOptions_SetPubAckWait(m_connOptions.get(), ms));
    updateStatus(s);
    if (s == NatsStatus::Enum::Ok) {
        if (ms != m_pubAckWait) {
            m_pubAckWait = ms;
            emit pubAckWaitChanged();
        }
    }
}

void StanConnectionOptions::setUrl(const QString &url)
{
    auto urlCstr = asUtf8CString(url);

    auto s = NatsStatus::fromC(stanConnOptions_SetURL(m_connOptions.get(), urlCstr.get()));
    updateStatus(s);
    if (s == NatsStatus::Enum::Ok) {
        if (url != m_url) {
            m_url = url;
            emit urlChanged();
        }
    }
}

void StanConnectionOptions::setNatsOptions(NatsOptions *opts)
{
    opts->setParent(this);

    auto s = NatsStatus::fromC(stanConnOptions_SetNATSOptions(m_connOptions.get(), opts->options().get()));
    updateStatus(s);
    if (s == NatsStatus::Enum::Ok) {
        if (opts != m_natsOptions) {
            m_natsOptions = opts;
            emit natsOptionsChanged();
        }
    }
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

NatsOptions *StanConnectionOptions::natsOptions() const
{
    return m_natsOptions;
}

int StanConnectionOptions::connectionWait() const
{
    return m_connectionWait;
}

QString StanConnectionOptions::discoveryPrefix() const
{
    return m_discoveryPrefix;
}

int StanConnectionOptions::maxPubAcksInflight() const
{
    return m_maxPubAcksInflight;
}

float StanConnectionOptions::maxPubAcksInflightPercentage() const
{
    return m_maxPubAcksInflightPercentage;
}

int StanConnectionOptions::pingsInterval() const
{
    return m_pingsInterval;
}

int StanConnectionOptions::pingsMaxOut() const
{
    return m_pingsMaxOut;
}

int StanConnectionOptions::pubAckWait() const
{
    return m_pubAckWait;
}

QString StanConnectionOptions::url() const
{
    return m_url;
}

} // namespace MessageQueue
