#include "fakecardreader.h"
#include "mfrc522cardreader.h"

#include <QProcessEnvironment>
#include <QTextStream>
#include <QThread>

FakeCardReader::FakeCardReader(QObject *parent)
    : CardReader(parent)
{
}

void FakeCardReader::start()
{
    while (!m_shutDown) {
        QString uid = QProcessEnvironment::systemEnvironment()
                          .value("FAKE_CARD_READER_EMIT_UID", "");

        if (!uid.isEmpty())
            emit cardRead(uid);

        QThread::sleep(1);
    }
}

void FakeCardReader::shutDown()
{
    m_shutDown = true;
}
