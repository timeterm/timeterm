#ifndef CARDREADER_H
#define CARDREADER_H

#include <QObject>
#include <QString>

class CardReader: public QObject
{
    Q_OBJECT

public:
    explicit CardReader(QObject *parent = nullptr);
    ~CardReader() override = default;

public slots:
    virtual void start() = 0;
    virtual void shutDown() = 0;

signals:
    void cardRead(const QString &uid);
};

#endif// CARDREADER_H
