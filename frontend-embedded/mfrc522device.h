#ifdef RASPBERRYPI

#ifndef MFRC522DEVICE_H
#define MFRC522DEVICE_H

#include "cardreader.h"
#include <QObject>
#include <QThread>
#include <mfrc522/mfrc522.h>

class Mfrc522Device: public CardReader
{
    Q_OBJECT

public:
    explicit Mfrc522Device(QObject *parent = nullptr);
    ~Mfrc522Device() override = default;

    static QString makeUidString(Mfrc522::Device::Uid uid);

public slots:
    void start() override;
    void shutDown() override;

private:
    Mfrc522::Device m_mfrcDev;
    bool m_shutDown = false;
};

#endif// MFRC522DEVICE_H

#endif// RASPBERRPI
