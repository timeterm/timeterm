#include "stanmessage.h"

namespace MessageQueue
{

StanMessage::StanMessage(QString channel, stanMsg *message)
    : m_channel(std::move(channel))
    , m_data(stanMsg_GetData(message), stanMsg_GetDataLength(message))
{
}

QString StanMessage::channel() const
{
    return m_channel;
}

QByteArray const &StanMessage::data() const
{
    return m_data;
}

} // namespace MessageQueue