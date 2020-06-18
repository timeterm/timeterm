#ifndef PROTOCLIENT_H
#define PROTOCLIENT_H

#include "messagequeueclient.h"

#include <timeterm_proto/messages.pb.h>

namespace MessageQueue
{

class ProtoClient: public MessageQueueClient
{
    Q_OBJECT

public slots:
    void disownTokenProto(const timeterm_proto::messages::DisownTokenMessage &msg);
    void retrieveNewTokenProto(const timeterm_proto::messages::RetrieveNewTokenMessage &msg);
};

}

#endif // PROTOCLIENT_H
