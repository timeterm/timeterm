#include "usbmount.h"

#ifndef TIMETERMOS

std::optional<QString> tryMountConfig()
{
    return std::nullopt;
}

std::optional<QString> tryUnmountConfig()
{
    return std::nullopt;
}

#endif
