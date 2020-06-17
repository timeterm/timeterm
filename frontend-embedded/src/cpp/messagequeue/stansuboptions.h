#ifndef STANSUBOPTIONS_H
#define STANSUBOPTIONS_H

#include "enums.h"
#include "scopedpointer.h"

#include <nats.h>

#include <QObject>
#include <QSharedPointer>

namespace MessageQueue
{

class StanSubOptions: public QObject
{
    Q_OBJECT
    Q_PROPERTY(NatsStatus lastStatus READ lastStatus)

public:
    explicit StanSubOptions(QObject *parent = nullptr);

    Q_INVOKABLE NatsStatus setDurableName(const QString &durableName);
    Q_INVOKABLE NatsStatus deliverAllAvailable();
    Q_INVOKABLE NatsStatus startWithLastReceived();
    Q_INVOKABLE NatsStatus startAtSequence(quint64 sequence);
    Q_INVOKABLE NatsStatus setManualAckMode(bool manualAck);
    Q_INVOKABLE NatsStatus setMaxInflight(int inflight);
    Q_INVOKABLE NatsStatus setAckWait(qint64 ms);

    QSharedPointer<stanSubOptions> subOptions();

    [[nodiscard]] NatsStatus lastStatus() const;

signals:
    void errorOccurred(NatsStatus status, const QString &message);

private:
    void updateStatus(NatsStatus s);

    QSharedPointer<stanSubOptions> m_subOptions;
    NatsStatus m_lastStatus = NatsStatus::Ok;
};

} // namespace MessageQueue

#endif // STANSUBOPTIONS_H
