#pragma once

#include "cardreader.h"
#include "fakecardreader.h"

#include <QObject>
#include <QThread>

#ifdef RASPBERRYPI
#include "mfrc522cardreader.h"
#endif

class CardReaderController: public QObject
{
    Q_OBJECT
    QThread cardReaderThread;

public:
    static CardReader *defaultCardReader(QObject *parent = nullptr)
    {
#ifdef RASPBERRYPI
        qDebug() << "Running on Raspberry Pi, using Mfrc522CardReader";
        return new Mfrc522CardReader(parent);
#else
        qDebug() << "Not running on a supported embedded device, using FakeCardReader";
        return new FakeCardReader(parent);
#endif
    }

    explicit CardReaderController(CardReader *cardReader = defaultCardReader(),
                                  QObject *parent = nullptr);
    ~CardReaderController() override;

signals:
    void cardRead(const QString &uid);
    void runCardReader(QPrivateSignal);

private:
    CardReader *m_cardReader;
};
