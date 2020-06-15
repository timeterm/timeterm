#include "apiclient.h"

#include <QJsonObject>
#include <QJsonParseError>
#include <QNetworkReply>

ApiClient::ApiClient(QObject *parent)
    : QObject(parent),
      m_qnam(new QNetworkAccessManager(this))
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

QString ApiClient::apiKey() const
{
    return m_apiKey;
}

void ApiClient::setApiKey(const QString &apiKey)
{
    if (apiKey != m_apiKey) {
        m_apiKey = apiKey;
        emit apiKeyChanged();
    }
}

void ApiClient::getTimetable()
{
}

void ApiClient::getCurrentUser()
{
    auto req = QNetworkRequest(m_baseUrl.resolved(QUrl("user/self")));
    setAuthHeaders(req);

    auto reply = m_qnam->get(req);
    connectReply(reply, &ApiClient::handleCurrentUserReply);
}

void ApiClient::connectReply(QNetworkReply *reply, ReplyHandler handler)
{
    m_handlers[reply] = handler;

    connect(reply, &QNetworkReply::finished, this, &ApiClient::replyFinished);
    connect(reply, &QNetworkReply::errorOccurred, this, &ApiClient::handleReplyError);
}

void ApiClient::setAuthHeaders(QNetworkRequest &req)
{
    req.setRawHeader("X-Api-Key", m_apiKey.toLocal8Bit());
    req.setRawHeader("X-Card-Uid", m_cardId.toLocal8Bit());
}

void ApiClient::replyFinished()
{
    auto reply = qobject_cast<QNetworkReply *>(QObject::sender());

    (this->*m_handlers[reply])(reply);
    m_handlers.remove(reply);

    reply->deleteLater();
}

void ApiClient::handleCurrentUserReply(QNetworkReply *reply)
{
    auto bytes = reply->readAll();
    auto json = QJsonDocument::fromJson(bytes);

    if (!json.isObject())
        return;

    auto user = TimetermUser();
    user.read(json.object());

    emit currentUserReceived(user);
}

void ApiClient::handleReplyError(QNetworkReply::NetworkError error)
{
    auto reply = qobject_cast<QNetworkReply *>(QObject::sender());

    m_handlers.remove(reply);

    reply->deleteLater();
}

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

void ZermeloAppointment::setId(qint64 id)
{
    if (id != m_id) {
        m_id = id;
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
    }
}

bool ZermeloAppointment::isCanceled() const
{
    return m_isCanceled;
}

void TimetermUser::setCardUid(const QString &cardUid)
{
    if (cardUid != m_cardUid) {
        m_cardUid = cardUid;
    }
}

QString TimetermUser::cardUid() const
{
    return m_cardUid;
}

void TimetermUser::setOrganizationId(const QString &organizationId)
{
    if (organizationId != m_organizationId) {
        m_organizationId = organizationId;
    }
}

QString TimetermUser::organizationId() const
{
    return m_organizationId;
}

void TimetermUser::setName(const QString &name)
{
    if (name != m_name) {
        m_name = name;
    }
}

QString TimetermUser::name() const
{
    return m_name;
}

void TimetermUser::setStudentCode(const QString &studentCode)
{
    if (studentCode != m_studentCode) {
        m_studentCode = studentCode;
    }
}

QString TimetermUser::studentCode() const
{
    return m_studentCode;
}

void TimetermUser::read(const QJsonObject &json)
{
    if (json.contains("cardUid") && json["cardUid"].isString())
        setCardUid(json["cardUid"].toString());

    if (json.contains("name") && json["name"].isString())
        setName(json["name"].toString());

    if (json.contains("organizationId") && json["organizationId"].isString())
        setOrganizationId(json["organizationId"].toString());

    if (json.contains("studentCode") && json["studentCode"].isString())
        setStudentCode(json["studentCode"].toString());
}

void TimetermUser::write(QJsonObject &json) const
{
    json["cardUid"] = cardUid();
    json["name"] = name();
    json["organizationId"] = organizationId();
    json["studentCode"] = studentCode();
}
