#include "disowntokenmessage.h"

namespace MessageQueue
{

void DisownTokenMessage::setDeviceId(const QString &deviceId)
{
    m_deviceId = deviceId;
}

QString DisownTokenMessage::deviceId() const
{
    return m_deviceId;
}

void DisownTokenMessage::setTokenHash(const QString &tokenHash)
{
    m_tokenHash = tokenHash;
}

QString DisownTokenMessage::tokenHash() const
{
    return m_tokenHash;
}

void DisownTokenMessage::setTokenHashAlg(const QString &tokenHashAlg)
{
    m_tokenHashAlg = tokenHashAlg;
}

QString DisownTokenMessage::tokenHashAlg() const
{
    return m_tokenHashAlg;
}

}