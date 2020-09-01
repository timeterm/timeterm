#pragma once

#include <QObject>

namespace MessageQueue
{

class DisownTokenMessage
{
    Q_GADGET
    Q_PROPERTY(QString deviceId READ deviceId)
    Q_PROPERTY(QString tokenHash READ tokenHash)
    Q_PROPERTY(QString tokenHashAlg READ tokenHashAlg)

public:
    void setDeviceId(const QString &deviceId);
    [[nodiscard]] QString deviceId() const;
    void setTokenHash(const QString &tokenHash);
    [[nodiscard]] QString tokenHash() const;
    void setTokenHashAlg(const QString &tokenHashAlg);
    [[nodiscard]] QString tokenHashAlg() const;

private:
    QString m_deviceId;
    QString m_tokenHash;
    QString m_tokenHashAlg;
};

}

Q_DECLARE_METATYPE(MessageQueue::DisownTokenMessage)
