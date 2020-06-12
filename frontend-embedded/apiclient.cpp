#include "apiclient.h"

ApiClient::ApiClient(QObject *parent) : QObject(parent)
{
}

QString ApiClient::cardId() const
{
    return m_cardId;
}

void ApiClient::setCardId(const QString &cardId)
{
    if (cardId != m_cardId) {
        m_cardId = cardId;
        emit cardIdChanged();
    }
}

Appointments::Appointments(QObject *parent) : QObject(parent)
{
}

void Appointments::append(Appointment *appointment)
{
    m_appointments.append(appointment);
    emit appointmentsChanged();
}

void Appointments::append(const QList<Appointment *> &appointments)
{
    m_appointments.append(appointments);
    emit appointmentsChanged();
}

QQmlListProperty<Appointment> Appointments::appointments()
{
    return {this, &m_appointments};
}

Appointment::Appointment(QObject *parent) : QObject(parent)
{
}

void Appointment::setId(qint64 id)
{
    if (id != m_id) {
        m_id = id;
        emit idChanged();
    }
}

qint64 Appointment::id() const
{
    return m_id;
}

void Appointment::setAppointmentInstance(qint64 appointmentInstance)
{
    if (appointmentInstance != m_appointmentInstance) {
        m_appointmentInstance = appointmentInstance;
        emit appointmentInstanceChanged();
    }
}

qint64 Appointment::appointmentInstance() const
{
    return m_appointmentInstance;
}

void Appointment::setStartTimeSlot(qint32 startTimeSlot)
{
    if (startTimeSlot != m_startTimeSlot) {
        m_startTimeSlot = startTimeSlot;
        emit startTimeSlotChanged();
    }
}

qint32 Appointment::startTimeSlot() const
{
    return m_startTimeSlot;
}

void Appointment::setEndTimeSlot(qint32 endTimeSlot)
{
    if (endTimeSlot != m_endTimeSlot) {
        m_endTimeSlot = endTimeSlot;
        emit endTimeSlotChanged();
    }
}

qint32 Appointment::endTimeSlot() const
{
    return m_endTimeSlot;
}

void Appointment::setCapacity(qint32 capacity)
{
    if (capacity != m_capacity) {
        m_capacity = capacity;
        emit capacityChanged();
    }
}

qint32 Appointment::capacity() const
{
    return m_capacity;
}

void Appointment::setAvailableSpace(qint32 availableSpace)
{
    if (availableSpace != m_availableSpace) {
        m_availableSpace = availableSpace;
        emit availableSpaceChanged();
    }
}

qint32 Appointment::availableSpace() const
{
    return m_availableSpace;
}

void Appointment::setStartTime(const QDateTime &startTime)
{
    if (startTime != m_startTime) {
        m_startTime = startTime;
        emit startTimeChanged();
    }
}

QDateTime Appointment::startTime() const
{
    return m_startTime;
}

void Appointment::setEndTime(const QDateTime &endTime)
{
    if (endTime != m_endTime) {
        m_endTime = endTime;
        emit endTimeChanged();
    }
}

QDateTime Appointment::endTime() const
{
    return m_endTime;
}

void Appointment::setSubjects(const QStringList &subjects)
{
    if (subjects != m_subjects) {
        m_subjects = subjects;
        emit subjectsChanged();
    }
}

QStringList Appointment::subjects() const
{
    return m_subjects;
}

void Appointment::setLocations(const QStringList &locations)
{
    if (locations != m_locations) {
        m_locations = locations;
        emit locationsChanged();
    }
}

QStringList Appointment::locations() const
{
    return m_locations;
}

void Appointment::setTeachers(const QStringList &teachers)
{
    if (teachers != m_teachers) {
        m_teachers = teachers;
        emit teachersChanged();
    }
}

QStringList Appointment::teachers() const
{
    return m_teachers;
}

void Appointment::setIsOnline(bool isOnline)
{
    if (isOnline != m_isOnline) {
        m_isOnline = isOnline;
        emit isOnlineChanged();
    }
}

bool Appointment::isOnline() const
{
    return m_isOnline;
}

void Appointment::setIsOptional(bool isOptional)
{
    if (isOptional != m_isOptional) {
        m_isOptional = isOptional;
        emit isOptionalChanged();
    }
}

bool Appointment::isOptional() const
{
    return m_isOptional;
}

void Appointment::setIsStudentEnrolled(bool isStudentEnrolled)
{
    if (isStudentEnrolled != m_isStudentEnrolled) {
        m_isStudentEnrolled = isStudentEnrolled;
        emit isStudentEnrolledChanged();
    }
}

bool Appointment::isStudentEnrolled() const
{
    return m_isStudentEnrolled;
}

void Appointment::setIsCanceled(bool isCanceled)
{
    if (isCanceled != m_isCanceled) {
        m_isCanceled = isCanceled;
        emit isCanceledChanged();
    }
}

bool Appointment::isCanceled() const
{
    return m_isCanceled;
}
