#pragma once

#include "enums.h"

#include <nats.h>

#include <QObject>
#include <QSharedPointer>

namespace MessageQueue
{

class StanSubOptions: public QObject
{
    Q_OBJECT
    Q_PROPERTY(QString durableName READ durableName WRITE setDurableName NOTIFY durableNameChanged)
    Q_PROPERTY(bool deliverAllAvailable READ deliverAllAvailable WRITE setDeliverAllAvailable NOTIFY deliverAllAvailableChanged)
    Q_PROPERTY(bool startWithLastReceived READ startWithLastReceived WRITE setStartWithLastReceived NOTIFY startWithLastReceivedChanged)
    Q_PROPERTY(quint64 startAtSequence READ startAtSequence WRITE setStartAtSequence NOTIFY startAtSequenceChanged)
    Q_PROPERTY(bool manualAckMode READ manualAckMode WRITE setManualAckMode NOTIFY manualAckModeChanged)
    Q_PROPERTY(int maxInflight READ maxInflight WRITE setMaxInflight NOTIFY maxInflightChanged)
    Q_PROPERTY(qint64 ackWaitMs READ ackWaitMs WRITE setAckWaitMs NOTIFY ackWaitMsChanged)
    Q_PROPERTY(QString channel READ channel WRITE setChannel NOTIFY channelChanged)

public:
    explicit StanSubOptions(QObject *parent = nullptr);

    [[nodiscard]] QString durableName() const;
    void setDurableName(const QString &durableName);
    [[nodiscard]] bool deliverAllAvailable() const;
    void setDeliverAllAvailable(bool deliverAllAvailable);
    [[nodiscard]] bool startWithLastReceived() const;
    void setStartWithLastReceived(bool startWithLastReceived);
    [[nodiscard]] quint64 startAtSequence() const;
    void setStartAtSequence(bool startAtSequence);
    [[nodiscard]] bool manualAckMode() const;
    void setManualAckMode(bool manualAckMode);
    [[nodiscard]] int maxInflight() const;
    void setMaxInflight(int maxInflight);
    [[nodiscard]] qint64 ackWaitMs() const;
    void setAckWaitMs(qint64 ackWaitMs);
    [[nodiscard]] QString channel() const;
    void setChannel(const QString &channel);

    NatsStatus::Enum build(stanSubOptions **ppSubOpts);

signals:
    void durableNameChanged();
    void deliverAllAvailableChanged();
    void startWithLastReceivedChanged();
    void startAtSequenceChanged();
    void manualAckModeChanged();
    void maxInflightChanged();
    void ackWaitMsChanged();
    void channelChanged();

private:
    NatsStatus::Enum configureSubOpts(stanSubOptions *pSubOpts);

    QString m_durableName;
    QString m_channel;
    bool m_deliverAllAvailable = false;
    bool m_startWithLastReceived = false;
    quint64 m_startAtSequence = 0;
    bool m_isStartAtSequenceSet = false;
    bool m_manualAckMode = false;
    int m_maxInflight = 0;
    bool m_isMaxInflightSet = false;
    qint64 m_ackWaitMs = 0;
    bool m_isAckWaitMsSet = false;
};

} // namespace MessageQueue

Q_DECLARE_METATYPE(QSharedPointer<stanSubscription *>)
