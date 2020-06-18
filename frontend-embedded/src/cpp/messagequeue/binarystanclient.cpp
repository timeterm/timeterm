#include "binarystanclient.h"

namespace MessageQueue
{

void BinaryStanClient::handleMessage(const MessageQueue::StanMessage &message)
{
    BinaryClient::handleMessage(message.channel(), message.data());
}

}
