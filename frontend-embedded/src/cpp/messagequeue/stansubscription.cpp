#include "stansubscription.h"

namespace MessageQueue
{

StanSubscription::StanSubscription(QObject *parent)
    : QObject(parent)
{
}

void StanSubscription::setSubscription(stanSubscription *sub)
{
    m_stanSub.reset(sub);
}

} // namespace MessageQueue