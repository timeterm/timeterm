#ifndef STANMESSAGE_H
#define STANMESSAGE_H

#include <nats.h>

#include <QObject>
#include <QSharedPointer>

class StanMessage
{
    Q_GADGET

public:
    explicit StanMessage(stanMsg *message) : m_stanMsg(message, StanMessage::deleter)
    {
    }

private:
    static void deleter(stanMsg *message) {
        stanMsg_Destroy(message);
    }

    QSharedPointer<stanMsg> m_stanMsg;
};

#endif//STANMESSAGE_H
