#ifndef STANSUBSCRIPTION_H
#define STANSUBSCRIPTION_H

#include <QObject>

#include "stanmessage.h"

class StanSubscription: public QObject
{
    Q_OBJECT

public:
    explicit StanSubscription(QObject *parent = nullptr);

signals:
    void messageReceived(StanMessage &message);

private:
    stanSubscription *m_stanSub;
};

#endif//STANSUBSCRIPTION_H
