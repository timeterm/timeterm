#include "decoders.h"

namespace MessageQueue
{

Decoder::Decoder(QObject *parent)
    : QObject(parent)
{
}

void Decoder::setFn(const Decoder::Fn &fn)
{
    m_convert = fn;
}

void Decoder::decodeMessage(const QSharedPointer<natsMsg *> &msg)
{
    if (m_convert) {
        QVariant result = m_convert(*msg);
        if (result.isValid() && !result.isNull()) {
            emit newMessage(result);
        }
    }
}

Decoder *Decoder::clone()
{
    auto decoder = new Decoder(this->parent());
    decoder->setFn(this->m_convert);
    return decoder;
}

DisownTokenMessageDecoder::DisownTokenMessageDecoder(QObject *parent)
    : Decoder(parent)
{
    setFn(convertProto<timeterm_proto::mq::DisownTokenMessage, convertDisownTokenMessage>);
}

QVariant DisownTokenMessageDecoder::convertDisownTokenMessage(const timeterm_proto::mq::DisownTokenMessage &msg)
{
    MessageQueue::DisownTokenMessage m;

    m.setDeviceId(QString::fromStdString(msg.device_id()));
    m.setTokenHash(QString::fromStdString(msg.token_hash()));
    m.setTokenHashAlg(QString::fromStdString(msg.token_hash_alg()));

    return QVariant::fromValue(m);
}

RetrieveNewTokenMessageDecoder::RetrieveNewTokenMessageDecoder(QObject *parent)
    : Decoder(parent)
{
    setFn(convertProto<timeterm_proto::mq::RetrieveNewTokenMessage, convertRetrieveNewTokenMessage>);
}

QVariant RetrieveNewTokenMessageDecoder::convertRetrieveNewTokenMessage(const timeterm_proto::mq::RetrieveNewTokenMessage &msg)
{
    MessageQueue::RetrieveNewTokenMessage m;

    m.setDeviceId(QString::fromStdString(msg.device_id()));
    m.setCurrentTokenHash(QString::fromStdString(msg.current_token_hash()));
    m.setCurrentTokenHashAlg(QString::fromStdString(msg.current_token_hash_alg()));

    return QVariant::fromValue(m);
}

}
