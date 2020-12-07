#pragma once

#include <QObject>

#ifdef TIMETERMOS
#include "ttsystemd.h"
#endif

class Systemd: public QObject
{
    Q_OBJECT

public:
    explicit Systemd(QObject *parent = nullptr);

public slots:
    void rebootDevice();
};
