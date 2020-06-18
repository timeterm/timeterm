#ifndef BINARYPROTOCLIENT_H
#define BINARYPROTOCLIENT_H

#include <QObject>
#include <timeterm_proto/messages.pb.h>

namespace MessageQueue
{

class BinaryProtoClient: QObject
{
    Q_OBJECT

public:
    explicit BinaryProtoClient(QObject *parent = nullptr);

public slots:
    void handleMessage(const QString &channel, const QByteArray &data);

signals:
    void disownTokenProto(const timeterm_proto::messages::DisownTokenMessage &msg);
    void retrieveNewTokenProto(const timeterm_proto::messages::RetrieveNewTokenMessage &msg);
};

}

#endif // BINARYPROTOCLIENT_H
