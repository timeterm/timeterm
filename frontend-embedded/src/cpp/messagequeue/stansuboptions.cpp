#include "stansuboptions.h"
#include "enums.h"
#include "strings.h"

namespace MessageQueue
{

NatsStatus newStanSubOptions(QSharedPointer<stanSubOptions> &ptr)
{
    stanSubOptions *stanSubOpts = nullptr;
    auto s = static_cast<NatsStatus>(stanSubOptions_Create(&stanSubOpts));
    if (s == NatsStatus::Ok)
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

NatsStatus StanSubOptions::setDurableName(const QString &durableName)
{
    auto durableNameCstr = asUtf8CString(durableName);

    return asNatsStatus(stanSubOptions_SetDurableName(m_subOptions.get(), durableNameCstr.get()));
}

NatsStatus StanSubOptions::deliverAllAvailable()
{
    return asNatsStatus(stanSubOptions_DeliverAllAvailable(m_subOptions.get()));
}

NatsStatus StanSubOptions::startWithLastReceived()
{
    return asNatsStatus(stanSubOptions_StartWithLastReceived(m_subOptions.get()));
}

NatsStatus StanSubOptions::startAtSequence(quint64 sequence)
{
    return asNatsStatus(stanSubOptions_StartAtSequence(m_subOptions.get(), sequence));
}

NatsStatus StanSubOptions::setManualAckMode(bool manualAck)
{
    return asNatsStatus(stanSubOptions_SetManualAckMode(m_subOptions.get(), manualAck));
}

NatsStatus StanSubOptions::setMaxInflight(int inflight)
{
    return asNatsStatus(stanSubOptions_SetMaxInflight(m_subOptions.get(), inflight));
}

NatsStatus StanSubOptions::setAckWait(qint64 ms)
{
    return asNatsStatus(stanSubOptions_SetAckWait(m_subOptions.get(), ms));
}

NatsStatus StanSubOptions::lastStatus() const
{
    return m_lastStatus;
}

void StanSubOptions::updateStatus(NatsStatus s)
{
    m_lastStatus = s;
    if (s == NatsStatus::Ok)
        return;

    const char *text = natsStatus_GetText(asCNatsStatus(s));
    auto statusStr = QString::fromLocal8Bit(text);
    emit errorOccurred(s, statusStr);
}

}