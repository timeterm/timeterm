#include "natscallbackhandlersingleton.h"
#include "enums.h"
#include "util/scopeguard.h"

#include <QDebug>

namespace MessageQueue
{

const char ACK_ACK[] = "+ACK";
const char ACK_NAK[] = "-NAK";
const char ACK_PROGRESS[] = "+WPI";
const char ACK_NEXT[] = "+NXT";
const char ACK_TERM[] = "+TERM";

NatsCallbackHandlerSingleton &NatsCallbackHandlerSingleton::singleton()
{
    static NatsCallbackHandlerSingleton instance;
    return instance;
}

void NatsCallbackHandlerSingleton::setMsgHandler(natsSubscription *sub, NatsMsgHandler handler)
{
    m_msgHandlers[sub] = std::move(handler);
}

void NatsCallbackHandlerSingleton::removeMsgHandler(natsSubscription *sub)
{
    m_msgHandlers.remove(sub);
}

void NatsCallbackHandlerSingleton::onMsg(natsConnection *nc, natsSubscription *sub, natsMsg *msg, void *closure)
{
    qDebug() << "Got message on topic" << natsMsg_GetSubject(msg);
    NatsCallbackHandlerSingleton::singleton().onMsg(nc, sub, msg);
}

enum class AckMode
{
    Ack,
    Nak,
    Progress,
    Next,
    Term,
};

NatsStatus::Enum jsAck(natsConnection *nc, natsMsg *msg, AckMode mode)
{
    const char *ackMode = ACK_ACK;

    switch (mode) {
    case AckMode::Ack:
        break; // the default
    case AckMode::Nak:
        ackMode = ACK_NAK;
        break;
    case AckMode::Progress:
        ackMode = ACK_PROGRESS;
        break;
    case AckMode::Next:
        ackMode = ACK_NEXT;
        break;
    case AckMode::Term:
        ackMode = ACK_TERM;
        break;
    }

    natsMsg *reply = nullptr;
    auto status = natsConnection_RequestString(&reply, nc, natsMsg_GetReply(msg), ackMode, 1000);
    natsMsg_Destroy(reply);
    return NatsStatus::fromC(status);
}

void NatsCallbackHandlerSingleton::onMsg(natsConnection *nc, natsSubscription *sub, natsMsg *msg)
{
    auto cleanUp = [msg]() {
        qDebug() << "Destroying message";
        natsMsg_Destroy(msg);
        qDebug() << "Destroyed message";
    };

    auto realOnMsg = [&]() {
        if (!m_msgHandlers.contains(sub)) {
            qInfo() << "No registered handlers for message on topic" << natsMsg_GetSubject(msg);

            if (natsMsg_GetReply(msg) != nullptr && strlen(natsMsg_GetReply(msg)) > 0) {
                qDebug() << "Acknowledging message";
                jsAck(nc, msg, AckMode::Nak);
            }

            return;
        }

        qDebug() << "Found message handler for subscription";
        m_msgHandlers[sub](msg);

        if (natsMsg_GetReply(msg) != nullptr && strlen(natsMsg_GetReply(msg)) > 0) {
            qDebug() << "Acknowledging message";
            jsAck(nc, msg, AckMode::Ack);
        }
    };

    after(realOnMsg, cleanUp);
}

} // namespace MessageQueue