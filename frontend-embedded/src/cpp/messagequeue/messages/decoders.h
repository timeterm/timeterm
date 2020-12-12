#pragma once

#include <QObject>
#include <QSharedPointer>
#include <QVariant>

#include <functional>
#include <nats.h>

namespace MessageQueue
{

class Decoder: public QObject
{
    Q_OBJECT

public:
    using Fn = std::function<QVariant(natsMsg *msg)>;

    explicit Decoder(QObject *parent = nullptr);

    void setFn(const Fn &fn);

    Decoder *clone();

signals:
    void newMessage(const QVariant &msg);

public slots:
    void decodeMessage(const QSharedPointer<natsMsg *> &msg);

private:
    Fn m_convert;
};

template<typename T, QVariant (*convert)(const T &)>
QVariant convertProto(natsMsg *msg)
{
    T message;
    if (message.ParseFromArray(natsMsg_GetData(msg), natsMsg_GetDataLength(msg)))
        return convert(message);
    return QVariant();
}

} // namespace MessageQueue
