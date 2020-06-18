#include "binarystanprotoclient.h"

namespace MessageQueue
{

void BinaryStanProtoClient::handleMessage(const MessageQueue::StanMessage &message)
{
    BinaryProtoClient::handleMessage(message.channel(), message.data());
}

}
