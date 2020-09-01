#pragma once

#include <QObject>

namespace MessageQueue
{

class RetrieveNewTokenMessage
{
    Q_GADGET
    Q_PROPERTY(QString deviceId READ deviceId)
    Q_PROPERTY(QString currentTokenHash READ currentTokenHash)
    Q_PROPERTY(QString currentTokenhashAlg READ currentTokenHashAlg)

public:
    void setDeviceId(const QString &deviceId);
    [[nodiscard]] QString deviceId() const;
    void setCurrentTokenHash(const QString &currentTokenHash);
    [[nodiscard]] QString currentTokenHash() const;
    void setCurrentTokenHashAlg(const QString &currentTokenHashAlg);
    [[nodiscard]] QString currentTokenHashAlg() const;

private:
    QString m_deviceId;
    QString m_currentTokenHash;
    QString m_currentTokenHashAlg;
};

}

Q_DECLARE_METATYPE(MessageQueue::RetrieveNewTokenMessage)
