#include "apiclient.h"
#include "createdevice.h"

#include <optional>
#include <utility>

#include <QJsonArray>
#include <QJsonObject>
#include <QJsonParseError>
#include <QNetworkReply>
#include <QUrlQuery>

ApiClient::ApiClient(QObject *parent)
    : QObject(parent)
    , m_qnam(new QNetworkAccessManager(this))
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
    connectReply(reply, [this](QNetworkReply *reply) {
        return handleGetAppointmentsReply(reply);
    });
}

void ApiClient::getCurrentUser()
{
    auto req = QNetworkRequest(m_baseUrl.resolved(QUrl("user/self")));
    setAuthHeaders(req);

    auto reply = m_qnam->get(req);
    connectReply(reply, [this](QNetworkReply *reply) {
        handleGetCurrentUserReply(reply);
    });
}

void ApiClient::connectReply(QNetworkReply *reply, ReplyHandler handler)
{
    m_handlers[reply] = std::move(handler);

    reply->setParent(this);
    connect(reply, &QNetworkReply::finished, this, &ApiClient::replyFinished);
    connect(reply, &QNetworkReply::errorOccurred, this, &ApiClient::handleReplyError);
}

void ApiClient::setAuthHeaders(QNetworkRequest &req)
{
    req.setRawHeader("X-Api-Key", m_apiKey.toLocal8Bit());
    if (m_cardId != "")
        req.setRawHeader("X-Card-Uid", m_cardId.toLocal8Bit());
}

void ApiClient::replyFinished()
{
    auto reply = qobject_cast<QNetworkReply *>(QObject::sender());

    auto handler = m_handlers[reply];
    if (handler != nullptr) {
        handler(reply);
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

void ApiClient::createDevice()
{
    auto reqData = CreateDeviceRequest();
    reqData.name = "Nieuw apparaat";

    QJsonObject reqJson;
    reqData.write(reqJson);

    auto reqBytes = QJsonDocument(reqJson).toJson();
    auto req = QNetworkRequest(m_baseUrl.resolved(QUrl("device")));
    setAuthHeaders(req);
    req.setHeader(QNetworkRequest::ContentTypeHeader, "application/json");

    auto reply = m_qnam->post(req, reqBytes);
    connectReply(reply, [this](QNetworkReply *reply) {
      handleCreateDeviceReply(reply);
    });
}

void ApiClient::handleCreateDeviceReply(QNetworkReply *reply)
{
    auto rsp = readJsonObject<CreateDeviceResponse>(reply);
    if (!rsp.has_value())
        return;

    emit deviceCreated(rsp.value());
}

void ApiClient::getNatsCreds(const QString& deviceId)
{
    auto url = m_baseUrl.resolved(QUrl(QStringLiteral("device/%1/config/natscreds").arg(deviceId)));
    auto req = QNetworkRequest(url);
    setAuthHeaders(req);

    auto reply = m_qnam->get(req);
    connectReply(reply, [this](QNetworkReply *reply) {
      return handleNatsCredsReply(reply);
    });
}

void ApiClient::handleNatsCredsReply(QNetworkReply *reply) {
    auto rsp = readJsonObject<NatsCredsResponse>(reply);
    if (!rsp.has_value())
        return;

    emit natsCredsReceived(rsp.value());
}
