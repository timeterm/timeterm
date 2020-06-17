#include "stanmessage.h"

namespace MessageQueue
{

StanMessage::StanMessage(QString channel, stanMsg *message)
    : m_channel(std::move(channel))
    , m_stanMsg(message, StanMessage::deleter)
{
}

void StanMessage::deleter(stanMsg *message)
{
    stanMsg_Destroy(message);
}

QString StanMessage::channel() const
{
    return m_channel;
}

} // namespace MessageQueue