#ifndef MESSAGEQUEUECLIENT_H
#define MESSAGEQUEUECLIENT_H

#include "messages/retrievenewtokenmessage.h"
#include "messages/disowntokenmessage.h"

#include <QObject>

namespace MessageQueue
{

class MessageQueueClient: public QObject
{
    Q_OBJECT

public:
    explicit MessageQueueClient(QObject *parent = nullptr);
    ~MessageQueueClient() override = default;

signals:
    void disownToken(const DisownTokenMessage &msg);
    void retrieveNewToken(const RetrieveNewTokenMessage &msg);
};

}

#endif // MESSAGEQUEUECLIENT_H
