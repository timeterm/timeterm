#include "stansuboptions.h"
#include "enums.h"
#include "strings.h"

namespace MessageQueue
{

NatsStatus::Enum newStanSubOptions(QSharedPointer<stanSubOptions> &ptr)
{
    stanSubOptions *stanSubOpts = nullptr;
    auto s = static_cast<NatsStatus::Enum>(stanSubOptions_Create(&stanSubOpts));
    if (s == NatsStatus::Enum::Ok)
        ptr.reset(stanSubOpts);
    return s;
}

StanSubOptions::StanSubOptions(QObject *parent)
    : QObject(parent)
{
    // TODO: maybe don't ignore the error?
    newStanSubOptions(m_subOptions);
}

QSharedPointer<stanSubOptions> StanSubOptions::subOptions()
{
    return m_subOptions;
}

NatsStatus::Enum StanSubOptions::setDurableName(const QString &durableName)
{
    auto durableNameCstr = asUtf8CString(durableName);

    return NatsStatus::as(stanSubOptions_SetDurableName(m_subOptions.get(), durableNameCstr.get()));
}

NatsStatus::Enum StanSubOptions::deliverAllAvailable()
{
    return NatsStatus::as(stanSubOptions_DeliverAllAvailable(m_subOptions.get()));
}

NatsStatus::Enum StanSubOptions::startWithLastReceived()
{
    return NatsStatus::as(stanSubOptions_StartWithLastReceived(m_subOptions.get()));
}

NatsStatus::Enum StanSubOptions::startAtSequence(quint64 sequence)
{
    return NatsStatus::as(stanSubOptions_StartAtSequence(m_subOptions.get(), sequence));
}

NatsStatus::Enum StanSubOptions::setManualAckMode(bool manualAck)
{
    return NatsStatus::as(stanSubOptions_SetManualAckMode(m_subOptions.get(), manualAck));
}

NatsStatus::Enum StanSubOptions::setMaxInflight(int inflight)
{
    return NatsStatus::as(stanSubOptions_SetMaxInflight(m_subOptions.get(), inflight));
}

NatsStatus::Enum StanSubOptions::setAckWait(qint64 ms)
{
    return NatsStatus::as(stanSubOptions_SetAckWait(m_subOptions.get(), ms));
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

}