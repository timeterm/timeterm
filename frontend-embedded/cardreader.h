#ifndef CARDREADER_H
#define CARDREADER_H

#include <QObject>
#include <QString>

class CardReader : public QObject
{
    Q_OBJECT

public:
    explicit CardReader(QObject *parent = nullptr);
    virtual ~CardReader() = default;

public slots:
    virtual void start() = 0;
    virtual void shutDown() = 0;

signals:
    virtual void cardRead(const QString &uid) = 0;
};

#endif // CARDREADER_H
