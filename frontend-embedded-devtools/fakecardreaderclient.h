#pragma once

#include <QObject>
#include <QtNetwork>

class FakeCardReaderClient: public QObject
{
    Q_OBJECT
public:
    explicit FakeCardReaderClient(QObject *parent = nullptr);

public slots:
    void sendCardUid(const QString &serverName, const QString &uid);

private:
    QLocalSocket *m_sock;
};
