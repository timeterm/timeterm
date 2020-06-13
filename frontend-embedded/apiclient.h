#ifndef APICLIENT_H
#define APICLIENT_H

#include <QDateTime>
#include <QList>
#include <QNetworkAccessManager>
#include <QObject>
#include <QQmlListProperty>

class ZermeloAppointment
{
    Q_GADGET
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
    explicit ZermeloAppointment(QObject *parent = nullptr);

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

class ZermeloAppointments
{
    Q_GADGET
    Q_PROPERTY(QList<ZermeloAppointment> data READ data NOTIFY dataChanged)

public:
    void append(const ZermeloAppointment& appointment);

    QList<ZermeloAppointment> data();

signals:
    void dataChanged();

private:
    void append(const QList<ZermeloAppointment> &appointments);

    QList<ZermeloAppointment> m_data;
};

class TimetermUser
{
    Q_GADGET

public:
    void setCardUid(const QString &cardUid);
    QString cardUid();
    void setOrganizationId(const QString &organizationId);
    QString organizationId();
    void setName(const QString &name);
    QString name();
    void setStudentCode(const QString &studentCode);
    QString studentCode();

signals:
    void cardUidChanged();
    void organizationIdChanged();
    void nameChanged();
    void studentCodeChanged();

private:
    QString m_cardUid;
    QString m_organizationId;
    QString m_name;
    QString m_studentCode;
};

class ApiClient: public QObject
{
    Q_OBJECT
    Q_PROPERTY(QString cardId WRITE setCardId READ cardId NOTIFY cardIdChanged)
    Q_PROPERTY(QString apiKey WRITE setApiKey READ apiKey NOTIFY apiKeyChanged)

public:
    explicit ApiClient(QObject *parent = nullptr);

    void setCardId(const QString &cardId);
    QString cardId() const;
    void setApiKey(const QString &apiKey);
    QString apiKey() const;

    Q_INVOKABLE void getCurrentUser();

signals:
    void cardIdChanged();
    void apiKeyChanged();

    void currentUserReceived(TimetermUser);

private:
    QUrl m_baseUrl = QUrl("https://timeterm.nl/api/v1/");
    QString m_cardId;
    QString m_apiKey;
    QNetworkAccessManager *m_qnam;
    void setAuthHeaders(QNetworkRequest &req);
};

#endif//APICLIENT_H
