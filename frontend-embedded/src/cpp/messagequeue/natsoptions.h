#pragma once

#include "enums.h"

#include <QObject>
#include <QSharedPointer>

#include <nats.h>

namespace MessageQueue
{

class NatsOptions: public QObject
{
    Q_OBJECT
    Q_PROPERTY(QString url READ url WRITE setUrl NOTIFY urlChanged)

public:
    explicit NatsOptions(QObject *parent = nullptr);

    NatsStatus::Enum build(natsOptions **ppOpts);

    [[nodiscard]] QString url() const;
    void setUrl(const QString &url);

signals:
    void urlChanged();

private:
    NatsStatus::Enum configureOpts(natsOptions *pOpts);

    QString m_url;
};

} // namespace MessageQueue
