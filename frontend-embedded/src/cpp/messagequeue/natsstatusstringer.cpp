#include "natsstatusstringer.h"

#include <QDebug>

namespace MessageQueue
{

NatsStatusStringer::NatsStatusStringer(QObject *parent)
    : QObject(parent)
{
}

QString NatsStatusStringer::stringify(int status)
{
    auto text = natsStatus_GetText(static_cast<natsStatus>(status));
    return QString::fromUtf8(text);
}

}
