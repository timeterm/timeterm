#include "fakecardreader.h"
#include "mfrc522cardreader.h"

#include <QProcessEnvironment>
#include <QtNetwork>

FakeCardReader::FakeCardReader(QObject *parent)
    : CardReader(parent),
      m_server(new QLocalServer(this))
{
    connect(m_server, &QLocalServer::newConnection, this, &FakeCardReader::handleConnection);
    connect(this, &FakeCardReader::shutDownInternal, m_server, &QLocalServer::close);
}

void FakeCardReader::start()
{
    auto randomNumber = QRandomGenerator().generate();
    auto serverName = "fake_card_reader_" + QString::number(randomNumber);

    if (!m_server->listen(serverName))
        throw std::runtime_error("Could not create fake card reader socket");

    qDebug() << "Fake card reader server listening with server name " << serverName;
}

void FakeCardReader::shutDown()
{
    emit shutDownInternal();
}

void FakeCardReader::handleConnection()
{
    QLocalSocket *conn = m_server->nextPendingConnection();
    connect(conn, &QLocalSocket::disconnected, conn, &QLocalSocket::deleteLater);
    connect(this, &FakeCardReader::shutDownInternal, conn, &QLocalSocket::close);

    conn->waitForConnected();
    conn->waitForReadyRead();

    QDataStream in(conn);
    in.setVersion(QDataStream::Qt_5_14);
    quint32 blockSize = 0;

    if (conn->bytesAvailable() < (int) sizeof(quint32))
        return;
    in >> blockSize;

    if (conn->bytesAvailable() < blockSize || in.atEnd())
        return;

    QString cardUid;
    in >> cardUid;

    emit cardRead(cardUid);
    
    conn->deleteLater();
}
