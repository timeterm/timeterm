#include "stancallbackhandlersingleton.h"

#include <QDebug>

namespace MessageQueue
{

StanCallbackHandlerSingleton &StanCallbackHandlerSingleton::singleton()
{
    static StanCallbackHandlerSingleton instance;
    return instance;
}

void StanCallbackHandlerSingleton::setMsgHandler(stanSubscription *sub, StanMsgHandler handler)
{
    m_msgHandlers[sub] = std::move(handler);
}

void StanCallbackHandlerSingleton::removeMsgHandler(stanSubscription *sub)
{
    m_msgHandlers.remove(sub);
}

void StanCallbackHandlerSingleton::setConnectionLostHandler(stanConnection *conn, StanConnLostHandler handler)
{
    m_connLostHandlers[conn] = std::move(handler);
}

void StanCallbackHandlerSingleton::removeConnectionLostHandler(stanConnection *conn)
{
    m_connLostHandlers.remove(conn);
}

void StanCallbackHandlerSingleton::onMsg(stanConnection *sc, stanSubscription *sub, const char *channel, stanMsg *msg, void */* closure */)
{
    qDebug() << "StanCallbackHandlerSingleton: got message on channel" << channel;
    StanCallbackHandlerSingleton::singleton().onMsg(sc, sub, channel, msg);
}

void StanCallbackHandlerSingleton::onConnLost(stanConnection *sc, const char *errTxt, void */* closure */)
{
    StanCallbackHandlerSingleton::singleton().onConnLost(sc, errTxt);
}

void StanCallbackHandlerSingleton::onMsg(stanConnection *, stanSubscription *sub, const char *channel, stanMsg *msg)
{
    if (!m_msgHandlers.contains(sub))
        return;
    qDebug() << "StanCallbackHandlerSingleton: found message handler for subscription";
    m_msgHandlers[sub](channel, msg);
    qDebug() << "StanCallbackHandlerSingleton: destroying message";
    stanMsg_Destroy(msg);
    qDebug() << "StanCallbackHandlerSingleton: destroyed message";
}

void StanCallbackHandlerSingleton::onConnLost(stanConnection *sc, const char *errTxt)
{
    if (!m_connLostHandlers.contains(sc))
        return;
    m_connLostHandlers[sc](errTxt);
}

} // namespace MessageQueue
