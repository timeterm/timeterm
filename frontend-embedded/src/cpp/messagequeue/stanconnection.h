#pragma once

#include "enums.h"
#include "stanconnectionoptions.h"
#include "stansuboptions.h"

#include <QHash>
#include <QObject>

#include <QtQml/qqml.h>
#include <functional>
#include <nats.h>

namespace MessageQueue
{

class StanConnection: public QObject
{
    Q_OBJECT
    Q_PROPERTY(MessageQueue::NatsStatus::Enum lastStatus READ lastStatus)
    Q_PROPERTY(QString cluster READ cluster WRITE setCluster NOTIFY clusterChanged)
    Q_PROPERTY(QString clientId READ clientId WRITE setClientId NOTIFY clientIdChanged)
    Q_PROPERTY(MessageQueue::StanConnectionOptions *connectionOptions READ connectionOptions WRITE setConnectionOptions NOTIFY connectionOptionsChanged)

public:
    explicit StanConnection(QObject *parent = nullptr);
    ~StanConnection() override;

    [[nodiscard]] NatsStatus::Enum lastStatus() const;
    void setCluster(const QString &cluster);
    [[nodiscard]] QString cluster() const;
    void setClientId(const QString &clientId);
    [[nodiscard]] QString clientId() const;
    void setConnectionOptions(StanConnectionOptions *options);
    [[nodiscard]] StanConnectionOptions *connectionOptions() const;

    Q_INVOKABLE void connect();
    NatsStatus::Enum subscribe(StanSubOptions *opts, stanSubscription **ppStanSub, QSharedPointer<stanConnection*> &spConn);

signals:
    void errorOccurred(MessageQueue::NatsStatus::Enum s, const QString &message);
    void clusterChanged();
    void clientIdChanged();
    void connectionOptionsChanged();
    void connected();
    void connectionLost();
    void setConnectionPrivate(const QSharedPointer<stanConnection*> &conn, QPrivateSignal);
    void lastStatusChanged();

private slots:
    void setConnection(const QSharedPointer<stanConnection*> &conn);

private:
    void updateStatus(NatsStatus::Enum s);

    NatsStatus::Enum m_lastStatus = NatsStatus::Enum::Ok;

    QString m_cluster;
    QString m_clientId;
    StanConnectionOptions *m_options;

    QSharedPointer<stanConnection *> m_stanConnection;
};

} // namespace MessageQueue

Q_DECLARE_METATYPE(QSharedPointer<stanConnection*>)
