#pragma once

#include <QObject>

#ifdef TIMTERMOS
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
