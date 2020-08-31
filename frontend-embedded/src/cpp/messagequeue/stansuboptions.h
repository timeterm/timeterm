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
    Q_PROPERTY(MessageQueue::NatsStatus::Enum lastStatus READ lastStatus)

public:
    explicit StanSubOptions(QObject *parent = nullptr);

    Q_INVOKABLE MessageQueue::NatsStatus::Enum setDurableName(const QString &durableName);
    Q_INVOKABLE MessageQueue::NatsStatus::Enum deliverAllAvailable();
    Q_INVOKABLE MessageQueue::NatsStatus::Enum startWithLastReceived();
    Q_INVOKABLE MessageQueue::NatsStatus::Enum startAtSequence(quint64 sequence);
    Q_INVOKABLE MessageQueue::NatsStatus::Enum setManualAckMode(bool manualAck);
    Q_INVOKABLE MessageQueue::NatsStatus::Enum setMaxInflight(int inflight);
    Q_INVOKABLE MessageQueue::NatsStatus::Enum setAckWait(qint64 ms);

    QSharedPointer<stanSubOptions> subOptions();

    [[nodiscard]] MessageQueue::NatsStatus::Enum lastStatus() const;

signals:
    void errorOccurred(NatsStatus::Enum status, const QString &message);

private:
    void updateStatus(NatsStatus::Enum s);

    QSharedPointer<stanSubOptions> m_subOptions;
    NatsStatus::Enum m_lastStatus = NatsStatus::Enum::Ok;
};

} // namespace MessageQueue

#endif // STANSUBOPTIONS_H
