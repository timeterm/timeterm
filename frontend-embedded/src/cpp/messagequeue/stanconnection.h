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

public:
    explicit StanConnection(QObject *parent = nullptr);

    [[nodiscard]] NatsStatus lastStatus() const;

    Q_INVOKABLE void connect();
    Q_INVOKABLE StanSubscription *subscribe(const QString &subscribe, StanSubOptions *opts);

signals:
    void errorOccurred(NatsStatus s, const QString &message);

private:
    void updateStatus(NatsStatus s);

    NatsStatus m_lastStatus = NatsStatus::Ok;

    QString m_cluster = "test-cluster";
    QString m_clientId = "client";

    StanConnectionScopedPointer m_stanConnection;
    NatsOptionsScopedPointer m_natsOpts;
    StanConnOptionsScopedPointer m_connOpts;
};

} // namespace MessageQueue

#endif // STANCONNECTION_H
