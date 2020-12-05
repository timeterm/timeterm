#include "zermeloappointments.h"

#include <QJsonArray>
#include <QJsonObject>

void ZermeloAppointments::append(const ZermeloAppointment &appointment)
{
    m_data.append(appointment);
}

void ZermeloAppointments::append(const QList<ZermeloAppointment> &appointments)
{
    m_data.append(appointments);
}

QList<ZermeloAppointment> ZermeloAppointments::data()
{
    return m_data;
}

void ZermeloAppointments::read(const QJsonObject &json)
{
    if (json.contains("data") && json["data"].isArray()) {
        for (const auto &item : json["data"].toArray()) {
            if (!item.isObject())
                continue;

            auto appointment = ZermeloAppointment();
            appointment.read(item.toObject());
            m_data.append(appointment);
        }
    }
}

void ZermeloAppointments::write(QJsonObject &json) const
{
    QJsonArray arr;
    for (const auto &appointment : m_data) {
        auto obj = QJsonObject();
        appointment.write(obj);
        arr.append(obj);
    }
    json["data"] = arr;
}
