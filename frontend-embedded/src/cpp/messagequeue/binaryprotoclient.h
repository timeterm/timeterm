#ifndef BINARYPROTOCLIENT_H
#define BINARYPROTOCLIENT_H

#include "stanmessage.h"
#include <QObject>
#include <src/cpp/messagequeue/messages/disowntokenmessage.h>
#include <src/cpp/messagequeue/messages/retrievenewtokenmessage.h>
#include <timeterm_proto/messages.pb.h>

namespace MessageQueue
{

class BinaryProtoClient: public QObject
{
    Q_OBJECT

public:
    explicit BinaryProtoClient(QObject *parent = nullptr);

private:
    void handleDisownTokenProto(const timeterm_proto::messages::DisownTokenMessage &msg);
    void handleRetrieveNewTokenProto(const timeterm_proto::messages::RetrieveNewTokenMessage &msg);

public slots:
    void handleMessage(const MessageQueue::StanMessage &message);

signals:
    void disownToken(MessageQueue::DisownTokenMessage msg);
    void retrieveNewToken(const MessageQueue::RetrieveNewTokenMessage &msg);
};

} // namespace MessageQueue

#endif // BINARYPROTOCLIENT_H
