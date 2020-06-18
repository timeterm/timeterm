#ifndef BINARYSTANPROTOCLIENT_H
#define BINARYSTANPROTOCLIENT_H

#include "binaryprotoclient.h"
#include "stanmessage.h"

namespace MessageQueue
{

class BinaryStanProtoClient: public BinaryProtoClient
{
public slots:
    void handleMessage(const StanMessage &message);
};

}

#endif // BINARYSTANPROTOCLIENT_H
