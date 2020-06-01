#pragma once

#include <cstring>
#include <functional>
#include <linux/spi/spidev.h>
#include <string>

//! The Mfrc522 namespace.
namespace Mfrc522 { // NOLINT

//! The Spi namespace.
namespace Spi {}

} // namespace Mfrc522

namespace Mfrc522::Spi {

//! DeviceOpenOptions contains options for opening a device.
struct DeviceOpenOptions
{
    //! The path to the device to open.
    //! Can be set using withDevicePath()
    std::string device = "/dev/spidev0.0";

    //! `SPI_MODE_0`, `SPI_MODE_1`, `SPI_MODE_3` or `SPI_MODE_4`.
    //! Can be set using withMode()
    //! \sa https://www.kernel.org/doc/Documentation/spi/spidev
    uint8_t mode = SPI_MODE_0;

    //! Speed in Hz (Hertz).
    //! Can be set using withSpeed()
    uint32_t speed = 500000;

    //! Amount of bits per word.
    //! Can be set using withBits()
    uint8_t bits = 8;

    //! Delay in µs (microseconds).
    //! Can be set using withDelay()
    uint16_t delay = 0;
};

//! DeviceOpenOption is an option for the Device constructor Device::Device().
using DeviceOpenOption = std::function<void(DeviceOpenOptions &)>;

//! withMode sets the SPI mode.
//! See DeviceOpenOptions::mode for options.
DeviceOpenOption withMode(uint8_t mode);

//! withDevicePath sets the SPI device to use.
DeviceOpenOption withDevicePath(const std::string &device);

//! withSpeed sets the speed (in Hz) for communicating with the SPI device.
DeviceOpenOption withSpeed(uint32_t hz);

//! withBits sets DeviceOpenOptions::bits.
DeviceOpenOption withBits(uint8_t bits);

//! withDelay sets DeviceOpenOptions::delay.
DeviceOpenOption withDelay(uint16_t delay);

//! Device is a handle for a SPI device.
class Device
{
public:
    //! The constructor for Device.
    Device(std::initializer_list<DeviceOpenOption> opts = {});

    ~Device();

    [[nodiscard]] std::vector<uint8_t> transfer(const std::vector<uint8_t> &bytes) const;
    int transferN(const char *buf, uint32_t len, const char *rx = nullptr) const;
    [[nodiscard]] uint8_t transfer1(uint8_t byte) const;

private:
    DeviceOpenOptions m_options;
    int m_fd = 0;
};

class DeviceOpenException : public std::runtime_error
{
public:
    explicit DeviceOpenException(int err);

private:
    int m_errno;
};

class DeviceConfigureException : public std::runtime_error
{
public:
    explicit DeviceConfigureException(const std::string &msg, int err);

private:
    int m_errno;
};

class PayloadTooLargeException : public std::runtime_error
{
public:
    PayloadTooLargeException();
};

class SpiSendMessageException : public std::runtime_error
{
public:
    explicit SpiSendMessageException(int err);

private:
    int m_errno;
};

} // namespace Mfrc522::Spi
