#include "stansubscription.h"
#include "stancallbackhandlersingleton.h"

namespace MessageQueue
{

StanSubscription::StanSubscription(QObject *parent)
    : QObject(parent)
{
}

void StanSubscription::setSubscription(stanSubscription *sub)
{
    m_stanSub.reset(sub);

    StanCallbackHandlerSingleton::singleton().setMsgHandler(sub, [this](const char *channel, stanMsg *msg) {
        emit messageReceived(StanMessage(QString::fromUtf8(channel), msg));
    });
}

} // namespace MessageQueue
