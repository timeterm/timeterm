#include "zermeloappointment.h"

#include <QJsonArray>
#include <QJsonObject>

void ZermeloAppointment::setId(qint64 id)
{
    if (id != m_id)
        m_id = id;
}

qint64 ZermeloAppointment::id() const
{
    return m_id;
}

void ZermeloAppointment::setAppointmentInstance(qint64 appointmentInstance)
{
    if (appointmentInstance != m_appointmentInstance)
        m_appointmentInstance = appointmentInstance;
}

qint64 ZermeloAppointment::appointmentInstance() const
{
    return m_appointmentInstance;
}

void ZermeloAppointment::setStartTimeSlot(QString startTimeSlot)
{
    if (startTimeSlot != m_startTimeSlot)
        m_startTimeSlot = startTimeSlot;
}

QString ZermeloAppointment::startTimeSlot() const
{
    return m_startTimeSlot;
}

void ZermeloAppointment::setEndTimeSlot(QString endTimeSlot)
{
    if (endTimeSlot != m_endTimeSlot)
        m_endTimeSlot = endTimeSlot;
}

QString ZermeloAppointment::endTimeSlot() const
{
    return m_endTimeSlot;
}

void ZermeloAppointment::setCapacity(qint32 capacity)
{
    if (capacity != m_capacity)
        m_capacity = capacity;
}

qint32 ZermeloAppointment::capacity() const
{
    return m_capacity;
}

void ZermeloAppointment::setAvailableSpace(qint32 availableSpace)
{
    if (availableSpace != m_availableSpace)
        m_availableSpace = availableSpace;
}

qint32 ZermeloAppointment::availableSpace() const
{
    return m_availableSpace;
}

void ZermeloAppointment::setStartTime(const QDateTime &startTime)
{
    if (startTime != m_startTime)
        m_startTime = startTime;
}

QDateTime ZermeloAppointment::startTime() const
{
    return m_startTime;
}

void ZermeloAppointment::setEndTime(const QDateTime &endTime)
{
    if (endTime != m_endTime)
        m_endTime = endTime;
}

QDateTime ZermeloAppointment::endTime() const
{
    return m_endTime;
}

void ZermeloAppointment::setSubjects(const QStringList &subjects)
{
    if (subjects != m_subjects)
        m_subjects = subjects;
}

QStringList ZermeloAppointment::subjects() const
{
    return m_subjects;
}

void ZermeloAppointment::setGroups(const QStringList &groups)
{
    if (groups != m_groups)
        m_groups = groups;
}

QStringList ZermeloAppointment::groups() const
{
    return m_groups;
}

void ZermeloAppointment::setLocations(const QStringList &locations)
{
    if (locations != m_locations)
        m_locations = locations;
}

QStringList ZermeloAppointment::locations() const
{
    return m_locations;
}

void ZermeloAppointment::setTeachers(const QStringList &teachers)
{
    if (teachers != m_teachers)
        m_teachers = teachers;
}

QStringList ZermeloAppointment::teachers() const
{
    return m_teachers;
}

void ZermeloAppointment::setIsOnline(bool isOnline)
{
    if (isOnline != m_isOnline)
        m_isOnline = isOnline;
}

bool ZermeloAppointment::isOnline() const
{
    return m_isOnline;
}

void ZermeloAppointment::setIsOptional(bool isOptional)
{
    if (isOptional != m_isOptional)
        m_isOptional = isOptional;
}

bool ZermeloAppointment::isOptional() const
{
    return m_isOptional;
}

void ZermeloAppointment::setIsStudentEnrolled(bool isStudentEnrolled)
{
    if (isStudentEnrolled != m_isStudentEnrolled)
        m_isStudentEnrolled = isStudentEnrolled;
}

bool ZermeloAppointment::isStudentEnrolled() const
{
    return m_isStudentEnrolled;
}

void ZermeloAppointment::setIsCanceled(bool isCanceled)
{
    if (isCanceled != m_isCanceled) {
        m_isCanceled = isCanceled;
    }
}

bool ZermeloAppointment::isCanceled() const
{
    return m_isCanceled;
}

void readStringArray(const QJsonArray &array, QStringList &into)
{
    for (const auto &item : array) {
        if (item.isString()) {
            into.append(item.toString());
        }
    }
}

void ZermeloAppointment::read(const QJsonObject &json)
{
    if (json.contains("id") && json["id"].isDouble())
        m_id = json["id"].toInt();

    if (json.contains("appointmentInstance") && json["appointmentInstance"].isDouble())
        m_appointmentInstance = json["appointmentInstance"].toInt();

    if (json.contains("startTimeSlot") && json["startTimeSlot"].isString())
        m_startTimeSlot = json["startTimeSlot"].toString();

    if (json.contains("endTimeSlot") && json["endTimeSlot"].isString())
        m_endTimeSlot = json["endTimeSlot"].toString();

    if (json.contains("capacity") && json["capacity"].isDouble())
        m_capacity = json["capacity"].toInt();

    if (json.contains("availableSpace") && json["availableSpace"].isDouble())
        m_availableSpace = json["availableSpace"].toInt();

    if (json.contains("startTime") && json["startTime"].isString())
        m_startTime = QDateTime::fromString(json["startTime"].toString());

    if (json.contains("endTime") && json["endTime"].isString())
        m_endTime = QDateTime::fromString(json["endTime"].toString());

    if (json.contains("subjects") && json["subjects"].isArray())
        readStringArray(json["subjects"].toArray(), m_subjects);

    if (json.contains("groups") && json["groups"].isArray())
        readStringArray(json["groups"].toArray(), m_groups);

    if (json.contains("locations") && json["locations"].isArray())
        readStringArray(json["locations"].toArray(), m_locations);

    if (json.contains("teachers") && json["teachers"].isArray())
        readStringArray(json["teachers"].toArray(), m_teachers);

    if (json.contains("isOnline") && json["isOnline"].isBool())
        m_isOnline = json["isOnline"].toBool();

    if (json.contains("isOptional") && json["isOptional"].isBool())
        m_isOptional = json["isOptional"].toBool();

    if (json.contains("isStudentEnrolled") && json["isStudentEnrolled"].isBool())
        m_isStudentEnrolled = json["isStudentEnrolled"].toBool();

    if (json.contains("isCanceled") && json["isCanceled"].isBool())
        m_isCanceled = json["isCanceled"].toBool();
}

QJsonArray stringListAsQJsonArray(const QStringList &list)
{
    QJsonArray arr;
    for (const auto &str : list) {
        arr.append(str);
    }
    return arr;
}

void ZermeloAppointment::write(QJsonObject &json) const
{
    json["id"] = m_id;
    json["appointmentInstance"] = m_appointmentInstance;
    json["startTimeSlot"] = m_startTimeSlot;
    json["endTimeSlot"] = m_endTimeSlot;
    json["capacity"] = m_capacity;
    json["availableSpace"] = m_availableSpace;
    json["startTime"] = m_startTime.toString();
    json["endTime"] = m_endTime.toString();
    json["subjects"] = stringListAsQJsonArray(m_subjects);
    json["groups"] = stringListAsQJsonArray(m_groups);
    json["locations"] = stringListAsQJsonArray(m_locations);
    json["teachers"] = stringListAsQJsonArray(m_teachers);
    json["isOnline"] = m_isOnline;
    json["isStudentEnrolled"] = m_isStudentEnrolled;
    json["isCanceled"] = m_isCanceled;
}
