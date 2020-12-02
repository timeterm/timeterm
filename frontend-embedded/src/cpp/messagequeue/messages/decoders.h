#pragma once

#include <QObject>
#include <QSharedPointer>
#include <QVariant>

#include <functional>
#include <messagequeue/jetstreamconsumer.h>
#include <nats.h>

#include "disowntokenmessage.h"

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

class DisownTokenMessageDecoder: public Decoder
{
    Q_OBJECT

public:
    explicit DisownTokenMessageDecoder(QObject *parent = nullptr);

    static QVariant convertDisownTokenMessage(const timeterm_proto::mq::DisownTokenMessage &msg);
};

class RetrieveNewTokenMessageDecoder: public Decoder
{
    Q_OBJECT

public:
    explicit RetrieveNewTokenMessageDecoder(QObject *parent = nullptr);

    static QVariant convertRetrieveNewTokenMessage(const timeterm_proto::mq::RetrieveNewTokenMessage &msg);
};

} // namespace MessageQueue
