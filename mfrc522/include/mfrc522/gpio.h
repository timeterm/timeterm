#pragma once

#include <cstring>
#include <fcntl.h>
#include <netinet/in.h>
#include <stdexcept>
#include <sys/epoll.h>
#include <unistd.h>
#include <sys/mman.h>

namespace Mfrc522::Gpio {

enum class Mode {
    Unknown = -1,
    Board = 10,
    Bcm = 11,
};

struct RpiInfo
{
    int p1Revision = -1;
    std::string ram;
    std::string manufacturer;
    std::string processor;
    std::string type;
    char revision[1024] = {0};
};

struct Callback
{
    unsigned int gpio;
    void (*func)(unsigned int gpio);
    Callback *next;
};
Callback *callbacks = nullptr;

int threadRunning = 0;
int epfdThread = -1;
int epfdBlocking = -1;
int eventOccurred[54] = {0};
const char *stredge[4] = {"none", "rising", "falling", "both"};

#define BLOCK_SIZE (4 * 1024)

static volatile uint32_t *gpioMap;

struct Gpios
{
    uint32_t gpio;
    int valueFd;
    int exported;
    int edge;
    int initialThread;
    int initialWait;
    int threadAdded;
    int bounceTime;
    uint64_t lastCall;
    Gpios *next;
};

Gpios *gpioList = nullptr;

int gpioDirection[54];

Mode mode = Mode::Unknown;
RpiInfo rpiInfo;

const int pinToGpioRev1[41] = {0};
const int pinToGpioRev2[41] = {0};
const int pinToGpioRev3[41] = {0};
const int (*pinToGpio)[41];

void getRpiInfo(RpiInfo *info)
{
    FILE *fp;
    char buffer[1024];
    char hardware[1024];
    char revision[1024];
    int found = 0;
    int len;

    if ((fp = fopen("/proc/device-tree/system/linux,revision", "r"))) {
        uint32_t n;
        if (fread(&n, sizeof(n), 1, fp) != 1) {
            fclose(fp);
            throw std::runtime_error("could not read from file");
        }
        sprintf(revision, "%x", ntohl(n));
        found = 1;
    } else if ((fp = fopen("/proc/cpuinfo", "r"))) {
        while (!feof(fp) && fgets(buffer, sizeof(buffer), fp)) {
            sscanf(buffer, "Hardware	: %s", hardware);
            if (strcmp(hardware, "BCM2708") == 0 || strcmp(hardware, "BCM2709") == 0
                || strcmp(hardware, "BCM2835") == 0 || strcmp(hardware, "BCM2836") == 0
                || strcmp(hardware, "BCM2837") == 0) {
                found = 1;
            }
            sscanf(buffer, "Revision	: %s", revision);
        }
    } else
        throw std::runtime_error("could not open device information");
    fclose(fp);

    if (!found)
        throw std::runtime_error("could not find device information");

    if ((len = strlen(revision)) == 0)
        throw std::runtime_error("invalid device revision");

    if (len >= 6 && strtol((char[]){revision[len - 6], 0}, nullptr, 16) & 8) {
        // new scheme
        //info->rev = revision[len-1]-'0';
        strcpy(info->revision, revision);
        switch (revision[len - 3]) {
        case '0':
            switch (revision[len - 2]) {
            case '0':
                info->type = "Model A";
                info->p1Revision = 2;
                break;
            case '1':
                info->type = "Model B";
                info->p1Revision = 2;
                break;
            case '2':
                info->type = "Model A+";
                info->p1Revision = 3;
                break;
            case '3':
                info->type = "Model B+";
                info->p1Revision = 3;
                break;
            case '4':
                info->type = "Pi 2 Model B";
                info->p1Revision = 3;
                break;
            case '5':
                info->type = "Alpha";
                info->p1Revision = 3;
                break;
            case '6':
                info->type = "Compute Module 1";
                info->p1Revision = 0;
                break;
            case '8':
                info->type = "Pi 3 Model B";
                info->p1Revision = 3;
                break;
            case '9':
                info->type = "Zero";
                info->p1Revision = 3;
                break;
            case 'a':
                info->type = "Compute Module 3";
                info->p1Revision = 0;
                break;
            case 'c':
                info->type = "Zero W";
                info->p1Revision = 3;
                break;
            case 'd':
                info->type = "Pi 3 Model B+";
                info->p1Revision = 3;
                break;
            case 'e':
                info->type = "Pi 3 Model A+";
                info->p1Revision = 3;
                break;
            default:
                info->type = "Unknown";
                info->p1Revision = 3;
                break;
            }
            break;
        case '1':
            switch (revision[len - 2]) {
            case '0':
                info->type = "Compute Module 3+";
                info->p1Revision = 0;
                break;
            case '1':
                info->type = "Pi 4 Model B";
                info->p1Revision = 3;
                break;
            default:
                info->type = "Unknown";
                info->p1Revision = 3;
                break;
            }
            break;
        default:
            info->type = "Unknown";
            info->p1Revision = 3;
            break;
        }

        switch (revision[len - 4]) {
        case '0':
            info->processor = "BCM2835";
            break;
        case '1':
            info->processor = "BCM2836";
            break;
        case '2':
            info->processor = "BCM2837";
            break;
        case '3':
            info->processor = "BCM2711";
            break;
        default:
            info->processor = "Unknown";
            break;
        }
        switch (revision[len - 5]) {
        case '0':
            info->manufacturer = "Sony";
            break;
        case '1':
            info->manufacturer = "Egoman";
            break;
        case '2':
            info->manufacturer = "Embest";
            break;
        case '3':
            info->manufacturer = "Sony Japan";
            break;
        case '4':
            info->manufacturer = "Embest";
            break;
        case '5':
            info->manufacturer = "Stadium";
            break;
        default:
            info->manufacturer = "Unknown";
            break;
        }
        switch (strtol((char[]){revision[len - 6], 0}, nullptr, 16) & 7) {
        case 0:
            info->ram = "256M";
            break;
        case 1:
            info->ram = "512M";
            break;
        case 2:
            info->ram = "1G";
            break;
        case 3:
            info->ram = "2G";
            break;
        case 4:
            info->ram = "4G";
            break;
        default:
            info->ram = "Unknown";
            break;
        }
    } else {
        // old scheme
        info->ram = "Unknown";
        info->manufacturer = "Unknown";
        info->processor = "Unknown";
        info->type = "Unknown";
        strcpy(info->revision, revision);

        uint64_t rev;
        sscanf(revision, "%llx", &rev);
        rev = rev & 0xefffffff; // ignore preceeding 1000 for overvolt

        if (rev == 0x0002 || rev == 0x0003) {
            info->type = "Model B";
            info->p1Revision = 1;
            info->ram = "256M";
            info->manufacturer = "Egoman";
            info->processor = "BCM2835";
        } else if (rev == 0x0004) {
            info->type = "Model B";
            info->p1Revision = 2;
            info->ram = "256M";
            info->manufacturer = "Sony UK";
            info->processor = "BCM2835";
        } else if (rev == 0x0005) {
            info->type = "Model B";
            info->p1Revision = 2;
            info->ram = "256M";
            info->manufacturer = "Qisda";
            info->processor = "BCM2835";
        } else if (rev == 0x0006) {
            info->type = "Model B";
            info->p1Revision = 2;
            info->ram = "256M";
            info->manufacturer = "Egoman";
            info->processor = "BCM2835";
        } else if (rev == 0x0007) {
            info->type = "Model A";
            info->p1Revision = 2;
            info->ram = "256M";
            info->manufacturer = "Egoman";
            info->processor = "BCM2835";
        } else if (rev == 0x0008) {
            info->type = "Model A";
            info->p1Revision = 2;
            info->ram = "256M";
            info->manufacturer = "Sony UK";
            info->processor = "BCM2835";
        } else if (rev == 0x0009) {
            info->type = "Model A";
            info->p1Revision = 2;
            info->ram = "256M";
            info->manufacturer = "Qisda";
            info->processor = "BCM2835";
        } else if (rev == 0x000d) {
            info->type = "Model B";
            info->p1Revision = 2;
            info->ram = "512M";
            info->manufacturer = "Egoman";
            info->processor = "BCM2835";
        } else if (rev == 0x000e) {
            info->type = "Model B";
            info->p1Revision = 2;
            info->ram = "512M";
            info->manufacturer = "Sony UK";
            info->processor = "BCM2835";
        } else if (rev == 0x000f) {
            info->type = "Model B";
            info->p1Revision = 2;
            info->ram = "512M";
            info->manufacturer = "Qisda";
            info->processor = "BCM2835";
        } else if (rev == 0x0010) {
            info->type = "Model B+";
            info->p1Revision = 3;
            info->ram = "512M";
            info->manufacturer = "Sony UK";
            info->processor = "BCM2835";
        } else if (rev == 0x0011) {
            info->type = "Compute Module 1";
            info->p1Revision = 0;
            info->ram = "512M";
            info->manufacturer = "Sony UK";
            info->processor = "BCM2835";
        } else if (rev == 0x0012) {
            info->type = "Model A+";
            info->p1Revision = 3;
            info->ram = "256M";
            info->manufacturer = "Sony UK";
            info->processor = "BCM2835";
        } else if (rev == 0x0013) {
            info->type = "Model B+";
            info->p1Revision = 3;
            info->ram = "512M";
            info->manufacturer = "Embest";
            info->processor = "BCM2835";
        } else if (rev == 0x0014) {
            info->type = "Compute Module 1";
            info->p1Revision = 0;
            info->ram = "512M";
            info->manufacturer = "Embest";
            info->processor = "BCM2835";
        } else if (rev == 0x0015) {
            info->type = "Model A+";
            info->p1Revision = 3;
            info->ram = "Unknown";
            info->manufacturer = "Embest";
            info->processor = "BCM2835";
        } else { // don't know - assume revision 3 p1 connector
            info->p1Revision = 3;
        }
    }
}

// TODO: run cleanup at application shutdown (adapt from cleanup, event_cleanup_all)
void init()
{
    for (int i = 0; i < 64; i++) {
        gpioDirection[i] = -1;
    }

    getRpiInfo(&rpiInfo);

    switch (rpiInfo.p1Revision) {
    case 1:
        pinToGpio = &pinToGpioRev1;
        break;
    case 2:
        pinToGpio = &pinToGpioRev2;
        break;
    default:
        pinToGpio = &pinToGpioRev3;
    }
}

void setMode(Mode newMode)
{
    if (mode != Mode::Unknown && newMode != mode) {
        throw std::runtime_error("a different mode has already been set");
    }

    if (newMode != Mode::Board && newMode != Mode::Bcm) {
        throw std::runtime_error("invalid mode");
    }

    if (rpiInfo.p1Revision == 0 && newMode == Mode::Board) {
        throw std::runtime_error("Board numbering system is not applicable on the compute module");
    }

    mode = newMode;
}

void deleteGpio(unsigned int gpio)
{
    Gpios *g = gpioList;
    Gpios *prev = nullptr;

    while (g != nullptr) {
        if (g->gpio == gpio) {
            if (prev == nullptr)
                gpioList = g->next;
            else
                prev->next = g->next;
            free(g);
            return;
        } else {
            prev = g;
            g = g->next;
        }
    }
}

Gpios *getGpio(unsigned int gpio)
{
    Gpios *g = gpioList;
    while (g != nullptr) {
        if (g->gpio == gpio)
            return g;
        g = g->next;
    }
    return nullptr;
}

void removeCallbacks(unsigned int gpio)
{
    Callback *cb = callbacks;
    Callback *temp;
    Callback *prev = nullptr;

    while (cb != nullptr) {
        if (cb->gpio == gpio) {
            if (prev == nullptr)
                callbacks = cb->next;
            else
                prev->next = cb->next;
            temp = cb;
            cb = cb->next;
            free(temp);
        } else {
            prev = cb;
            cb = cb->next;
        }
    }
}

#define x_write(fd, buf, len) \
    do { \
        size_t x_write_len = (len); \
\
        if ((size_t) write((fd), (buf), x_write_len) != x_write_len) { \
            close(fd); \
            return (-1); \
        } \
    } while (/* CONSTCOND */ 0)

int gpioUnexport(unsigned int gpio)
{
    int fd, len;
    char str_gpio[3];

    if ((fd = open("/sys/class/gpio/unexport", O_WRONLY)) < 0)
        return -1;

    len = snprintf(str_gpio, sizeof(str_gpio), "%d", gpio);
    x_write(fd, str_gpio, len);
    close(fd);

    return 0;
}

int gpioSetEdge(unsigned int gpio, unsigned int edge)
{
    int fd;
    char filename[28];

    snprintf(filename, sizeof(filename), "/sys/class/gpio/gpio%d/edge", gpio);

    if ((fd = open(filename, O_WRONLY)) < 0)
        return -1;

    x_write(fd, stredge[edge], strlen(stredge[edge]) + 1);
    close(fd);
    return 0;
}

#define NO_EDGE 0

void removeEdgeDetect(unsigned int gpio)
{
    epoll_event ev{};
    Gpios *g = getGpio(gpio);

    if (g == nullptr)
        return;

    // delete epoll of fd

    ev.events = EPOLLIN | EPOLLET | EPOLLPRI;
    ev.data.fd = g->valueFd;
    epoll_ctl(epfdThread, EPOLL_CTL_DEL, g->valueFd, &ev);

    // delete callbacks for gpio
    removeCallbacks(gpio);

    gpioSetEdge(gpio, NO_EDGE);
    g->edge = NO_EDGE;

    if (g->valueFd != -1)
        close(g->valueFd);

    gpioUnexport(gpio);
    eventOccurred[gpio] = 0;

    deleteGpio(gpio);
}

void eventCleanup(int gpio)
{
    Gpios *g = gpioList;
    Gpios *next_gpio = nullptr;

    while (g != nullptr) {
        next_gpio = g->next;
        if ((gpio == -666) || ((int) g->gpio == gpio))
            removeEdgeDetect(g->gpio);
        g = next_gpio;
    }
    if (gpioList == nullptr) {
        if (epfdBlocking != -1) {
            close(epfdBlocking);
            epfdBlocking = -1;
        }
        if (epfdThread != -1) {
            close(epfdThread);
            epfdThread = -1;
        }
        threadRunning = 0;
    }
}

void cleanup()
{
    munmap((void *) gpioMap, BLOCK_SIZE);
}

void eventCleanupAll()
{
    eventCleanup(-666);
}

} // namespace Mfrc522::Gpio