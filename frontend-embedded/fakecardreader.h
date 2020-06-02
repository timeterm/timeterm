#ifndef FAKECARDREADER_H
#define FAKECARDREADER_H

#include "cardreader.h"
#include <QtNetwork/QLocalServer>

class FakeCardReader: public CardReader
{
    Q_OBJECT

public:
    explicit FakeCardReader(QObject *parent = nullptr);
    ~FakeCardReader() override = default;

public slots:
    void start() override;
    void shutDown() override;
    void handleConnection();

signals:
    void shutDownInternal();

private:
    QLocalServer *m_server;
};

#endif// FAKECARDREADER_H
