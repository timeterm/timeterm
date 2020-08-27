#include "protoclient.h"

namespace MessageQueue {

void ProtoClient::disownTokenProto(const timeterm_proto::messages::DisownTokenMessage &msg)
{
    DisownTokenMessage disownTokenMessage;

    disownTokenMessage.setDeviceId(QString::fromStdString(msg.device_id()));
    disownTokenMessage.setTokenHash(QString::fromStdString(msg.token_hash()));
    disownTokenMessage.setTokenHashAlg(QString::fromStdString(msg.token_hash_alg()));

    emit disownToken(disownTokenMessage);
}

void ProtoClient::retrieveNewTokenProto(const timeterm_proto::messages::RetrieveNewTokenMessage &msg)
{
    RetrieveNewTokenMessage retrieveNewTokenMessage;

    retrieveNewTokenMessage.setDeviceId(QString::fromStdString(msg.device_id()));
    retrieveNewTokenMessage.setCurrentTokenHash(QString::fromStdString(msg.current_token_hash()));
    retrieveNewTokenMessage.setCurrentTokenHashAlg(QString::fromStdString(msg.current_token_hash_alg()));

    emit retrieveNewToken(retrieveNewTokenMessage);
}

}
