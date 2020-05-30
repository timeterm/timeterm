#pragma once

#include "gpio.h"
#include "spi.h"

namespace Mfrc522 {

const uint8_t MAX_LEN = 16;

enum class Status {
    Ok,
    NoTagErr,
    Err,
};

enum class Register : uint8_t {
    Command = 0x01,
    CommIEn = 0x02,
    DivlEn = 0x03,
    CommIrq = 0x04,
    DivIrq = 0x05,
    Error = 0x06,
    Status1 = 0x07,
    Status2 = 0x08,
    FifoData = 0x09,
    FifoLevel = 0x0A,
    WaterLevel = 0x0B,
    Control = 0x0C,
    BitFraming = 0x0D,
    Coll = 0x0E,
    Mode = 0x11,
    TxMode = 0x12,
    RxMode = 0x13,
    TxControl = 0x14,
    TxAuto = 0x15,
    TxSel = 0x16,
    RxSel = 0x17,
    RxThreshold = 0x18,
    Demod = 0x19,
    Mifare = 0x1C,
    SerialSpeed = 0x1F,
    CRCResultM = 0x21,
    CRCResultL = 0x22,
    ModWidth = 0x24,
    RFCfg = 0x26,
    GsN = 0x27,
    CWGsP = 0x28,
    ModGsP = 0x29,
    TMode = 0x2A,
    TPrescaler = 0x2B,
    TReloadH = 0x2C,
    TReloadL = 0x2D,
    TCounterValueH = 0x2E,
    TCounterValueL = 0x2F,
    TestSel1 = 0x31,
    TestSel2 = 0x32,
    TestPinEn = 0x33,
    TestPinValue = 0x34,
    TestBus = 0x35,
    AutoTest = 0x36,
    Version = 0x37,
    AnalogTest = 0x38,
    TestDAC1 = 0x39,
    TestDAC2 = 0x3A,
    TestADC = 0x3B,
};

enum class PcdCommand : uint8_t  {
    Idle = 0x00,
    Authent = 0x0E,
    Receive = 0x08,
    Transmit = 0x04,
    Transceive = 0x0C,
    ResetPhase = 0x0F,
    CalcCrc = 0x03,
};

enum class PiccCommand : uint8_t  {
    ReqIdl = 0x26,
    ReqAll = 0x52,
    Anticoll = 0x93,
    SelectTag = 0x93,
    Authent1A = 0x60,
    Authent1B = 0x61,
    Read = 0x30,
    Write = 0xA0,
    Decrement = 0xC0,
    Increment = 0xC1,
    Restore = 0xC2,
    Transfer = 0xB0,
    Halt = 0x50,
};

class Device
{
public:
    Device(std::initializer_list<Spi::DeviceOpenOption> options = {
               Spi::withSpeed(1000000),
           });

    void init();
    void write(Register reg, uint8_t cmd);
    void writePcdCommand(PcdCommand cmd);
    uint8_t read(Register reg);
    void setBitMask(Register reg, uint8_t mask);
    void clearBitMask(Register reg, uint8_t mask);
    void antennaOn();
    void antennaOff();
    std::tuple<Status, std::vector<uint8_t>, size_t> toCard(PcdCommand cmd,
                                                            const std::vector<uint8_t> &data);
    std::tuple<Status, size_t> request(PiccCommand reqMode);
    std::tuple<Status, std::vector<uint8_t>> antiColl();
    void reset();

private:
    Spi::Device m_spiDevice;
};

} // namespace Mfrc522
