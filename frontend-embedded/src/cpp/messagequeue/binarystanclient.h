#ifndef BINARYSTANCLIENT_H
#define BINARYSTANCLIENT_H

#include "binaryclient.h"
#include "stanmessage.h"

namespace MessageQueue
{

class BinaryStanClient: public BinaryClient
{
public slots:
    void handleMessage(const StanMessage &message);
};

}

#endif // BINARYSTANCLIENT_H
