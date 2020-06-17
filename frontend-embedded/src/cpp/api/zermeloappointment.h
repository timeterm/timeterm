#ifndef ZERMELOAPPOINTMENT_H
#define ZERMELOAPPOINTMENT_H

#include <QDateTime>
#include <QObject>

class ZermeloAppointment
{
Q_GADGET
    Q_PROPERTY(qint64 id WRITE setId READ id)
    Q_PROPERTY(qint64 appointmentInstance WRITE setAppointmentInstance READ appointmentInstance)
    Q_PROPERTY(qint32 startTimeSlot WRITE setStartTimeSlot READ startTimeSlot)
    Q_PROPERTY(qint32 endTimeSlot WRITE setEndTimeSlot READ endTimeSlot)
    Q_PROPERTY(qint32 capacity WRITE setCapacity READ capacity)
    Q_PROPERTY(qint32 availableSpace WRITE setAvailableSpace READ availableSpace)
    Q_PROPERTY(QDateTime startTime WRITE setStartTime READ startTime)
    Q_PROPERTY(QDateTime endTime WRITE setEndTime READ endTime)
    Q_PROPERTY(QStringList subjects WRITE setSubjects READ subjects)
    Q_PROPERTY(QStringList locations WRITE setLocations READ locations)
    Q_PROPERTY(QStringList teachers WRITE setTeachers READ teachers)
    Q_PROPERTY(bool isOnline WRITE setIsOnline READ isOnline)
    Q_PROPERTY(bool isOptional WRITE setIsOptional READ isOptional)
    Q_PROPERTY(bool isStudentEnrolled WRITE setIsStudentEnrolled READ isStudentEnrolled)
    Q_PROPERTY(bool isCanceled WRITE setIsCanceled READ isCanceled)

public:
    void setId(qint64 id);
    [[nodiscard]] qint64 id() const;
    void setAppointmentInstance(qint64 appointmentInstance);
    [[nodiscard]] qint64 appointmentInstance() const;
    void setStartTimeSlot(qint32 startTimeSlot);
    [[nodiscard]] qint32 startTimeSlot() const;
    void setEndTimeSlot(qint32 endTimeSlot);
    [[nodiscard]] qint32 endTimeSlot() const;
    void setCapacity(qint32 capacity);
    [[nodiscard]] qint32 capacity() const;
    void setAvailableSpace(qint32 availableSpace);
    [[nodiscard]] qint32 availableSpace() const;
    void setStartTime(const QDateTime &startTime);
    [[nodiscard]] QDateTime startTime() const;
    void setEndTime(const QDateTime &endTime);
    [[nodiscard]] QDateTime endTime() const;
    void setSubjects(const QStringList &subjects);
    [[nodiscard]] QStringList subjects() const;
    void setLocations(const QStringList &locations);
    [[nodiscard]] QStringList locations() const;
    void setTeachers(const QStringList &teachers);
    [[nodiscard]] QStringList teachers() const;
    void setIsOnline(bool isOnline);
    [[nodiscard]] bool isOnline() const;
    void setIsOptional(bool isOptional);
    [[nodiscard]] bool isOptional() const;
    void setIsStudentEnrolled(bool isStudentEnrolled);
    [[nodiscard]] bool isStudentEnrolled() const;
    void setIsCanceled(bool isCanceled);
    [[nodiscard]] bool isCanceled() const;

    void read(const QJsonObject &json);
    void write(QJsonObject &json) const;

private:
    qint64 m_id = 0;
    qint64 m_appointmentInstance = 0;
    qint32 m_startTimeSlot = 0;
    qint32 m_endTimeSlot = 0;
    qint32 m_capacity = 0;
    qint32 m_availableSpace = 0;
    QDateTime m_startTime;
    QDateTime m_endTime;
    QStringList m_subjects;
    QStringList m_locations;
    QStringList m_teachers;
    bool m_isOnline = false;
    bool m_isOptional = false;
    bool m_isStudentEnrolled = false;
    bool m_isCanceled = false;
};

Q_DECLARE_METATYPE(ZermeloAppointment)

#endif // ZERMELOAPPOINTMENT_H
