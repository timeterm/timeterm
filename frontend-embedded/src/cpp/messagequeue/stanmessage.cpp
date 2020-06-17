#include "stanmessage.h"

namespace MessageQueue
{

StanMessage::StanMessage(stanMsg *message)
    : m_stanMsg(message, StanMessage::deleter)
{
}

void StanMessage::deleter(stanMsg *message)
{
    stanMsg_Destroy(message);
}

} // namespace MessageQueue