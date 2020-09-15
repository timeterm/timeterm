#include "natsoptions.h"
#include "strings.h"

#include <QDebug>

namespace MessageQueue
{

MessageQueue::NatsOptions::NatsOptions(QObject *parent)
    : QObject(parent)
{
}

#define CHECK_NATS_STATUS(status)                                    \
    do {                                                             \
        if ((status) != NATS_OK) return NatsStatus::fromC((status)); \
    } while (0)

NatsStatus::Enum NatsOptions::build(natsOptions **ppOpts)
{
    auto status = natsOptions_Create(ppOpts);
    CHECK_NATS_STATUS(status);

    auto s = configureOpts(*ppOpts);
    if (s != NatsStatus::Enum::Ok) {
        natsOptions_Destroy(*ppOpts);
        return s;
    }

    return s;
}

void natsConnectionLostCb(natsConnection *, void *)
{
    qCritical() << "NATS connection lost";
}

NatsStatus::Enum NatsOptions::configureOpts(natsOptions *pOpts)
{
    natsStatus s;

    if (m_url != "") {
        auto urlCstr = asUtf8CString(m_url);
        s = natsOptions_SetURL(pOpts, urlCstr.get());
        CHECK_NATS_STATUS(s);
    }

    s = natsOptions_UseOldRequestStyle(pOpts, true);
    CHECK_NATS_STATUS(s);

    s = natsOptions_SetAllowReconnect(pOpts, true);
    CHECK_NATS_STATUS(s);

    s = natsOptions_SetReconnectWait(pOpts, 5000);
    CHECK_NATS_STATUS(s);

    s = natsOptions_SetRetryOnFailedConnect(pOpts, true, nullptr, nullptr);
    CHECK_NATS_STATUS(s);

    natsOptions_SetDisconnectedCB(pOpts, natsConnectionLostCb, nullptr);

    return NatsStatus::fromC(s);
}

QString NatsOptions::url() const
{
    return m_url;
}

void NatsOptions::setUrl(const QString &url)
{
    if (url != m_url) {
        m_url = url;
        emit urlChanged();
    }
}

} // namespace MessageQueue
