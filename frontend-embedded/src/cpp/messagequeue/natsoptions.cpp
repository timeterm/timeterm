#include "natsoptions.h"

namespace MessageQueue
{

MessageQueue::NatsOptions::NatsOptions(QObject *parent)
    : QObject(parent)
    , m_options(nullptr, natsOptions_Destroy)
{
    natsOptions *options;
    auto s = natsOptions_Create(&options);
    if (s == NATS_OK)
        m_options.reset(options);

    updateStatus(NatsStatus::fromC(s));

    if (s != NATS_OK)
        natsOptions_Destroy(options);
}

QSharedPointer<natsOptions> NatsOptions::options()
{
    return m_options;
}

NatsStatus::Enum NatsOptions::lastStatus() const
{
    return m_lastStatus;
}

void NatsOptions::updateStatus(NatsStatus::Enum s)
{
    m_lastStatus = s;
    if (s == NatsStatus::Enum::Ok)
        return;

    const char *text = natsStatus_GetText(NatsStatus::asC(s));
    auto statusStr = QString::fromLocal8Bit(text);
    emit errorOccurred(s, statusStr);
}

} // namespace MessageQueue
