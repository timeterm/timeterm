#include "natsoptions.h"
#include "strings.h"

#include <api/natscreds.h>

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

NatsStatus::Enum NatsOptions::configureOpts(natsOptions *pOpts)
{
    natsStatus s;

    if (m_url != "") {
        auto urlCstr = asUtf8CString(m_url);
        s = natsOptions_SetURL(pOpts, urlCstr.get());
        CHECK_NATS_STATUS(s);
    }

    auto credsFilePathCstr = asUtf8CString(createNatsCredsPath());
    s = natsOptions_SetUserCredentialsFromFiles(pOpts, credsFilePathCstr.get(), nullptr);
    CHECK_NATS_STATUS(s);

    // For JetStream compat.
    s = natsOptions_UseOldRequestStyle(pOpts, true);
    CHECK_NATS_STATUS(s);

    s = natsOptions_SetAllowReconnect(pOpts, false);
    CHECK_NATS_STATUS(s);

    s = natsOptions_SetMaxReconnect(pOpts, -1);

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
