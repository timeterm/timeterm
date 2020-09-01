#include "stancallbackhandlersingleton.h"

#include <QDebug>
#include <src/cpp/util/defer.h>

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

void StanCallbackHandlerSingleton::onMsg(stanConnection *sc, stanSubscription *sub, const char *channel, stanMsg *msg, void * /* closure */)
{
    qDebug() << "Got message on channel" << channel;
    StanCallbackHandlerSingleton::singleton().onMsg(sc, sub, channel, msg);
}

void StanCallbackHandlerSingleton::onConnLost(stanConnection *sc, const char *errTxt, void * /* closure */)
{
    StanCallbackHandlerSingleton::singleton().onConnLost(sc, errTxt);
}

void StanCallbackHandlerSingleton::onMsg(stanConnection *, stanSubscription *sub, const char *channel, stanMsg *msg)
{
    auto cleanUp = [msg]() {
        qDebug() << "Destroying message";
        stanMsg_Destroy(msg);
        qDebug() << "Destroyed message";
    };

    auto realOnMsg = [&]() {
        if (!m_msgHandlers.contains(sub)) {
            qInfo() << "No registered handler for message on channel" << channel;
            return;
        }

        qDebug() << "Found message handler for subscription";
        m_msgHandlers[sub](channel, msg);
    };

    after(realOnMsg, cleanUp);
}

void StanCallbackHandlerSingleton::onConnLost(stanConnection *sc, const char *errTxt)
{
    if (!m_connLostHandlers.contains(sc))
        return;
    m_connLostHandlers[sc](errTxt);
}

} // namespace MessageQueue
