#ifndef CARDREADERCONTROLLER_H
#define CARDREADERCONTROLLER_H

#include "cardreader.h"
#include "fakecardreader.h"
#include <QObject>
#include <QThread>

#ifdef RASPBERRYPI
#include "mfrc522device.h"
#endif

class CardReaderController: public QObject
{
    Q_OBJECT
    QThread cardReaderThread;

public:
    explicit CardReaderController(CardReader *cardReader, QObject *parent = nullptr);
    ~CardReaderController() override;

    static CardReader *defaultCardReader() {
#ifdef RASPBERRYPI
        return new Mfrc522Device();
#else
        return new FakeCardReader();
#endif
    }

signals:
    void cardRead(const QString &uid);
    void runCardReader(QPrivateSignal);

private:
    CardReader *m_cardReader;
};

#endif// CARDREADERCONTROLLER_H
