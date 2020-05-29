#pragma once

#include <cstring>
#include <functional>
#include <string>

/*!
 * Based on https://github.com/lthiery/SPI-Py,
 * the non-working library https://github.com/mxgxw/MFRC522-python,
 * and its working fork https://github.com/ukleinek/MFRC522-python.
 *
 * Parts licensed under GPL-2.0 and LGPL-3.0 because the Python library are
 * that as well, and that's only nice to do.
 */

//! The Spi namespace.
namespace Spi {

//! DeviceOpenOptions contains options for opening a device.
struct DeviceOpenOptions
{
    //! The path to the device to open.
    //! Can be set using withDevicePath()
    std::string device = "/dev/spidev0.0";

    //! `SPI_MODE_0`, `SPI_MODE_1`, `SPI_MODE_3` or `SPI_MODE_4`.
    //! Can be set using withMode()
    //! \sa https://www.kernel.org/doc/Documentation/spi/spidev
    uint8_t mode = 0;

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

    void transfer();

private:
    DeviceOpenOptions m_options;
    std::optional<int> m_fd = std::nullopt;
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
    std::string m_what;
};

} // namespace Spi
