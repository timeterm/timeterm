#include <mfrc522/spi.h>

#include <cstdint>
#include <fcntl.h>
#include <iostream>
#include <linux/spi/spidev.h>
#include <linux/types.h>
#include <sys/ioctl.h>
#include <unistd.h>

namespace Mfrc522::Spi {

DeviceOpenOption withMode(uint8_t mode)
{
    return [=](DeviceOpenOptions &opts) { opts.mode = mode; };
}

DeviceOpenOption withDevicePath(const std::string &device)
{
    return [=](DeviceOpenOptions &opts) { opts.device = device; };
}

DeviceOpenOption withSpeed(uint32_t hz)
{
    return [=](DeviceOpenOptions &opts) { opts.speed = hz; };
}

DeviceOpenOption withBits(uint8_t bits)
{
    return [=](DeviceOpenOptions &opts) { opts.bits = bits; };
}

DeviceOpenOption withDelay(uint16_t delay)
{
    return [=](DeviceOpenOptions &opts) { opts.delay = delay; };
}

#pragma clang diagnostic push
#pragma ide diagnostic ignored "hicpp-signed-bitwise"
Device::Device(std::initializer_list<DeviceOpenOption> opts)
{
    for (auto &opt : opts) {
        opt(m_options);
    }

    int fd = open(m_options.device.c_str(), O_RDWR);
    if (fd == -1) {
        throw DeviceOpenException(errno);
    }

    int ret = ioctl(fd, SPI_IOC_WR_MODE, &m_options.mode);
    if (ret == -1) {
        throw DeviceConfigureException("can't write SPI mode", errno);
    }

    ret = ioctl(fd, SPI_IOC_RD_MODE, &m_options.mode);
    if (ret == -1) {
        throw DeviceConfigureException("can't read SPI mode", errno);
    }

    ret = ioctl(fd, SPI_IOC_WR_BITS_PER_WORD, &m_options.bits);
    if (ret == -1) {
        throw DeviceConfigureException("can't write bits per word", errno);
    }

    ret = ioctl(fd, SPI_IOC_RD_BITS_PER_WORD, &m_options.bits);
    if (ret == -1) {
        throw DeviceConfigureException("can't read bits per word", errno);
    }

    ret = ioctl(fd, SPI_IOC_WR_MAX_SPEED_HZ, &m_options.speed);
    if (ret == -1) {
        throw DeviceConfigureException("can't write max speed (in Hz)", errno);
    }

    ret = ioctl(fd, SPI_IOC_RD_MAX_SPEED_HZ, &m_options.speed);
    if (ret == -1) {
        throw DeviceConfigureException("can't read max speed (in Hz)", errno);
    }

    m_fd = fd;
}
#pragma clang diagnostic pop

Device::~Device()
{
    close(m_fd);
}

#pragma clang diagnostic push
#pragma ide diagnostic ignored "hicpp-signed-bitwise"
std::vector<uint8_t> Device::transfer(const std::vector<uint8_t> &tx)
{
    if (tx.size() > UINT32_MAX) {
        throw PayloadTooLargeException();
    }

    std::vector<uint8_t> rx;
    // Make rx the size of tx, because we're reading as much bytes
    // as we're writing (apparently).
    rx.resize(tx.size());

    struct spi_ioc_transfer transfer = {
        .tx_buf = (uintptr_t) tx.data(),
        .rx_buf = (uintptr_t) rx.data(),
        .len = (uint32_t) rx.size(),
        .speed_hz = m_options.speed,
        .delay_usecs = m_options.delay,
        .bits_per_word = m_options.bits,
        .cs_change = 0,
    };

    int ret = ioctl(m_fd, SPI_IOC_MESSAGE(1), &transfer);
    if (ret < 1) {
        throw SpiSendMessageException(errno);
    }

    return rx;
}
#pragma clang diagnostic pop

DeviceOpenException::DeviceOpenException(int err)
    : std::runtime_error(std::string("could not open device: ") + strerror(err))
    , m_errno(err)
{}

DeviceConfigureException::DeviceConfigureException(const std::string &msg, int err)
    : std::runtime_error("could not configure device: " + msg + ": " + strerror(err))
    , m_errno(err)
{}

PayloadTooLargeException::PayloadTooLargeException()
    : std::runtime_error("payload is too large (max is UINT32_MAX)")
{}

SpiSendMessageException::SpiSendMessageException(int err)
    : std::runtime_error(std::string("could not send SPI message: ") + strerror(err))
    , m_errno(err)
{}

} // namespace Mfrc522::Spi