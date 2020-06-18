#ifndef STANCONNECTIONOPTIONS_H
#define STANCONNECTIONOPTIONS_H

#include "enums.h"
#include "natsoptions.h"
#include "scopedpointer.h"

#include <nats.h>

#include <QObject>
#include <QSharedPointer>

namespace MessageQueue
{

class StanConnectionOptions: public QObject
{
Q_OBJECT
    Q_PROPERTY(NatsStatus::Enum lastStatus READ lastStatus)

public:
    explicit StanConnectionOptions(QObject *parent = nullptr);

    Q_INVOKABLE NatsStatus::Enum setNatsOptions(NatsOptions *opts);
    Q_INVOKABLE NatsStatus::Enum setConnectionWait(qint64 wait);
    Q_INVOKABLE NatsStatus::Enum setDiscoveryPrefix(const QString &prefix);
    Q_INVOKABLE NatsStatus::Enum setMaxPubAcksInflight(int maxPubAcksInflight, float percentage);
    Q_INVOKABLE NatsStatus::Enum setPings(int interval, int maxOut);
    Q_INVOKABLE NatsStatus::Enum setPubAckWait(qint64 ms);
    Q_INVOKABLE NatsStatus::Enum setUrl(const QString &url);

    QSharedPointer<stanConnOptions> connectionOptions();

    [[nodiscard]] NatsStatus::Enum lastStatus() const;

signals:
    void errorOccurred(NatsStatus::Enum status, const QString &message);

private:
    void updateStatus(NatsStatus::Enum s);

    QSharedPointer<stanConnOptions> m_connOptions;
    NatsStatus::Enum m_lastStatus = NatsStatus::Enum::Ok;
};

} // namespace MessageQueue

#endif // STANCONNECTIONOPTIONS_H
