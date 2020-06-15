#include "apiclient.h"

#include <QJsonArray>
#include <QJsonObject>
#include <QJsonParseError>
#include <QNetworkReply>
#include <QUrlQuery>
#include <optional>

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

void setTimetableQueryParams(QUrl &url, const QDateTime &start, const QDateTime &end)
{
    auto query = QUrlQuery(url);
    query.setQueryItems({
        {"startTime", start.toString()},
        {"endTime", end.toString()},
    });
    url.setQuery(query);
}

void ApiClient::getAppointments(const QDateTime &start, const QDateTime &end)
{
    auto url = m_baseUrl.resolved(QUrl("zermelo/appointment"));
    setTimetableQueryParams(url, start, end);

    auto req = QNetworkRequest(url);
    setAuthHeaders(req);

    auto reply = m_qnam->get(req);
    connectReply(reply, &ApiClient::handleGetAppointmentsReply);
}

void ApiClient::getCurrentUser()
{
    auto req = QNetworkRequest(m_baseUrl.resolved(QUrl("user/self")));
    setAuthHeaders(req);

    auto reply = m_qnam->get(req);
    connectReply(reply, &ApiClient::handleGetCurrentUserReply);
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

    auto methodPointer = m_handlers[reply];
    if (methodPointer != nullptr) {
        (this->*methodPointer)(reply);
        m_handlers.remove(reply);
    }

    reply->deleteLater();
}

template<typename T>
std::optional<T> readJsonObject(QNetworkReply *reply)
{
    auto bytes = reply->readAll();
    auto json = QJsonDocument::fromJson(bytes);

    if (!json.isObject())
        return std::nullopt;

    auto decoded = T();
    decoded.read(json.object());
    return decoded;
}

void ApiClient::handleGetCurrentUserReply(QNetworkReply *reply)
{
    auto user = readJsonObject<TimetermUser>(reply);
    if (!user.has_value())
        return;

    emit currentUserReceived(user.value());
}

void ApiClient::handleGetAppointmentsReply(QNetworkReply *reply)
{
    auto user = readJsonObject<ZermeloAppointments>(reply);
    if (!user.has_value())
        return;

    emit timetableReceived(user.value());
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

void ZermeloAppointment::setStartTimeSlot(qint32 startTimeSlot)
{
    if (startTimeSlot != m_startTimeSlot)
        m_startTimeSlot = startTimeSlot;
}

qint32 ZermeloAppointment::startTimeSlot() const
{
    return m_startTimeSlot;
}

void ZermeloAppointment::setEndTimeSlot(qint32 endTimeSlot)
{
    if (endTimeSlot != m_endTimeSlot)
        m_endTimeSlot = endTimeSlot;
}

qint32 ZermeloAppointment::endTimeSlot() const
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

    if (json.contains("startTimeSlot") && json["startTimeSlot"].isDouble())
        m_startTimeSlot = json["startTimeSlot"].toInt();

    if (json.contains("endTimeSlot") && json["endTimeSlot"].isDouble())
        m_endTimeSlot = json["endTimeSlot"].toInt();

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
    json["locations"] = stringListAsQJsonArray(m_locations);
    json["teachers"] = stringListAsQJsonArray(m_teachers);
    json["isOnline"] = m_isOnline;
    json["isStudentEnrolled"] = m_isStudentEnrolled;
    json["isCanceled"] = m_isCanceled;
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
