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

}
