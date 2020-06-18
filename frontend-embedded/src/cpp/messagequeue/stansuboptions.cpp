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
        ptr.reset(stanSubOpts);

    if (s != NATS_OK)
        stanSubOptions_Destroy(stanSubOpts);

    return NatsStatus::fromC(s);
}

StanSubOptions::StanSubOptions(QObject *parent)
    : QObject(parent)
    , m_subOptions(nullptr, stanSubOptions_Destroy)
{
    updateStatus(newStanSubOptions(m_subOptions));
}

QSharedPointer<stanSubOptions> StanSubOptions::subOptions()
{
    return m_subOptions;
}

NatsStatus::Enum StanSubOptions::setDurableName(const QString &durableName)
{
    auto durableNameCstr = asUtf8CString(durableName);

    auto s = NatsStatus::fromC(stanSubOptions_SetDurableName(m_subOptions.get(), durableNameCstr.get()));
    updateStatus(s);
    return s;
}

NatsStatus::Enum StanSubOptions::deliverAllAvailable()
{
    auto s = NatsStatus::fromC(stanSubOptions_DeliverAllAvailable(m_subOptions.get()));
    updateStatus(s);
    return s;
}

NatsStatus::Enum StanSubOptions::startWithLastReceived()
{
    auto s = NatsStatus::fromC(stanSubOptions_StartWithLastReceived(m_subOptions.get()));
    updateStatus(s);
    return s;
}

NatsStatus::Enum StanSubOptions::startAtSequence(quint64 sequence)
{
    auto s = NatsStatus::fromC(stanSubOptions_StartAtSequence(m_subOptions.get(), sequence));
    updateStatus(s);
    return s;
}

NatsStatus::Enum StanSubOptions::setManualAckMode(bool manualAck)
{
    auto s = NatsStatus::fromC(stanSubOptions_SetManualAckMode(m_subOptions.get(), manualAck));
    updateStatus(s);
    return s;
}

NatsStatus::Enum StanSubOptions::setMaxInflight(int inflight)
{
    auto s = NatsStatus::fromC(stanSubOptions_SetMaxInflight(m_subOptions.get(), inflight));
    updateStatus(s);
    return s;
}

NatsStatus::Enum StanSubOptions::setAckWait(qint64 ms)
{
    auto s = NatsStatus::fromC(stanSubOptions_SetAckWait(m_subOptions.get(), ms));
    updateStatus(s);
    return s;
}

NatsStatus::Enum StanSubOptions::lastStatus() const
{
    return m_lastStatus;
}

void StanSubOptions::updateStatus(NatsStatus::Enum s)
{
    m_lastStatus = s;
    if (s == NatsStatus::Enum::Ok)
        return;

    const char *text = natsStatus_GetText(NatsStatus::asC(s));
    auto statusStr = QString::fromLocal8Bit(text);
    emit errorOccurred(s, statusStr);
}

} // namespace MessageQueue