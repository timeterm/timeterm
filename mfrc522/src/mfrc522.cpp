#include "mfrc522/mfrc522.h"
#include <iostream>

// Physical pin 22, see the 'gpio readall' utility from wiringPi (2.52 required for RPi 4B!).
// Also very useful: https://www.ics.com/blog/gpio-programming-using-sysfs-interface
const uint8_t RESET_PIN = 25;

namespace Mfrc522 {

// Ok.
Device::Device(std::initializer_list<Spi::DeviceOpenOption> options)
    : m_spiDevice(options)
{
    Gpio::exportPin(RESET_PIN, Gpio::PinDirection::Out);
    Gpio::writePin(RESET_PIN, 1);

    init();
}

// Ok.
void Device::reset()
{
    writePcdCommand(PcdCommand::ResetPhase);
}

// Ok.
std::tuple<Status, std::vector<uint8_t>> Device::antiColl()
{
    write(Register::BitFraming, 0x00);

    std::vector<uint8_t> reqData = {
        static_cast<uint8_t>(PiccCommand::Anticoll),
        0x20,
    };
    auto [status, backData, backBits] = toCard(PcdCommand::Transceive, reqData);

    if (status == Status::Ok) {
        if (backData.size() == 5) {
            uint8_t serNumCheck = 0;
            for (int i = 0; i < 4; i++) {
                serNumCheck = serNumCheck ^ backData[i];
            }
            if (serNumCheck != backData[4]) {
                status = Status::Err;
            }
        } else {
            status = Status::Err;
        }
    }

    return std::make_tuple(status, backData);
}

// Ok.
std::tuple<Status, size_t> Device::request(PiccCommand reqMode)
{
    write(Register::BitFraming, 0x07);

    std::vector<uint8_t> reqData = {static_cast<uint8_t>(reqMode)};
    auto [status, _, backBits] = toCard(PcdCommand::Transceive, reqData);

    if ((status != Status::Ok) | (backBits != 0x10)) {
        status = Status::Err;
    }

    return std::make_tuple(status, backBits);
}

// Ok.
std::tuple<Status, std::vector<uint8_t>, size_t> Device::toCard(PcdCommand cmd,
                                                                const std::vector<uint8_t> &data)
{
    uint8_t irqEn = 0x00;
    uint8_t waitIRq = 0x00;
    Status status = Status::Err;
    uint8_t lastBits = 0;
    size_t backLen = 0;
    std::vector<uint8_t> backData;

    if (cmd == PcdCommand::Authent) {
        irqEn = 0x12;
        waitIRq = 0x10;
    }
    if (cmd == PcdCommand::Transceive) {
        irqEn = 0x77;
        waitIRq = 0x30;
    }

    write(Register::CommIEn, irqEn | 0x80u);
    clearBitMask(Register::CommIrq, 0x80);
    setBitMask(Register::FifoLevel, 0x80);

    writePcdCommand(PcdCommand::Idle);

    for (auto byte : data) {
        write(Register::FifoData, byte);
    }

    writePcdCommand(cmd);

    if (cmd == PcdCommand::Transceive) {
        setBitMask(Register::BitFraming, 0x80);
    }

    int i = 2000;
    uint8_t n;
    while (true) {
        n = read(Register::CommIrq);
        i--;
        if (!((i != 0) && ~(n & 0x01u) && ~static_cast<uint8_t>(n & waitIRq))) {
            break;
        }
    }

    clearBitMask(Register::BitFraming, 0x80);

    if (i != 0) {
        if ((read(Register::Error) & 0x1Bu) == 0) {
            status = Status::Ok;

            if (static_cast<uint8_t>(n & irqEn) & 0x01u) {
                status = Status::NoTagErr;
            }

            if (cmd == PcdCommand::Transceive) {
                n = read(Register::FifoLevel);
                lastBits = read(Register::Control) & 0x07u;
                if (lastBits != 0) {
                    backLen = (n - 1) * 8 + lastBits;
                } else {
                    backLen = n * 8;
                }

                if (n == 0) {
                    n = 1;
                }
                if (n > MAX_LEN) {
                    n = MAX_LEN;
                }

                for (i = 0; i < n; i++) {
                    backData.push_back(read(Register::FifoData));
                }
            }
        } else {
            status = Status::Err;
        }
    }

    return std::make_tuple(status, backData, backLen);
}

// Ok.
void Device::antennaOff()
{
    clearBitMask(Register::TxControl, 0x03);
}

// Ok.
void Device::antennaOn()
{
    auto tmp = read(Register::TxControl);
    if (~(tmp & 0x03u)) {
        setBitMask(Register::TxControl, 0x03);
    }
}

// Ok.
void Device::clearBitMask(Register reg, uint8_t mask)
{
    auto tmp = read(reg);
    write(reg, tmp & static_cast<uint8_t>(~mask));
}

// Ok.
void Device::setBitMask(Register reg, uint8_t mask)
{
    auto tmp = read(reg);
    write(reg, tmp | mask);
}

// Ok.
uint8_t Device::read(Register reg)
{
    auto addr = static_cast<uint8_t>((static_cast<uint8_t>(static_cast<uint8_t>(reg) << 1u) & 0x7Eu)
                                     | 0x80u);

    auto byte = m_spiDevice.transfer({addr, 0})[1];
    std::cout << "Got a byte: " << std::hex << std::to_string(byte) << std::dec << std::endl;

    return byte;
}

// Ok.
void Device::writePcdCommand(PcdCommand cmd)
{
    write(Register::Command, static_cast<uint8_t>(cmd));
}

// Ok.
void Device::write(Register reg, uint8_t cmd)
{
    auto addr = static_cast<uint8_t>(static_cast<uint8_t>(static_cast<uint8_t>(reg) << 1u) & 0x7Eu);

    auto _ = m_spiDevice.transfer({addr, cmd});
}

// Ok.
void Device::init()
{
    Gpio::writePin(RESET_PIN, 1);

    reset();

    write(Register::TMode, 0x8D);
    write(Register::TPrescaler, 0x3E);
    write(Register::TReloadL, 30);
    write(Register::TReloadH, 0);

    write(Register::TxAuto, 0x40);
    write(Register::Mode, 0x3D);

    antennaOn();
}

} // namespace Mfrc522