#include "binaryprotoclient.h"
#include <src/cpp/messagequeue/messages/disowntokenmessage.h>
#include <src/cpp/messagequeue/messages/retrievenewtokenmessage.h>

#include <QDebug>

namespace MessageQueue
{

BinaryProtoClient::BinaryProtoClient(QObject *parent)
    : QObject(parent)
{
}

void BinaryProtoClient::handleMessage(const MessageQueue::StanMessage &rawMsg)
{
    if (rawMsg.channel() == "timeterm.disown-token") {
        timeterm_proto::messages::DisownTokenMessage msg;

        if (msg.ParseFromArray(rawMsg.data().data(), rawMsg.data().length()))
            handleDisownTokenProto(msg);
    } else if (rawMsg.channel() == "timeterm.retrieve-new-token") {
        timeterm_proto::messages::RetrieveNewTokenMessage msg;

        if (msg.ParseFromArray(rawMsg.data().data(), rawMsg.data().length()))
            handleRetrieveNewTokenProto(msg);
    }
}

void BinaryProtoClient::handleDisownTokenProto(const timeterm_proto::messages::DisownTokenMessage &msg)
{
    DisownTokenMessage disownTokenMessage;

    disownTokenMessage.setDeviceId(QString::fromStdString(msg.device_id()));
    disownTokenMessage.setTokenHash(QString::fromStdString(msg.token_hash()));
    disownTokenMessage.setTokenHashAlg(QString::fromStdString(msg.token_hash_alg()));

    emit disownToken(disownTokenMessage);
}

void BinaryProtoClient::handleRetrieveNewTokenProto(const timeterm_proto::messages::RetrieveNewTokenMessage &msg)
{
    RetrieveNewTokenMessage retrieveNewTokenMessage;

    retrieveNewTokenMessage.setDeviceId(QString::fromStdString(msg.device_id()));
    retrieveNewTokenMessage.setCurrentTokenHash(QString::fromStdString(msg.current_token_hash()));
    retrieveNewTokenMessage.setCurrentTokenHashAlg(QString::fromStdString(msg.current_token_hash_alg()));

    emit retrieveNewToken(retrieveNewTokenMessage);
}

} // namespace MessageQueue
