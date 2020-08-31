#include "stansuboptions.h"
#include "enums.h"
#include "strings.h"

namespace MessageQueue
{

NatsStatus::Enum newStanSubOptions(QSharedPointer<stanSubOptions> &ptr)
{
    stanSubOptions *stanSubOpts = nullptr;
    auto s = stanSubOptions_Create(&stanSubOpts);
    if (s == NATS_OK)
        ptr.reset(stanSubOpts, stanSubOptions_Destroy);

    if (s != NATS_OK)
        stanSubOptions_Destroy(stanSubOpts);

    return NatsStatus::fromC(s);
}

StanSubOptions::StanSubOptions(QObject *parent)
    : QObject(parent)
{
}

#define CHECK_NATS_STATUS(status)                           \
    do {                                                    \
        if (s != NATS_OK) return NatsStatus::fromC(status); \
    } while (0)

NatsStatus::Enum StanSubOptions::build(stanSubOptions **ppSubOpts)
{
    auto status = stanSubOptions_Create(ppSubOpts);
    if (status != NATS_OK)
        return NatsStatus::fromC(status);

    auto s = configureSubOpts(*ppSubOpts);
    if (s != NatsStatus::Enum::Ok) {
        stanSubOptions_Destroy(*ppSubOpts);
        return s;
    }

    return s;
}

NatsStatus::Enum StanSubOptions::configureSubOpts(stanSubOptions *pSubOpts)
{
    natsStatus s;

    if (m_durableName != "") {
        auto durableNameCstr = asUtf8CString(m_durableName);
        s = stanSubOptions_SetDurableName(pSubOpts, durableNameCstr.get());
        CHECK_NATS_STATUS(s);
    }

    if (m_deliverAllAvailable) {
        s = stanSubOptions_DeliverAllAvailable(pSubOpts);
        CHECK_NATS_STATUS(s);
    }

    if (m_startWithLastReceived) {
        s = stanSubOptions_StartWithLastReceived(pSubOpts);
        CHECK_NATS_STATUS(s);
    }

    if (m_isStartAtSequenceSet) {
        s = stanSubOptions_StartAtSequence(pSubOpts, m_startAtSequence);
        CHECK_NATS_STATUS(s);
    }

    s = stanSubOptions_SetManualAckMode(pSubOpts, m_manualAckMode);
    CHECK_NATS_STATUS(s);

    if (m_isMaxInflightSet) {
        s = stanSubOptions_SetMaxInflight(pSubOpts, m_maxInflight);
        CHECK_NATS_STATUS(s);
    }

    if (m_isAckWaitMsSet) {
        s = stanSubOptions_SetAckWait(pSubOpts, m_ackWaitMs);
        CHECK_NATS_STATUS(s);
    }

    return NatsStatus::fromC(s);
}

QString StanSubOptions::durableName() const
{
    return m_durableName;
}

void StanSubOptions::setDurableName(const QString &durableName)
{
    if (durableName != m_durableName) {
        m_durableName = durableName;
        emit durableNameChanged();
    }
}

bool StanSubOptions::deliverAllAvailable() const
{
    return m_deliverAllAvailable;
}

void StanSubOptions::setDeliverAllAvailable(bool deliverAllAvailable)
{
    if (deliverAllAvailable != m_deliverAllAvailable) {
        m_deliverAllAvailable = deliverAllAvailable;
        emit deliverAllAvailableChanged();
    }
}

bool StanSubOptions::startWithLastReceived() const
{
    return m_startWithLastReceived;
}

void StanSubOptions::setStartWithLastReceived(bool startWithLastReceived)
{
    if (startWithLastReceived != m_startWithLastReceived) {
        m_startWithLastReceived = startWithLastReceived;
        emit startWithLastReceivedChanged();
    }
}

quint64 StanSubOptions::startAtSequence() const
{
    return m_startAtSequence;
}

void StanSubOptions::setStartAtSequence(bool startAtSequence)
{
    if (startAtSequence != m_startAtSequence) {
        m_startAtSequence = startAtSequence;
        m_isStartAtSequenceSet = true;
        emit startAtSequenceChanged();
    }
}

bool StanSubOptions::manualAckMode() const
{
    return m_manualAckMode;
}

void StanSubOptions::setManualAckMode(bool manualAckMode)
{
    if (manualAckMode != m_manualAckMode) {
        m_manualAckMode = manualAckMode;
        emit manualAckModeChanged();\
    }
}

int StanSubOptions::maxInflight() const
{
    return m_maxInflight;
}

void StanSubOptions::setMaxInflight(int maxInflight)
{
    if (maxInflight != m_maxInflight) {
        m_maxInflight = maxInflight;
        m_isMaxInflightSet = true;
        emit maxInflightChanged();
    }
}

qint64 StanSubOptions::ackWaitMs() const
{
    return m_ackWaitMs;
}

void StanSubOptions::setAckWaitMs(qint64 ackWaitMs)
{
    if (ackWaitMs != m_ackWaitMs) {
        m_ackWaitMs = ackWaitMs;
        m_isAckWaitMsSet = true;
        emit ackWaitMsChanged();
    }
}

QString StanSubOptions::channel() const
{
    return m_channel;
}

void StanSubOptions::setChannel(const QString &channel)
{
    if (channel != m_channel) {
        m_channel = channel;
        emit channelChanged();
    }
}

} // namespace MessageQueue