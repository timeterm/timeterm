#pragma once

#include "createdevice.h"
#include "natscreds.h"
#include "servicesresponse.h"
#include "timetermuser.h"
#include "zermeloappointments.h"

#include <QJsonObject>
#include <QNetworkAccessManager>
#include <QNetworkReply>
#include <QObject>
#include <devcfg/connmanserviceconfig.h>

class ApiClient: public QObject
{
    Q_OBJECT
    Q_PROPERTY(QString cardId WRITE setCardId READ cardId NOTIFY cardIdChanged)
    Q_PROPERTY(QString apiKey WRITE setApiKey READ apiKey NOTIFY apiKeyChanged)

    using ReplyHandler = std::function<void(QNetworkReply *)>;

public:
    explicit ApiClient(QObject *parent = nullptr);
    ~ApiClient() override = default;

    void setCardId(const QString &cardId);
    [[nodiscard]] QString cardId() const;
    void setApiKey(const QString &apiKey);
    [[nodiscard]] QString apiKey() const;

    Q_INVOKABLE void getCurrentUser();
    Q_INVOKABLE void getAppointments(const QDateTime &start, const QDateTime &end);
    Q_INVOKABLE void createDevice();
    Q_INVOKABLE void getNatsCreds(const QString &deviceId);
    Q_INVOKABLE void doHeartbeat(const QString &deviceId);
    Q_INVOKABLE void getAllNetworkingServices(const QString &deviceId);
    Q_INVOKABLE void updateChoice(const QVariant &unenrollFromParticipationId, const QVariant &enrollIntoParticipationId);

signals:
    void cardIdChanged();
    void apiKeyChanged();
    void currentUserReceived(TimetermUser);
    void timetableReceived(ZermeloAppointments);
    void deviceCreated(CreateDeviceResponse);
    void natsCredsReceived(NatsCredsResponse);
    void heartbeatSucceeded();
    void choiceUpdateSucceeded();
    void newNetworkingServices(NetworkingServicesResponse);

private slots:
    void replyFinished();
    void handleReplyError(QNetworkReply::NetworkError error);

private:
    void connectReply(QNetworkReply *reply, ReplyHandler handler);
    void handleGetCurrentUserReply(QNetworkReply *reply);
    void handleGetAppointmentsReply(QNetworkReply *reply);
    void handleCreateDeviceReply(QNetworkReply *reply);
    void handleNatsCredsReply(QNetworkReply *reply);
    void handleHeartbeatReply(QNetworkReply *reply);
    void handleChoiceUpdateReply(QNetworkReply *reply);
    void handleNewNetworkingServices(QNetworkReply *reply);
    void setAuthHeaders(QNetworkRequest &req);

    QUrl m_baseUrl = QUrl("https://api.timeterm.nl/");
    QString m_cardId;
    QString m_apiKey;
    QNetworkAccessManager *m_qnam;
    QHash<QNetworkReply *, ReplyHandler> m_handlers;
};

class ApiError
{
    Q_GADGET
    Q_PROPERTY(QString message MEMBER message)

public:
    void read(const QJsonObject &json);

    QString message;
};
