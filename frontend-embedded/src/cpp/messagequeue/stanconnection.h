#ifndef STANCONNECTION_H
#define STANCONNECTION_H

#include "enums.h"
#include "scopedpointer.h"
#include "stansuboptions.h"
#include "stansubscription.h"

#include <QHash>
#include <QObject>

#include <functional>
#include <nats.h>

namespace MessageQueue
{

using StanConnectionDeleter = ScopedPointerDestroyerDeleter<stanConnection, natsStatus, &stanConnection_Destroy>;
using StanConnOptionsDeleter = ScopedPointerDestroyerDeleter<stanConnOptions, void, &stanConnOptions_Destroy>;
using NatsOptionsDeleter = ScopedPointerDestroyerDeleter<natsOptions, void, &natsOptions_Destroy>;

using StanConnectionScopedPointer = QScopedPointer<stanConnection, StanConnectionDeleter>;
using NatsOptionsScopedPointer = QScopedPointer<natsOptions, NatsOptionsDeleter>;
using StanConnOptionsScopedPointer = QScopedPointer<stanConnOptions, StanConnOptionsDeleter>;

class StanConnection: public QObject
{
    Q_OBJECT
    Q_PROPERTY(NatsStatus lastStatus READ lastStatus)
    Q_PROPERTY(QString cluster READ cluster WRITE setCluster NOTIFY clusterChanged)
    Q_PROPERTY(QString clientId READ clientId WRITE setClientId NOTIFY clientIdChanged)

public:
    explicit StanConnection(QObject *parent = nullptr);

    [[nodiscard]] NatsStatus lastStatus() const;
    void setCluster(const QString &cluster);
    QString cluster() const;
    void setClientId(const QString &clientId);
    QString clientId() const;

    Q_INVOKABLE void connect();
    Q_INVOKABLE StanSubscription *subscribe(const QString &subscribe, StanSubOptions *opts);

signals:
    void errorOccurred(NatsStatus s, const QString &message);
    void clusterChanged();
    void clientIdChanged();

private:
    void updateStatus(NatsStatus s);

    NatsStatus m_lastStatus = NatsStatus::Ok;

    QString m_cluster;
    QString m_clientId;

    StanConnectionScopedPointer m_stanConnection;
    NatsOptionsScopedPointer m_natsOpts;
    StanConnOptionsScopedPointer m_connOpts;
};

} // namespace MessageQueue

#endif // STANCONNECTION_H
