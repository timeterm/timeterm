#include "fakecardreaderclient.h"

FakeCardReaderClient::FakeCardReaderClient(QObject *parent)
    : QObject(parent),
      m_sock(new QLocalSocket(this))
{
}

void FakeCardReaderClient::sendCardUid(const QString &uid)
{
    QByteArray block;
    QDataStream out(&block, QIODevice::WriteOnly);
    out.setVersion(QDataStream::Qt_5_14);

    out << static_cast<quint32>(uid.length());
    out << uid;

    m_sock->abort();
    m_sock->connectToServer("fake_card_reader");
    m_sock->write(block);
    m_sock->flush();
    m_sock->disconnectFromServer();
}
