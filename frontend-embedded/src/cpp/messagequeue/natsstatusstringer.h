#pragma once

#include "enums.h"
#include <QObject>

namespace MessageQueue
{

class NatsStatusStringer: public QObject
{
    Q_OBJECT

public:
    explicit NatsStatusStringer(QObject *parent = nullptr);

    Q_INVOKABLE static QString stringify(int status);
};

} // namespace MessageQueue
