#ifdef RASPBERRYPI

#ifndef MFRC522CARDREADER_H
#define MFRC522CARDREADER_H

#include "cardreader.h"
#include <QObject>
#include <QThread>
#include <mfrc522/mfrc522.h>

class Mfrc522CardReader: public CardReader
{
    Q_OBJECT

public:
    explicit Mfrc522CardReader(QObject *parent = nullptr);
    ~Mfrc522CardReader() override = default;

    static QString makeUidString(Mfrc522::Device::Uid uid);

public slots:
    void start() override;
    void shutDown() override;

private:
    Mfrc522::Device m_mfrcDev;
    bool m_shutDown = false;
};

#endif// MFRC522CARDREADER_H

#endif// RASPBERRPI
