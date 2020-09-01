#pragma once

#include "enums.h"

#include <nats.h>

#include <QObject>
#include <QSharedPointer>

namespace MessageQueue
{

class NatsOptions: public QObject
{
    Q_OBJECT
    Q_PROPERTY(MessageQueue::NatsStatus::Enum lastStatus READ lastStatus)

public:
    explicit NatsOptions(QObject *parent = nullptr);

    QSharedPointer<natsOptions> options();

    [[nodiscard]] NatsStatus::Enum lastStatus() const;

signals:
    void errorOccurred(NatsStatus::Enum status, const QString &message);

private:
    void updateStatus(NatsStatus::Enum s);

    QSharedPointer<natsOptions> m_options;
    NatsStatus::Enum m_lastStatus = NatsStatus::Enum::Ok;
};

} // namespace MessageQueue
