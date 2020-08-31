#include "stansubscription.h"
#include "stancallbackhandlersingleton.h"

#include <QDebug>

namespace MessageQueue
{

StanSubscription::StanSubscription(QObject *parent)
    : QObject(parent)
{
}

void StanSubscription::setSubscription(stanSubscription *sub)
{
    m_sub = sub;

    StanCallbackHandlerSingleton::singleton().setMsgHandler(sub, [this](const char *channel, stanMsg *msg) {
        qDebug() << "StanSubscription#setSubscription/msgHandler: emitting messageReceived for message on channel" << channel;
        emit messageReceived(StanMessage(QString::fromUtf8(channel), msg));
        qDebug() << "StanSubscription#setSubscription/msgHandler: emitted messageReceived for message on channel" << channel;
    });
}

StanSubscription::~StanSubscription()
{
    if (m_sub != nullptr) {
        stanSubscription_Destroy(m_sub);
    }
}

} // namespace MessageQueue
