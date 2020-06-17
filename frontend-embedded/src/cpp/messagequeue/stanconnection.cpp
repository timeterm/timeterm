#include "stanconnection.h"
#include "enums.h"
#include "stancallbackhandlersingleton.h"

namespace MessageQueue
{

NatsStatus newNatsOptions(NatsOptionsScopedPointer &ptr)
{
    natsOptions *natsOpts = nullptr;
    auto s = static_cast<NatsStatus>(natsOptions_Create(&natsOpts));
    if (s == NatsStatus::Ok)
        ptr.reset(natsOpts);
    return s;
}

NatsStatus newStanConnOptions(StanConnOptionsScopedPointer &ptr)
{
    stanConnOptions *stanConnOpts = nullptr;
    auto s = static_cast<NatsStatus>(stanConnOptions_Create(&stanConnOpts));
    if (s == NatsStatus::Ok)
        ptr.reset(stanConnOpts);
    return s;
}

StanConnection::StanConnection(QObject *parent)
    : QObject(parent)
{
    updateStatus(newNatsOptions(m_natsOpts));
    if (m_lastStatus == NatsStatus::Ok)
        newStanConnOptions(m_connOpts);
}

QScopedArrayPointer<char> asUtf8(const QString &str)
{
    auto bytes = str.toUtf8();
    auto stdString = bytes.toStdString();
    auto src = stdString.c_str();

    auto dst = QScopedArrayPointer<char>(new char[strlen(src) + 1]);
    strcpy(dst.get(), src);

    // For some reason we can't move out the QScopedArrayPointer normally, so using a tiny hack.
    return QScopedArrayPointer<char>(dst.take());
}

void StanConnection::connect()
{
    auto cluster = asUtf8(m_cluster);
    auto clientId = asUtf8(m_clientId);
    stanConnection *stanConn = nullptr;

    auto s = static_cast<NatsStatus>(stanConnection_Connect(&stanConn, cluster.get(), clientId.get(), m_connOpts.get()));
    m_stanConnection.reset(stanConn);

    updateStatus(s);
}

NatsStatus StanConnection::lastStatus() const
{
    return m_lastStatus;
}

void StanConnection::updateStatus(NatsStatus s)
{
    m_lastStatus = s;
    if (s == NatsStatus::Ok)
        return;

    const char *text = natsStatus_GetText(static_cast<natsStatus>(s));
    auto statusStr = QString::fromLocal8Bit(text);
    emit errorOccurred(s, statusStr);
}

StanSubscription *StanConnection::subscribe(const QString &channel, StanSubOptions *opts)
{
    stanSubOptions *subOptions = nullptr;
    stanSubOptions_Create(&subOptions);

    auto channelCstr = asUtf8(channel);

    stanSubscription *sub = nullptr;
    stanConnection_Subscribe(&sub, m_stanConnection.get(), channelCstr.get(), StanCallbackHandlerSingleton::onMsg, nullptr, subOptions);

    auto rsub = new StanSubscription(this);
    rsub->setSubscription(sub);

    return rsub;
}

} // namespace MessageQueue
