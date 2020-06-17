#ifndef STANCONNECTION_H
#define STANCONNECTION_H

#include <QHash>
#include <QObject>
#include <functional>
#include <nats.h>

class StanConnection: public QObject
{
    Q_OBJECT

public:
    explicit StanConnection(QObject *parent = nullptr);

    stanConnection *m_stanConnection = nullptr;
};

#endif//STANCONNECTION_H
