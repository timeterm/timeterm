#pragma once

#include "timetermuser.h"
#include "zermeloappointments.h"

#include <QNetworkAccessManager>
#include <QNetworkReply>
#include <QObject>

class FakeApiClient: public QObject
{
    Q_OBJECT

public:
    explicit FakeApiClient(QObject *parent = nullptr);
    ~FakeApiClient() override = default;

    Q_INVOKABLE void getCurrentUser();
    Q_INVOKABLE void getAppointments(const QDateTime &start, const QDateTime &end);

signals:
    void currentUserReceived(TimetermUser);
    void timetableReceived(ZermeloAppointments);
};
