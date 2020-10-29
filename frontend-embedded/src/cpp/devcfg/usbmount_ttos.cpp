#ifdef TIMETERMOS

#include "usbmount.h"
#include <sys/mount.h>
#include <sys/stat.h>

std::optional<QString> tryMountConfig()
{
    struct stat st = {0};
    if (stat("/mnt/config", &st) == -1)
        mkdir("/mnt/config", 0700);

    if (mount("/dev/sda1", "/mnt/config", "vfat", MS_RDONLY, nullptr)) {
        if (errno == EBUSY) {
            return "Mountpoint is busy";
        }
        return QStringLiteral("Could not mount: ") + strerror(errno);
    }

    return std::nullopt;
}

std::optional<QString> tryUnmountConfig()
{
    if (umount("/mnt/config")) {
        if (errno == EBUSY) {
            return "Mountpoint is busy";
        }
        return QStringLiteral("Could not unmount: ") + strerror(errno);
    }

    return std::nullopt;
}

#endif // TIMETERMOS
