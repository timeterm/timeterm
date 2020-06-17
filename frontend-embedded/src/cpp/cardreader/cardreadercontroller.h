#ifndef CARDREADERCONTROLLER_H
#define CARDREADERCONTROLLER_H

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
    explicit CardReaderController(CardReader *cardReader = defaultCardReader(),
                                  QObject *parent = nullptr);
    ~CardReaderController() override;

    static CardReader *defaultCardReader()
    {
#ifdef RASPBERRYPI
        return new Mfrc522CardReader();
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

#endif // CARDREADERCONTROLLER_H
