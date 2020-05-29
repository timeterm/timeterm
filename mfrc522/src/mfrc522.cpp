#include "mfrc522/mfrc522.h"

const uint8_t RESET_PIN = 22;

Mfrc522::Device::Device(std::initializer_list<Spi::DeviceOpenOption> options)
{
    m_spiDevice = Spi::Device(options);

    Gpio::exportPin(RESET_PIN, Gpio::PinDirection::Out);
    Gpio::writePin(RESET_PIN, 1);
}
