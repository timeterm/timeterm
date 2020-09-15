#pragma once

#include <functional>

#include <QHash>

#include <nats.h>

namespace MessageQueue
{

using NatsMsgHandler = std::function<void(natsMsg *msg)>;

class NatsCallbackHandlerSingleton
{
public:
    static NatsCallbackHandlerSingleton &singleton();

    NatsCallbackHandlerSingleton(NatsCallbackHandlerSingleton const &) = delete;
    void operator=(NatsCallbackHandlerSingleton const &) = delete;

    void setMsgHandler(natsSubscription *sub, NatsMsgHandler handler);
    void removeMsgHandler(natsSubscription *sub);

    static void onMsg(natsConnection *nc, natsSubscription *sub, natsMsg *msg, void *closure);

private:
    NatsCallbackHandlerSingleton() = default;

    void onMsg(natsConnection *nc, natsSubscription *sub, natsMsg *msg);

    QHash<natsSubscription *, NatsMsgHandler> m_msgHandlers;
};

} // namespace MessageQueue
