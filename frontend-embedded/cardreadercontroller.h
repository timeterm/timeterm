#ifndef CARDREADERCONTROLLER_H
#define CARDREADERCONTROLLER_H

#include "cardreader.h"
#include <QObject>
#include <QThread>

class CardReaderController: public QObject
{
    Q_OBJECT
    QThread cardReaderThread;

public:
    explicit CardReaderController(CardReader *cardReader, QObject *parent = nullptr);
    ~CardReaderController() override;

signals:
    void cardRead(const QString &uid);
    void runCardReader(QPrivateSignal);

private:
    CardReader *m_cardReader;
};

#endif// CARDREADERCONTROLLER_H
