#pragma once

#include <QObject>
#include "zermeloappointment.h"

class ZermeloAppointments
{
    Q_GADGET
    Q_PROPERTY(QList<ZermeloAppointment> data READ data)

public:
    void append(const ZermeloAppointment &appointment);

    QList<ZermeloAppointment> data();

    void read(const QJsonObject &json);
    void write(QJsonObject &json) const;

private:
    void append(const QList<ZermeloAppointment> &appointments);

    QList<ZermeloAppointment> m_data;
};

Q_DECLARE_METATYPE(ZermeloAppointments)
