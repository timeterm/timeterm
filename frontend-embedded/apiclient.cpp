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

ZermeloAppointments::ZermeloAppointments(QObject *parent) : QObject(parent)
{
}

void ZermeloAppointments::append(ZermeloAppointment *appointment)
{
    m_data.append(appointment);
    emit dataChanged();
}

void ZermeloAppointments::append(const QList<ZermeloAppointment *> &appointments)
{
    m_data.append(appointments);
    emit dataChanged();
}

QQmlListProperty<ZermeloAppointment> ZermeloAppointments::data()
{
    return {this, &m_data};
}

ZermeloAppointment::ZermeloAppointment(QObject *parent) : QObject(parent)
{
}

void ZermeloAppointment::setId(qint64 id)
{
    if (id != m_id) {
        m_id = id;
        emit idChanged();
    }
}

qint64 ZermeloAppointment::id() const
{
    return m_id;
}

void ZermeloAppointment::setAppointmentInstance(qint64 appointmentInstance)
{
    if (appointmentInstance != m_appointmentInstance) {
        m_appointmentInstance = appointmentInstance;
        emit appointmentInstanceChanged();
    }
}

qint64 ZermeloAppointment::appointmentInstance() const
{
    return m_appointmentInstance;
}

void ZermeloAppointment::setStartTimeSlot(qint32 startTimeSlot)
{
    if (startTimeSlot != m_startTimeSlot) {
        m_startTimeSlot = startTimeSlot;
        emit startTimeSlotChanged();
    }
}

qint32 ZermeloAppointment::startTimeSlot() const
{
    return m_startTimeSlot;
}

void ZermeloAppointment::setEndTimeSlot(qint32 endTimeSlot)
{
    if (endTimeSlot != m_endTimeSlot) {
        m_endTimeSlot = endTimeSlot;
        emit endTimeSlotChanged();
    }
}

qint32 ZermeloAppointment::endTimeSlot() const
{
    return m_endTimeSlot;
}

void ZermeloAppointment::setCapacity(qint32 capacity)
{
    if (capacity != m_capacity) {
        m_capacity = capacity;
        emit capacityChanged();
    }
}

qint32 ZermeloAppointment::capacity() const
{
    return m_capacity;
}

void ZermeloAppointment::setAvailableSpace(qint32 availableSpace)
{
    if (availableSpace != m_availableSpace) {
        m_availableSpace = availableSpace;
        emit availableSpaceChanged();
    }
}

qint32 ZermeloAppointment::availableSpace() const
{
    return m_availableSpace;
}

void ZermeloAppointment::setStartTime(const QDateTime &startTime)
{
    if (startTime != m_startTime) {
        m_startTime = startTime;
        emit startTimeChanged();
    }
}

QDateTime ZermeloAppointment::startTime() const
{
    return m_startTime;
}

void ZermeloAppointment::setEndTime(const QDateTime &endTime)
{
    if (endTime != m_endTime) {
        m_endTime = endTime;
        emit endTimeChanged();
    }
}

QDateTime ZermeloAppointment::endTime() const
{
    return m_endTime;
}

void ZermeloAppointment::setSubjects(const QStringList &subjects)
{
    if (subjects != m_subjects) {
        m_subjects = subjects;
        emit subjectsChanged();
    }
}

QStringList ZermeloAppointment::subjects() const
{
    return m_subjects;
}

void ZermeloAppointment::setLocations(const QStringList &locations)
{
    if (locations != m_locations) {
        m_locations = locations;
        emit locationsChanged();
    }
}

QStringList ZermeloAppointment::locations() const
{
    return m_locations;
}

void ZermeloAppointment::setTeachers(const QStringList &teachers)
{
    if (teachers != m_teachers) {
        m_teachers = teachers;
        emit teachersChanged();
    }
}

QStringList ZermeloAppointment::teachers() const
{
    return m_teachers;
}

void ZermeloAppointment::setIsOnline(bool isOnline)
{
    if (isOnline != m_isOnline) {
        m_isOnline = isOnline;
        emit isOnlineChanged();
    }
}

bool ZermeloAppointment::isOnline() const
{
    return m_isOnline;
}

void ZermeloAppointment::setIsOptional(bool isOptional)
{
    if (isOptional != m_isOptional) {
        m_isOptional = isOptional;
        emit isOptionalChanged();
    }
}

bool ZermeloAppointment::isOptional() const
{
    return m_isOptional;
}

void ZermeloAppointment::setIsStudentEnrolled(bool isStudentEnrolled)
{
    if (isStudentEnrolled != m_isStudentEnrolled) {
        m_isStudentEnrolled = isStudentEnrolled;
        emit isStudentEnrolledChanged();
    }
}

bool ZermeloAppointment::isStudentEnrolled() const
{
    return m_isStudentEnrolled;
}

void ZermeloAppointment::setIsCanceled(bool isCanceled)
{
    if (isCanceled != m_isCanceled) {
        m_isCanceled = isCanceled;
        emit isCanceledChanged();
    }
}

bool ZermeloAppointment::isCanceled() const
{
    return m_isCanceled;
}
