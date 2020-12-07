#pragma once

#include <QDateTime>
#include <QObject>

class ZermeloAppointment
{
    Q_GADGET
    Q_PROPERTY(qint64 id WRITE setId READ id)
    Q_PROPERTY(qint64 appointmentInstance WRITE setAppointmentInstance READ appointmentInstance)
    Q_PROPERTY(QString startTimeSlotName WRITE setStartTimeSlotName READ startTimeSlotName)
    Q_PROPERTY(QString endTimeSlotName WRITE setEndTimeSlotName READ endTimeSlotName)
    Q_PROPERTY(qint32 capacity WRITE setCapacity READ capacity)
    Q_PROPERTY(qint32 availableSpace WRITE setAvailableSpace READ availableSpace)
    Q_PROPERTY(QDateTime startTime WRITE setStartTime READ startTime)
    Q_PROPERTY(QDateTime endTime WRITE setEndTime READ endTime)
    Q_PROPERTY(QStringList subjects WRITE setSubjects READ subjects)
    Q_PROPERTY(QStringList groups WRITE setGroups READ groups)
    Q_PROPERTY(QStringList locations WRITE setLocations READ locations)
    Q_PROPERTY(QStringList teachers WRITE setTeachers READ teachers)
    Q_PROPERTY(bool isOnline WRITE setIsOnline READ isOnline)
    Q_PROPERTY(bool isOptional WRITE setIsOptional READ isOptional)
    Q_PROPERTY(bool isStudentEnrolled WRITE setIsStudentEnrolled READ isStudentEnrolled)
    Q_PROPERTY(bool isCanceled WRITE setIsCanceled READ isCanceled)
    Q_PROPERTY(QString content WRITE setContent READ content)
    Q_PROPERTY(QList<ZermeloAppointment> alternatives READ alternatives)

public:
    void setId(qint64 id);
    [[nodiscard]] qint64 id() const;
    void setAppointmentInstance(qint64 appointmentInstance);
    [[nodiscard]] qint64 appointmentInstance() const;
    void setStartTimeSlotName(const QString& startTimeSlotName);
    [[nodiscard]] QString startTimeSlotName() const;
    void setEndTimeSlotName(const QString& endTimeSlotName);
    [[nodiscard]] QString endTimeSlotName() const;
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
    void setGroups(const QStringList &groups);
    [[nodiscard]] QStringList groups() const;
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
    void setContent(const QString &content);
    [[nodiscard]] QString content() const;
    void appendAlternative(const ZermeloAppointment &appointment);
    QList<ZermeloAppointment> alternatives();

    void read(const QJsonObject &json);
    void write(QJsonObject &json) const;

private:
    void appendAlternatives(const QList<ZermeloAppointment> &alternatives);

    qint64 m_id = 0;
    qint64 m_appointmentInstance = 0;
    QString m_startTimeSlotName;
    QString m_endTimeSlotName;
    qint32 m_capacity = 0;
    qint32 m_availableSpace = 0;
    QDateTime m_startTime;
    QDateTime m_endTime;
    QStringList m_subjects;
    QStringList m_groups;
    QStringList m_locations;
    QStringList m_teachers;
    bool m_isOnline = false;
    bool m_isOptional = false;
    bool m_isStudentEnrolled = false;
    bool m_isCanceled = false;
    QString m_content;
    QList<ZermeloAppointment> m_alternatives;
};

Q_DECLARE_METATYPE(ZermeloAppointment)
