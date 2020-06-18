#ifndef BINARYCLIENT_H
#define BINARYCLIENT_H

#include <QObject>
#include <timeterm_proto/messages.pb.h>

namespace MessageQueue
{

class BinaryClient: QObject
{
    Q_OBJECT

public:
    explicit BinaryClient(QObject *parent = nullptr);

public slots:
    void handleMessage(const QString &channel, const QByteArray &data);

signals:
    void disownTokenProto(const timeterm_proto::messages::DisownTokenMessage &msg);
    void retrieveNewTokenProto(const timeterm_proto::messages::RetrieveNewTokenMessage &msg);
};

}

#endif // BINARYCLIENT_H
