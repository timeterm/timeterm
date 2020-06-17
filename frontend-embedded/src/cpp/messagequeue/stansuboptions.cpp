#include "stansuboptions.h"
#include "enums.h"

namespace MessageQueue
{

NatsStatus newStanSubOptions(StanSubOptionsScopedPointer &ptr)
{
    stanSubOptions *stanSubOpts = nullptr;
    auto s = static_cast<NatsStatus>(stanSubOptions_Create(&stanSubOpts));
    if (s == NatsStatus::Ok)
        ptr.reset(stanSubOpts);
    return s;
}

StanSubOptions::StanSubOptions(QObject *parent)
    : QObject(parent)
{
    // TODO: maybe don't ignore the error?
    newStanSubOptions(m_subOptions);
}

}