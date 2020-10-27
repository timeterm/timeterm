#ifdef TIMTERMOS

#include "usbmount.h"
#include <sys/mount.h>

std::optional<QString> tryMountConfig()
{
    if (mount("/dev/sda1", "/mnt/config", "vfat", MS_NOATIME, nullptr)) {
        if (errno == EBUSY) {
            return "Mountpoint is busy";
        }
        return QStringLiteral("Could not mount: ") + strerror(errno);
    }

    return std::nullopt;
}

std::optional<QString> tryUnmountConfig() {
    if (umount("/mnt/config")) {
        if (errno == EBUSY) {
            return "Mountpoint is busy";
        }
        return QStringLiteral("Could not unmount: ") + strerror(errno);
    }

    return std::nullopt;
}

#endif
