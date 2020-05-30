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

enum class Register {
    Reserved00,
    Command,
    CommIEn,
    DivlEn,
    CommIrq,
    DivIrq,
    Error,
    Status1,
    Status2,
    FifoData,
    FifoLevel,
    WaterLevel,
    Control,
    BitFraming,
    Coll,
    Reserved01,
    Reserved10,
    Mode,
    TxMode,
    RxMode,
    TxControl,
    TxAuto,
    TxSel,
    RxSel,
    RxThreshold,
    Demod,
    Reserved11,
    Reserved12,
    Mifare,
    Reserved13,
    Reserved14,
    SerialSpeed,
    Reserved20,
    CRCResultM,
    CRCResultL,
    Reserved21,
    ModWidth,
    Reserved22,
    RfCfg,
    GsNReg,
    CWGsP,
    ModGsP,
    TMode,
    TPrescaler,
    TReloadH,
    TReloadL,
    TCounterValueH,
    TCounterValueL,
    Reserved30,
    TestSel1,
    TestSel2,
    TestPinEn,
    TestPinValue,
    TestBus,
    AutoTest,
    Version,
    AnalogTest,
    TestDAC1,
    TestDAC2,
    TestADC,
    Reserved31,
    Reserved32,
    Reserved33,
    Reserved34,
};

enum class PcdCommand {
    Idle = 0x00,
    Authent = 0x0E,
    Receive = 0x08,
    Transmit = 0x04,
    Transceive = 0x0C,
    ResetPhase = 0x0F,
    CalcCrc = 0x03,
};

enum class PiccCommand {
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
