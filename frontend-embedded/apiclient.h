#ifndef APICLIENT_H
#define APICLIENT_H

#include <QDateTime>
#include <QList>
#include <QNetworkAccessManager>
#include <QNetworkReply>
#include <QObject>
#include <QQmlListProperty>

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

class TimetermUser
{
    Q_GADGET
    Q_PROPERTY(QString cardUid READ cardUid WRITE setCardUid)
    Q_PROPERTY(QString organizationId READ organizationId WRITE setOrganizationId)
    Q_PROPERTY(QString name READ name WRITE setName)
    Q_PROPERTY(QString studentCode READ studentCode WRITE setStudentCode)

public:
    void setCardUid(const QString &cardUid);
    [[nodiscard]] QString cardUid() const;
    void setOrganizationId(const QString &organizationId);
    [[nodiscard]] QString organizationId() const;
    void setName(const QString &name);
    [[nodiscard]] QString name() const;
    void setStudentCode(const QString &studentCode);
    [[nodiscard]] QString studentCode() const;

    void read(const QJsonObject &json);
    void write(QJsonObject &json) const;

private:
    QString m_cardUid;
    QString m_organizationId;
    QString m_name;
    QString m_studentCode;
};

Q_DECLARE_METATYPE(TimetermUser)

class ApiClient;

using ReplyHandler = void(ApiClient::*)(QNetworkReply*);

class ApiClient: public QObject
{
    Q_OBJECT
    Q_PROPERTY(QString cardId WRITE setCardId READ cardId NOTIFY cardIdChanged)
    Q_PROPERTY(QString apiKey WRITE setApiKey READ apiKey NOTIFY apiKeyChanged)

public:
    explicit ApiClient(QObject *parent = nullptr);

    void setCardId(const QString &cardId);
    [[nodiscard]] QString cardId() const;
    void setApiKey(const QString &apiKey);
    [[nodiscard]] QString apiKey() const;

    Q_INVOKABLE void getCurrentUser();
    Q_INVOKABLE void getTimetable(const QDateTime& start, const QDateTime& end);

signals:
    void cardIdChanged();
    void apiKeyChanged();
    void currentUserReceived(TimetermUser);
    void timetableReceived(ZermeloAppointments);

private slots:
    void replyFinished();
    void handleReplyError(QNetworkReply::NetworkError error);

private:
    void connectReply(QNetworkReply *reply, ReplyHandler handler);
    void handleCurrentUserReply(QNetworkReply* reply);
    void handleTimetableReply(QNetworkReply *reply);
    void setAuthHeaders(QNetworkRequest &req);

    QUrl m_baseUrl = QUrl("https://timeterm.nl/api/v1/");
    QString m_cardId;
    QString m_apiKey;
    QNetworkAccessManager *m_qnam;
    QHash<QNetworkReply *, ReplyHandler> m_handlers;
};

#endif//APICLIENT_H
