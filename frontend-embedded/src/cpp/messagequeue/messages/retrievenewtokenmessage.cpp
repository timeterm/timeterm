#include "retrievenewtokenmessage.h"

namespace MessageQueue
{

void RetrieveNewTokenMessage::setDeviceId(const QString &deviceId)
{
    m_deviceId = deviceId;
}

QString RetrieveNewTokenMessage::deviceId() const
{
    return m_deviceId;
}

void RetrieveNewTokenMessage::setCurrentTokenHash(const QString &currentTokenHash)
{
    m_currentTokenHash = currentTokenHash;
}

QString RetrieveNewTokenMessage::currentTokenHash() const
{
    return m_currentTokenHash;
}

void RetrieveNewTokenMessage::setCurrentTokenHashAlg(const QString &currentTokenHashAlg)
{
    m_currentTokenHashAlg = currentTokenHashAlg;
}

QString RetrieveNewTokenMessage::currentTokenHashAlg() const
{
    return m_currentTokenHashAlg;
}

} // namespace MessageQueue