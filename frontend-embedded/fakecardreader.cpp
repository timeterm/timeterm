#include "fakecardreader.h"
#include "mfrc522device.h"

#include <QTextStream>
#include <QThread>

FakeCardReader::FakeCardReader(QObject *parent)
    : CardReader(parent)
{
}

void FakeCardReader::start()
{
    while (!m_shutDown) {
        emit cardRead("test");

        QThread::sleep(1);
    }
}

void FakeCardReader::shutDown()
{
    m_shutDown = true;
}
