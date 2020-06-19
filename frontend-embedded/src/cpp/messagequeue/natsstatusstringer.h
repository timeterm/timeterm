#ifndef NATSSTATUSSTRINGER_H
#define NATSSTATUSSTRINGER_H

#include "enums.h"
#include <QObject>

namespace MessageQueue
{

class NatsStatusStringer: public QObject
{
    Q_OBJECT

public:
    explicit NatsStatusStringer(QObject *parent = nullptr);

    Q_INVOKABLE QString stringify(int status);
};

} // namespace MessageQueue

#endif // NATSSTATUSSTRINGER_H
