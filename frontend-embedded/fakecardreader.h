#ifndef FAKECARDREADER_H
#define FAKECARDREADER_H

#include "cardreader.h"

class FakeCardReader: public CardReader
{
    Q_OBJECT

public:
    explicit FakeCardReader(QObject *parent = nullptr);
    ~FakeCardReader() override = default;

public slots:
    void start() override;
    void shutDown() override;

private:
    bool m_shutDown = false;
};

#endif// FAKECARDREADER_H
