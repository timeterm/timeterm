#include "stanconnection.h"
#include "stancallbackhandlersingleton.h"

#include <utility>

StanConnection::StanConnection(QObject *parent) : QObject(parent)
{
    // TODO: actually use the error, maybe don't connect in the constructor

    natsStatus s;
    natsOptions *natsOpts = nullptr;

    s = natsOptions_Create(&natsOpts);
    if (s != NATS_OK)
        throw std::runtime_error("could not create NATS options");

    stanConnOptions *stanConnOpts = nullptr;
    s = stanConnOptions_Create(&stanConnOpts);
    if (s == NATS_OK)
        s = stanConnOptions_SetNATSOptions(stanConnOpts, natsOpts);

    if (s == NATS_OK)
        s = stanConnOptions_SetConnectionLostHandler(stanConnOpts, StanCallbackHandlerSingleton::onConnLost, nullptr);

    const char *cluster = "test-cluster";
    const char *clientId = "client";
    if (s == NATS_OK)
        s = stanConnection_Connect(&m_stanConnection, cluster, clientId, stanConnOpts);

    natsOptions_Destroy(natsOpts);
    stanConnOptions_Destroy(stanConnOpts);
}
