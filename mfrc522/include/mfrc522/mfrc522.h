#pragma once

#include "gpio.h"
#include "spi.h"

namespace Mfrc522 {

class Device
{
public:
    Device(std::initializer_list<Spi::DeviceOpenOption> options);

private:
    Spi::Device m_spiDevice;
};

} // namespace Mfrc522
