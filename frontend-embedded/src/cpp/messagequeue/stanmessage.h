#ifndef STANMESSAGE_H
#define STANMESSAGE_H

#include <nats.h>

#include <QObject>
#include <QSharedPointer>

#include "scopedpointer.h"

namespace MessageQueue
{

class StanMessage
{
    Q_GADGET

public:
    explicit StanMessage(stanMsg *message);

private:
    static void deleter(stanMsg *message);

    QSharedPointer<stanMsg> m_stanMsg;
};

} // namespace MessageQueue

#endif // STANMESSAGE_H
