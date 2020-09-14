#include "natsoptions.h"
#include "strings.h"

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
    natsStatus s = NATS_OK;

    if (m_url != "") {
        auto urlCstr = asUtf8CString(m_url);
        s = natsOptions_SetURL(pOpts, urlCstr.get());
        CHECK_NATS_STATUS(s);
    }

    s = natsOptions_UseOldRequestStyle(pOpts, true);
    CHECK_NATS_STATUS(s);

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
