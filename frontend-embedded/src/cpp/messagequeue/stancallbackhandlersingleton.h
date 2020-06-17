#ifndef STANCALLBACKHANDLERSINGLETON_H
#define STANCALLBACKHANDLERSINGLETON_H

#include <functional>
#include <nats.h>

#include <QHash>

namespace MessageQueue
{

using StanMsgHandler = std::function<void(const char *channel, stanMsg *msg)>;
using StanConnLostHandler = std::function<void(const char *errTxt)>;

class StanCallbackHandlerSingleton
{
public:
    static StanCallbackHandlerSingleton &singleton();

    StanCallbackHandlerSingleton(StanCallbackHandlerSingleton const &) = delete;
    void operator=(StanCallbackHandlerSingleton const &) = delete;

    void setMsgHandler(stanSubscription *conn, StanMsgHandler handler);
    void removeMsgHandler(stanSubscription *conn);

    void setConnectionLostHandler(stanConnection *conn, StanConnLostHandler handler);
    void removeConnectionLostHandler(stanConnection *conn);

    static void onMsg(stanConnection *sc, stanSubscription *sub, const char *channel, stanMsg *msg, void *closure);
    static void onConnLost(stanConnection *sc, const char *errTxt, void *closure);

private:
    StanCallbackHandlerSingleton() = default;

    void onMsg(stanConnection *sc, stanSubscription *sub, const char *channel, stanMsg *msg);
    void onConnLost(stanConnection *sc, const char *errTxt);

    QHash<stanSubscription *, StanMsgHandler> m_msgHandlers;
    QHash<stanConnection *, StanConnLostHandler> m_connLostHandlers;
};

} // namespace MessageQueue

#endif // STANCALLBACKHANDLERSINGLETON_H
