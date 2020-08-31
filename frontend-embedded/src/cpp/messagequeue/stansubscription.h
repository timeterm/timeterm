#ifndef STANSUBSCRIPTION_H
#define STANSUBSCRIPTION_H

#include <QObject>

#include "stanmessage.h"

namespace MessageQueue
{

using StanSubscriptionDeleter = ScopedPointerDestroyerDeleter<stanSubscription, void, stanSubscription_Destroy>;
using StanSubscriptionScopedPointer = QScopedPointer<stanSubscription, StanSubscriptionDeleter>;

class StanSubscription: public QObject
{
    Q_OBJECT

public:
    explicit StanSubscription(QObject *parent = nullptr);
    ~StanSubscription() override;

    void setSubscription(stanSubscription *sub);

signals:
    void messageReceived(const MessageQueue::StanMessage &message);

private:
    stanSubscription *m_sub = nullptr;
};

} // namespace MessageQueue

#endif // STANSUBSCRIPTION_H
