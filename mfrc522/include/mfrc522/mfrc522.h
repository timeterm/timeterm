#pragma once

#include "spi.h"

namespace Mfrc522 {

class Device
{
    Device(std::initializer_list<Spi::DeviceOpenOption> options)
    {
        m_spiDevice = Spi::Device(options);

        wi
    }

private:
    Spi::Device m_spiDevice;
};

} // namespace Mfrc522
