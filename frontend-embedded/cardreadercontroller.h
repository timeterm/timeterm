#ifndef CARDREADERCONTROLLER_H
#define CARDREADERCONTROLLER_H

#include <QObject>
#include <QThread>
#include "cardreader.h"

class CardReaderController : public QObject
{
    Q_OBJECT
    QThread cardReaderThread;
public:
    explicit CardReaderController(CardReader *cardReader, QObject *parent = nullptr);
    ~CardReaderController();

signals:
    void cardRead(const QString &uid);

private:
    CardReader *m_cardReader;
};

#endif // CARDREADERCONTROLLER_H
