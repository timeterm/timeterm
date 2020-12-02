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
    Q_PROPERTY(QString credsFilePath READ credsFilePath WRITE setCredsFilePath NOTIFY credsFilePathChanged)

public:
    explicit NatsOptions(QObject *parent = nullptr);

    NatsStatus::Enum build(natsOptions **ppOpts);

    [[nodiscard]] QString url() const;
    void setUrl(const QString &url);

    [[nodiscard]] QString credsFilePath() const;
    void setCredsFilePath(const QString &path);

signals:
    void urlChanged();
    void credsFilePathChanged();

private:
    NatsStatus::Enum configureOpts(natsOptions *pOpts);

    QString m_url;
    QString m_credsFilePath;
};

} // namespace MessageQueue
