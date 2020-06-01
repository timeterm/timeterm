#ifndef MFRC522DEVICE_H
#define MFRC522DEVICE_H

#include <QObject>
#include <QThread>
#include <mfrc522/mfrc522.h>
#include "cardreader.h"

class Mfrc522Device : public CardReader
{
    Q_OBJECT

public:
    explicit Mfrc522Device(QObject *parent = nullptr);
    ~Mfrc522Device() = default;

    QString makeUidString(Mfrc522::Device::Uid uid);

public slots:
    void start();
    void shutDown();

private:
    Mfrc522::Device m_mfrcDev;
    bool m_shutDown = false;
};

#endif // MFRC522DEVICE_H
