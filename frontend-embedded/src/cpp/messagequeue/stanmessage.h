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
    Q_PROPERTY(QString channel READ channel)

public:
    StanMessage() = default;
    StanMessage(QString channel, stanMsg *message);

    [[nodiscard]] QString channel() const;
    [[nodiscard]] QByteArray const &data() const;

private:
    QString m_channel;
    QByteArray m_data;
};

} // namespace MessageQueue

Q_DECLARE_METATYPE(MessageQueue::StanMessage)

#endif // STANMESSAGE_H
