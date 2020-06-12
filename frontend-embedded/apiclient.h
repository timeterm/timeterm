#ifndef APICLIENT_H
#define APICLIENT_H

#include <QDateTime>
#include <QList>
#include <QObject>
#include <QQmlListProperty>

class ApiClient: public QObject
{
    Q_OBJECT
    Q_PROPERTY(QString cardId WRITE setCardId READ cardId NOTIFY cardIdChanged)

public:
    explicit ApiClient(QObject *parent = nullptr);

    void setCardId(const QString &cardId);

    QString cardId() const;

signals:
    void cardIdChanged();

private:
    QString m_cardId;
};

class Appointment: public QObject
{
    Q_OBJECT
    Q_PROPERTY(qint64 id WRITE setId READ id NOTIFY idChanged)
    Q_PROPERTY(qint64 appointmentInstance WRITE setAppointmentInstance READ appointmentInstance NOTIFY appointmentInstanceChanged)
    Q_PROPERTY(qint32 startTimeSlot WRITE setStartTimeSlot READ startTimeSlot NOTIFY startTimeSlotChanged)
    Q_PROPERTY(qint32 endTimeSlot WRITE setEndTimeSlot READ endTimeSlot NOTIFY endTimeSlotChanged)
    Q_PROPERTY(qint32 capacity WRITE setCapacity READ capacity NOTIFY capacityChanged)
    Q_PROPERTY(qint32 availableSpace WRITE setAvailableSpace READ availableSpace NOTIFY availableSpaceChanged)
    Q_PROPERTY(QDateTime startTime WRITE setStartTime READ startTime NOTIFY startTimeChanged)
    Q_PROPERTY(QDateTime endTime WRITE setEndTime READ endTime NOTIFY endTimeChanged)
    Q_PROPERTY(QStringList subjects WRITE setSubjects READ subjects NOTIFY subjectsChanged)
    Q_PROPERTY(QStringList locations WRITE setLocations READ locations NOTIFY locationsChanged)
    Q_PROPERTY(QStringList teachers WRITE setTeachers READ teachers NOTIFY teachersChanged)
    Q_PROPERTY(bool isOnline WRITE setIsOnline READ isOnline NOTIFY isOnlineChanged)
    Q_PROPERTY(bool isOptional WRITE setIsOptional READ isOptional NOTIFY isOptionalChanged)
    Q_PROPERTY(bool isStudentEnrolled WRITE setIsStudentEnrolled READ isStudentEnrolled NOTIFY isStudentEnrolledChanged)
    Q_PROPERTY(bool isCanceled WRITE setIsCanceled READ isCanceled NOTIFY isCanceledChanged)

public:
    explicit Appointment(QObject *parent = nullptr);

    void setId(qint64 id);
    qint64 id() const;
    void setAppointmentInstance(qint64 appointmentInstance);
    qint64 appointmentInstance() const;
    void setStartTimeSlot(qint32 startTimeSlot);
    qint32 startTimeSlot() const;
    void setEndTimeSlot(qint32 endTimeSlot);
    qint32 endTimeSlot() const;
    void setCapacity(qint32 capacity);
    qint32 capacity() const;
    void setAvailableSpace(qint32 availableSpace);
    qint32 availableSpace() const;
    void setStartTime(const QDateTime &startTime);
    QDateTime startTime() const;
    void setEndTime(const QDateTime &endTime);
    QDateTime endTime() const;
    void setSubjects(const QStringList &subjects);
    QStringList subjects() const;
    void setLocations(const QStringList &locations);
    QStringList locations() const;
    void setTeachers(const QStringList &teachers);
    QStringList teachers() const;
    void setIsOnline(bool isOnline);
    bool isOnline() const;
    void setIsOptional(bool isOptional);
    bool isOptional() const;
    void setIsStudentEnrolled(bool isStudentEnrolled);
    bool isStudentEnrolled() const;
    void setIsCanceled(bool isCanceled);
    bool isCanceled() const;

signals:
    void idChanged();
    void appointmentInstanceChanged();
    void startTimeSlotChanged();
    void endTimeSlotChanged();
    void capacityChanged();
    void availableSpaceChanged();
    void startTimeChanged();
    void endTimeChanged();
    void subjectsChanged();
    void locationsChanged();
    void teachersChanged();
    void isOnlineChanged();
    void isOptionalChanged();
    void isStudentEnrolledChanged();
    void isCanceledChanged();

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

class Appointments: public QObject
{
    Q_OBJECT
    Q_PROPERTY(QQmlListProperty<Appointment> appointments READ appointments NOTIFY appointmentsChanged)

public:
    explicit Appointments(QObject *parent = nullptr);

    void append(Appointment *appointment);

    QQmlListProperty<Appointment> appointments();

signals:
    void appointmentsChanged();

private:
    QList<Appointment *> m_appointments;
    void append(const QList<Appointment *> &appointments);
};

#endif//APICLIENT_H
