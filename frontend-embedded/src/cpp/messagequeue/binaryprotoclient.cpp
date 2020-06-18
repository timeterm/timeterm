#include "binaryprotoclient.h"

namespace MessageQueue
{

BinaryProtoClient::BinaryProtoClient(QObject *parent)
    : QObject(parent)
{
}

void BinaryProtoClient::handleMessage(const QString &channel, const QByteArray &data)
{
    if (channel == "timeterm.disown-token") {
        timeterm_proto::messages::DisownTokenMessage msg;

        if (msg.ParseFromArray(data.data(), data.length()))
            emit disownTokenProto(msg);
    } else if (channel == "timeterm.retrieve-new-token") {
        timeterm_proto::messages::RetrieveNewTokenMessage msg;

        if (msg.ParseFromArray(data.data(), data.length()))
            emit retrieveNewTokenProto(msg);
    }
}

} // namespace MessageQueue
