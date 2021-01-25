#ifdef TIMETERMOS

#include "mfrc522cardreader.h"

#include <QTextStream>

Mfrc522CardReader::Mfrc522CardReader(QObject *parent)
    : CardReader(parent)
{
    m_mfrcDev.pcdInit();
    m_mfrcDev.pcdSetAntennaGain(0x07u << 4u);
}

void Mfrc522CardReader::start()
{
    while (!m_shutDown) {
        if (!m_mfrcDev.piccIsNewCardPresent())
            continue;

        if (!m_mfrcDev.piccReadCardSerial())
            continue;

        // A card has been read
        auto uidString = makeUidString(m_mfrcDev.getUid());

        emit cardRead(uidString);

        QThread::sleep(1);
    }
}

QString Mfrc522CardReader::makeUidString(Mfrc522::Device::Uid uid)
{
    QString result;

    QTextStream ts(&result);
    ts << Qt::hex;

    for (uint8_t i = 0; i < uid.size; ++i) {
        if (uid.uidByte[i] < 0x10) {
            ts << '0';
        }
        ts << uid.uidByte[i];
    }

    return result;
}

void Mfrc522CardReader::shutDown()
{
    m_shutDown = true;
}

#endif // TIMETERMOS
