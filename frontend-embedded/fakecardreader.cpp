#include "fakecardreader.h"
#include "mfrc522cardreader.h"

#include <QProcessEnvironment>
#include <QThread>
#include <QtNetwork/QLocalServer>
#include <QtNetwork/QLocalSocket>

FakeCardReader::FakeCardReader(QObject *parent)
    : CardReader(parent),
      m_server(new QLocalServer(this))
{
    connect(m_server, &QLocalServer::newConnection, this, &FakeCardReader::handleConnection);
    connect(this, &FakeCardReader::shutDownInternal, m_server, &QLocalServer::close);
}

void FakeCardReader::start()
{
    if (!m_server->listen("fake_card_reader")) {
        throw std::runtime_error("Could not create fake card reader socket");
    }
}

void FakeCardReader::shutDown()
{
    emit shutDownInternal();
}

void FakeCardReader::handleConnection()
{
    QLocalSocket *conn = m_server->nextPendingConnection();
    connect(this, &FakeCardReader::shutDownInternal, conn, &QLocalSocket::close);

    while (conn->isOpen()) {
        conn->waitForBytesWritten();

        if (!conn->canReadLine())
            continue;

        QByteArray line = conn->readLine();
        emit cardRead(line);
    }
}
