#include "stancallbackhandlersingleton.h"

StanCallbackHandlerSingleton &StanCallbackHandlerSingleton::singleton()
{
    static StanCallbackHandlerSingleton instance;
    return instance;
}

void StanCallbackHandlerSingleton::setMsgHandler(stanConnection *conn, StanMsgHandler handler)
{
    m_msgHandlers[conn] = std::move(handler);
}

void StanCallbackHandlerSingleton::removeMsgHandler(stanConnection *conn)
{
    m_msgHandlers.remove(conn);
}

void StanCallbackHandlerSingleton::setConnectionLostHandler(stanConnection *conn, StanConnLostHandler handler)
{
    m_connLostHandlers[conn] = std::move(handler);
}

void StanCallbackHandlerSingleton::removeConnectionLostHandler(stanConnection *conn)
{
    m_connLostHandlers.remove(conn);
}

void StanCallbackHandlerSingleton::onMsg(stanConnection *sc, stanSubscription *sub, const char *channel, stanMsg *msg, void *closure)
{
    StanCallbackHandlerSingleton::singleton().onMsg(sc, sub, channel, msg);
}

void StanCallbackHandlerSingleton::onConnLost(stanConnection *sc, const char *errTxt, void *closure)
{
    StanCallbackHandlerSingleton::singleton().onConnLost(sc, errTxt);
}

void StanCallbackHandlerSingleton::onMsg(stanConnection *sc, stanSubscription *sub, const char *channel, stanMsg *msg)
{
    if (!m_msgHandlers.contains(sc))
        return;
    m_msgHandlers[sc](sub, channel, msg);
}

void StanCallbackHandlerSingleton::onConnLost(stanConnection *sc, const char *errTxt)
{
    if (!m_connLostHandlers.contains(sc))
        return;
    m_connLostHandlers[sc](errTxt);
}