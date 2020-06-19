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
    Q_PROPERTY(MessageQueue::NatsStatus::Enum lastStatus READ lastStatus NOTIFY lastStatusChanged)
    Q_PROPERTY(NatsOptions *natsOptions READ natsOptions WRITE setNatsOptions NOTIFY natsOptionsChanged)
    Q_PROPERTY(int connectionWait READ connectionWait WRITE setConnectionWait NOTIFY connectionWaitChanged)
    Q_PROPERTY(QString discoveryPrefix READ discoveryPrefix WRITE setDiscoveryPrefix NOTIFY discoveryPrefixChanged)
    Q_PROPERTY(int maxPubAcksInflight READ maxPubAcksInflight WRITE setMaxPubAcksInflight NOTIFY maxPubAcksInflightChanged)
    Q_PROPERTY(float maxPubAcksInflightPercentage READ maxPubAcksInflightPercentage WRITE setMaxPubAcksInflightPercentage NOTIFY maxPubAcksInflightPercentageChanged)
    Q_PROPERTY(int pingsInterval READ pingsInterval WRITE setPingsInterval NOTIFY pingsIntervalChanged)
    Q_PROPERTY(int pingsMaxOut READ pingsMaxOut WRITE setPingsMaxOut NOTIFY pingsMaxOutChanged)
    Q_PROPERTY(int pubAckWait READ pubAckWait WRITE setPubAckWait NOTIFY pubAckWaitChanged)
    Q_PROPERTY(QString url READ url WRITE setUrl NOTIFY urlChanged)

public:
    explicit StanConnectionOptions(QObject *parent = nullptr);
    ~StanConnectionOptions() override = default;

    void setNatsOptions(NatsOptions *opts);
    [[nodiscard]] NatsOptions *natsOptions() const;
    void setConnectionWait(int wait);
    [[nodiscard]] int connectionWait() const;
    void setDiscoveryPrefix(const QString &prefix);
    [[nodiscard]] QString discoveryPrefix() const;
    void setMaxPubAcksInflight(int maxPubAcksInflight);
    [[nodiscard]] int maxPubAcksInflight() const;
    void setMaxPubAcksInflightPercentage(float percentage);
    [[nodiscard]] float maxPubAcksInflightPercentage() const;
    void setPingsInterval(int interval);
    [[nodiscard]] int pingsInterval() const;
    void setPingsMaxOut(int maxOut);
    [[nodiscard]] int pingsMaxOut() const;
    void setPubAckWait(int ms);
    [[nodiscard]] int pubAckWait() const;
    void setUrl(const QString &url);
    [[nodiscard]] QString url() const;

    QSharedPointer<stanConnOptions> connectionOptions();

    [[nodiscard]] NatsStatus::Enum lastStatus() const;

signals:
    void lastStatusChanged();
    void natsOptionsChanged();
    void connectionWaitChanged();
    void discoveryPrefixChanged();
    void maxPubAcksInflightChanged();
    void maxPubAcksInflightPercentageChanged();
    void pingsIntervalChanged();
    void pingsMaxOutChanged();
    void pubAckWaitChanged();
    void urlChanged();

    void errorOccurred(NatsStatus::Enum status, const QString &message);

private:
    void updateStatus(NatsStatus::Enum s);

    NatsStatus::Enum updateMaxPubAcksInflight();
    NatsStatus::Enum updatePings();

    NatsOptions *m_natsOptions = nullptr;
    int m_connectionWait = 0;
    QString m_discoveryPrefix;
    int m_maxPubAcksInflight = 0;
    float m_maxPubAcksInflightPercentage = 1;
    int m_pingsInterval = 0;
    int m_pingsMaxOut = 0;
    int m_pubAckWait = 0;
    QString m_url;

    QSharedPointer<stanConnOptions> m_connOptions;
    NatsStatus::Enum m_lastStatus = NatsStatus::Enum::Ok;
};

} // namespace MessageQueue

#endif // STANCONNECTIONOPTIONS_H
