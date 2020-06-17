#ifndef STANCALLBACKHANDLERSINGLETON_H
#define STANCALLBACKHANDLERSINGLETON_H

#include <functional>
#include <nats.h>

#include <QHash>

using StanMsgHandler = std::function<void(stanSubscription *sub, const char *channel, stanMsg *msg)>;
using StanConnLostHandler = std::function<void(const char *errTxt)>;

class StanCallbackHandlerSingleton
{
public:
    static StanCallbackHandlerSingleton &singleton();

    StanCallbackHandlerSingleton(StanCallbackHandlerSingleton const &) = delete;
    void operator=(StanCallbackHandlerSingleton const &) = delete;

    void setMsgHandler(stanConnection *conn, StanMsgHandler handler);
    void removeMsgHandler(stanConnection *conn);

    void setConnectionLostHandler(stanConnection *conn, StanConnLostHandler handler);
    void removeConnectionLostHandler(stanConnection *conn);

    static void onMsg(stanConnection *sc, stanSubscription *sub, const char *channel, stanMsg *msg, void *closure);
    static void onConnLost(stanConnection *sc, const char *errTxt, void *closure);

private:
    StanCallbackHandlerSingleton() = default;

    void onMsg(stanConnection *sc, stanSubscription *sub, const char *channel, stanMsg *msg);
    void onConnLost(stanConnection *sc, const char *errTxt);

    QHash<stanConnection *, StanMsgHandler> m_msgHandlers;
    QHash<stanConnection *, StanConnLostHandler> m_connLostHandlers;
};

#endif//STANCALLBACKHANDLERSINGLETON_H
