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
        return new Mfrc522CardReader();
#else
        return new FakeCardReader();
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
