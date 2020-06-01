#include <cstdint>
#include <cstdio>
#include <cstring>
#include <mfrc522/gpio.h>
#include <mfrc522/mfrc522.h>
#include <string>

#define RSTPIN 25
#define LOW 0
#define HIGH 1

using namespace std;

namespace Mfrc522 {

Device::Device()
    : m_spiDev({Spi::withSpeed(4000000)})
{
    Gpio::exportPin(RSTPIN, Gpio::PinDirection::Out);
    Gpio::writePin(RSTPIN, LOW);
}

//-----------------------------------------------------------------------------------
// Basic interface functions for communicating with the MFRC522
//-----------------------------------------------------------------------------------

void Device::PCD_WriteRegister(uint8_t reg, uint8_t value) const
{
    uint8_t data[2];
    data[0] = reg & 0x7Eu;
    data[1] = value;
    m_spiDev.transferNU(data, 2);
}

void Device::PCD_WriteRegister(uint8_t reg, uint8_t count, uint8_t *values) const
{
    for (uint8_t index = 0; index < count; index++) {
        PCD_WriteRegister(reg, values[index]);
    }
}

uint8_t Device::PCD_ReadRegister(uint8_t reg) const
{
    uint8_t data[2];
    data[0] = 0x80u | (reg & 0x7Eu);
    m_spiDev.transferNU(data, 2);
    return (uint8_t) data[1];
}

void Device::PCD_ReadRegister(uint8_t reg, uint8_t count, uint8_t *values, uint8_t rxAlign) const
{
    if (count == 0) {
        return;
    }

    uint8_t address
        = 0x80u
          | (reg
             & 0x7Eu); // MSB == 1 is for reading. LSB is not used in address. Datasheet section 8.1.2.3.
    uint8_t index = 0; // Index in values array.
    count--;           // One read is performed outside of the loop
    m_spiDev.transfer1(address);
    while (index < count) {
        if (index == 0 && rxAlign) { // Only update bit positions rxAlign..7 in values[0]
            // Create bit mask for bit positions rxAlign..7
            uint8_t mask = 0;
            for (uint8_t i = rxAlign; i <= 7; i++) {
                mask |= 1u << i;
            }
            // Read value and tell that we want to read the same address again.
            uint8_t value = m_spiDev.transfer1(address);
            // Apply mask to both current value of values[0] and the new data in value.
            values[0] = static_cast<uint8_t>(values[index] & static_cast<uint8_t>(~mask))
                        | static_cast<uint8_t>(value & mask);
        } else { // Normal case
            values[index] = m_spiDev.transfer1(address);
        }
        index++;
    }
    values[index] = m_spiDev.transfer1(address); // Read the final uint8_t. Send 0 to stop reading.
}

void Device::PCD_SetRegisterBitMask(uint8_t reg, uint8_t mask) const
{
    uint8_t tmp;
    tmp = PCD_ReadRegister(reg);
    PCD_WriteRegister(reg, tmp | mask); // set bit mask
}

void Device::PCD_ClearRegisterBitMask(uint8_t reg, uint8_t mask) const
{
    uint8_t tmp;
    tmp = PCD_ReadRegister(reg);
    PCD_WriteRegister(reg, tmp & static_cast<uint8_t>(~mask)); // clear bit mask
}

uint8_t Device::PCD_CalculateCRC(uint8_t *data, uint8_t length, uint8_t *result) const
{
    PCD_WriteRegister(CommandReg, PcdIdle);      // Stop any active command.
    PCD_WriteRegister(DivIrqReg, 0x04);           // Clear the CRCIRq interrupt request bit
    PCD_SetRegisterBitMask(FIFOLevelReg, 0x80);   // FlushBuffer = 1, FIFO initialization
    PCD_WriteRegister(FIFODataReg, length, data); // Write data to the FIFO
    PCD_WriteRegister(CommandReg, PcdCalcCrc);   // Start the calculation

    // Wait for the CRC calculation to complete. Each iteration of the while-loop takes 17.73�s.
    uint16_t i = 5000;
    uint8_t n;
    while (true) {
        n = PCD_ReadRegister(
            DivIrqReg); // DivIrqReg[7..0] bits are: Set2 reserved reserved MfinActIRq reserved CRCIRq reserved reserved
        if (n & 0x04u) { // CRCIRq bit set - calculation done
            break;
        }
        if (--i
            == 0) { // The emergency break. We will eventually terminate on this one after 89ms. Communication with the MFRC522 might be down.
            return STATUS_TIMEOUT;
        }
    }
    PCD_WriteRegister(CommandReg, PcdIdle); // Stop calculating CRC for new content in the FIFO.

    // Transfer the result from the registers to the result buffer
    result[0] = PCD_ReadRegister(CrcResultRegL);
    result[1] = PCD_ReadRegister(CrcResultRegH);
    return STATUS_OK;
}

//-----------------------------------------------------------------------------------
// Functions for manipulating the MFRC522
//-----------------------------------------------------------------------------------

void Device::PCD_Init() const
{
    if (Gpio::readPin(RSTPIN) == LOW) {
        // The MFRC522 chip is in power down mode.

        // Exit power down mode. This triggers a hard reset.
        Gpio::writePin(RSTPIN, HIGH);

        // Section 8.8.2 in the datasheet says the oscillator start-up time is the start
        // up time of the crystal + 37,74�s. Let us be generous: 50ms.
        delayMS(50);
    } else {
        // Perform a soft reset
        PCD_Reset();
    }

    // When communicating with a PICC we need a timeout if something goes wrong.
    // f_timer = 13.56 MHz / (2*TPreScaler+1) where TPreScaler = [TPrescaler_Hi:TPrescaler_Lo].
    // TPrescaler_Hi are the four low bits in TModeReg. TPrescaler_Lo is TPrescalerReg.
    PCD_WriteRegister(
        TModeReg,
        0x80); // TAuto=1; timer starts automatically at the end of the transmission in all communication modes at all speeds
    PCD_WriteRegister(
        TPrescalerReg,
        0xA9); // TPreScaler = TModeReg[3..0]:TPrescalerReg, ie 0x0A9 = 169 => f_timer=40kHz, ie a timer period of 25�s.
    PCD_WriteRegister(TReloadRegH, 0x03); // Reload timer with 0x3E8 = 1000, ie 25ms before timeout.
    PCD_WriteRegister(TReloadRegL, 0xE8);

    PCD_WriteRegister(
        TxASKReg,
        0x40); // Default 0x00. Force a 100 % ASK modulation independent of the ModGsPReg register setting
    PCD_WriteRegister(
        ModeReg,
        0x3D); // Default 0x3F. Set the preset value for the CRC coprocessor for the CalcCRC command to 0x6363 (ISO 14443-3 part 6.2.4)
    PCD_AntennaOn(); // Enable the antenna driver pins TX1 and TX2 (they were disabled by the reset)
}

void Device::PCD_Reset() const
{
    PCD_WriteRegister(CommandReg, PcdSoftReset); // Issue the SoftReset command.
    // The datasheet does not mention how long the SoftRest command takes to complete.
    // But the MFRC522 might have been in soft power-down mode (triggered by bit 4 of CommandReg)
    // Section 8.8.2 in the datasheet says the oscillator start-up time is the start up time of the crystal + 37,74�s. Let us be generous: 50ms.
    delayMS(50);
    // Wait for the PowerDown bit in CommandReg to be cleared
    while (PCD_ReadRegister(CommandReg) & (1u << 4u)) {
        // PCD still restarting - unlikely after waiting 50ms, but better safe than sorry.
    }
}

void Device::PCD_AntennaOn() const
{
    uint8_t value = PCD_ReadRegister(TxControlReg);
    if ((value & 0x03u) != 0x03) {
        PCD_WriteRegister(TxControlReg, value | 0x03u);
    }
}

[[maybe_unused]] void Device::PCD_AntennaOff() const
{
    PCD_ClearRegisterBitMask(TxControlReg, 0x03);
}

uint8_t Device::PCD_GetAntennaGain() const
{
    return PCD_ReadRegister(RFCfgReg) & (0x07u << 4u);
}

[[maybe_unused]] void Device::PCD_SetAntennaGain(uint8_t mask) const
{
    if (PCD_GetAntennaGain() != mask) {                         // only bother if there is a change
        PCD_ClearRegisterBitMask(RFCfgReg, (0x07u << 4u));      // clear needed to allow 000 pattern
        PCD_SetRegisterBitMask(RFCfgReg, mask & (0x07u << 4u)); // only set RxGain[2:0] bits
    }
}

[[maybe_unused]] bool Device::PCD_PerformSelfTest() const
{
    // This follows directly the steps outlined in 16.1.1
    // 1. Perform a soft reset.
    PCD_Reset();

    // 2. Clear the internal buffer by writing 25 bytes of 00h
    uint8_t ZEROES[25] = {0x00};
    PCD_SetRegisterBitMask(FIFOLevelReg, 0x80); // flush the FIFO buffer
    PCD_WriteRegister(FIFODataReg, 25, ZEROES); // write 25 bytes of 00h to FIFO
    PCD_WriteRegister(CommandReg, PcdMem);     // transfer to internal buffer

    // 3. Enable self-test
    PCD_WriteRegister(AutoTestReg, 0x09);

    // 4. Write 00h to FIFO buffer
    PCD_WriteRegister(FIFODataReg, 0x00);

    // 5. Start self-test by issuing the CalcCRC command
    PCD_WriteRegister(CommandReg, PcdCalcCrc);

    // 6. Wait for self-test to complete
    uint16_t i;
    uint8_t n;
    for (i = 0; i < 0xFF; i++) {
        n = PCD_ReadRegister(
            DivIrqReg); // DivIrqReg[7..0] bits are: Set2 reserved reserved MfinActIRq reserved CRCIRq reserved reserved
        if (n & 0x04u) { // CRCIRq bit set - calculation done
            break;
        }
    }
    PCD_WriteRegister(CommandReg, PcdIdle); // Stop calculating CRC for new content in the FIFO.

    // 7. Read out resulting 64 bytes from the FIFO buffer.
    uint8_t result[64];
    PCD_ReadRegister(FIFODataReg, 64, result, 0);

    // Auto self-test done
    // Reset AutoTestReg register to be 0 again. Required for normal operation.
    PCD_WriteRegister(AutoTestReg, 0x00);

    // Determine firmware version (see section 9.3.4.8 in spec)
    uint8_t version = PCD_ReadRegister(VersionReg);

    // Pick the appropriate reference values
    const uint8_t *reference;
    switch (version) {
    case 0x91: // Version 1.0
        reference = firmwareReferenceV1_0;
        break;
    case 0x92: // Version 2.0
        reference = firmwareReverenceV2_0;
        break;
    default: // Unknown version
        return false;
    }

    // Verify that the results match up to our expectations
    for (i = 0; i < 64; i++) {
        if (result[i] != reference[i]) {
            return false;
        }
    }
    // Test passed; all is good.
    return true;
}

//-----------------------------------------------------------------------------------
// Functions for communicating with PICCs
//-----------------------------------------------------------------------------------

uint8_t Device::PCD_TransceiveData(uint8_t *sendData,
                                   uint8_t sendLen,
                                   uint8_t *backData,
                                   uint8_t *backLen,
                                   uint8_t *validBits,
                                   uint8_t rxAlign,
                                   bool checkCrc) const
{
    uint8_t waitIRq = 0x30; // RxIRq and IdleIRq
    return PCD_CommunicateWithPICC(PcdTransceive,
                                   waitIRq,
                                   sendData,
                                   sendLen,
                                   backData,
                                   backLen,
                                   validBits,
                                   rxAlign,
                                   checkCrc);
}

uint8_t Device::PCD_CommunicateWithPICC(uint8_t command,
                                        uint8_t waitIRq,
                                        uint8_t *sendData,
                                        uint8_t sendLen,
                                        uint8_t *backData,
                                        uint8_t *backLen,
                                        uint8_t *validBits,
                                        uint8_t rxAlign,
                                        bool checkCrc) const
{
    uint8_t n, _validBits;
    unsigned int i;

    // Prepare values for BitFramingReg
    uint8_t txLastBits = validBits ? *validBits : 0;
    uint8_t bitFraming
        = (rxAlign << 4u)
          + txLastBits; // RxAlign = BitFramingReg[6..4]. TxLastBits = BitFramingReg[2..0]

    PCD_WriteRegister(CommandReg, PcdIdle);           // Stop any active command.
    PCD_WriteRegister(ComIrqReg, 0x7F);                // Clear all seven interrupt request bits
    PCD_SetRegisterBitMask(FIFOLevelReg, 0x80);        // FlushBuffer = 1, FIFO initialization
    PCD_WriteRegister(FIFODataReg, sendLen, sendData); // Write sendData to the FIFO
    PCD_WriteRegister(BitFramingReg, bitFraming);      // Bit adjustments
    PCD_WriteRegister(CommandReg, command);            // Execute the command
    if (command == PcdTransceive) {
        PCD_SetRegisterBitMask(BitFramingReg, 0x80); // StartSend=1, transmission of data starts
    }

    // Wait for the command to complete.
    // In PCD_Init() we set the TAuto flag in TModeReg. This means the timer automatically starts when the PCD stops transmitting.
    // Each iteration of the do-while-loop takes 17.86�s.
    i = 2000;
    while (true) {
        n = PCD_ReadRegister(
            ComIrqReg); // ComIrqReg[7..0] bits are: Set1 TxIRq RxIRq IdleIRq HiAlertIRq LoAlertIRq ErrIRq TimerIRq
        if (n & waitIRq) { // One of the interrupts that signal success has been set.
            break;
        }
        if (n & 0x01u) { // Timer interrupt - nothing received in 25ms
            return STATUS_TIMEOUT;
        }
        if (--i
            == 0) { // The emergency break. If all other condions fail we will eventually terminate on this one after 35.7ms. Communication with the MFRC522 might be down.
            return STATUS_TIMEOUT;
        }
    }

    // Stop now if any errors except collisions were detected.
    uint8_t errorRegValue = PCD_ReadRegister(
        ErrorReg); // ErrorReg[7..0] bits are: WrErr TempErr reserved BufferOvfl CollErr CRCErr ParityErr ProtocolErr
    if (errorRegValue & 0x13u) { // BufferOvfl ParityErr ProtocolErr
        return STATUS_ERROR;
    }

    // If the caller wants data back, get it from the MFRC522.
    if (backData && backLen) {
        n = PCD_ReadRegister(FIFOLevelReg); // Number of bytes in the FIFO
        if (n > *backLen) {
            return STATUS_NO_ROOM;
        }
        *backLen = n;                                        // Number of bytes returned
        PCD_ReadRegister(FIFODataReg, n, backData, rxAlign); // Get received data from FIFO
        _validBits
            = PCD_ReadRegister(ControlReg)
              & 0x07u; // RxLastBits[2:0] indicates the number of valid bits in the last received uint8_t. If this value is 000b, the whole uint8_t is valid.
        if (validBits) {
            *validBits = _validBits;
        }
    }

    // Tell about collisions
    if (errorRegValue & 0x08u) { // CollErr
        return STATUS_COLLISION;
    }

    // Perform CRC_A validation if requested.
    if (backData && backLen && checkCrc) {
        // In this case a MIFARE Classic NAK is not OK.
        if (*backLen == 1 && _validBits == 4) {
            return STATUS_MIFARE_NACK;
        }
        // We need at least the CRC_A value and all 8 bits of the last uint8_t must be received.
        if (*backLen < 2 || _validBits != 0) {
            return STATUS_CRC_WRONG;
        }
        // Verify CRC_A - do our own calculation and store the control in controlBuffer.
        uint8_t controlBuffer[2];
        n = PCD_CalculateCRC(&backData[0], *backLen - 2, &controlBuffer[0]);
        if (n != STATUS_OK) {
            return n;
        }
        if ((backData[*backLen - 2] != controlBuffer[0])
            || (backData[*backLen - 1] != controlBuffer[1])) {
            return STATUS_CRC_WRONG;
        }
    }

    return STATUS_OK;
}

uint8_t Device::PICC_RequestA(uint8_t *bufferAtqa, uint8_t *bufferSize) const
{
    return PICC_ReqaOrWupa(PICC_CMD_REQA, bufferAtqa, bufferSize);
}

uint8_t Device::PICC_WakeupA(uint8_t *bufferAtqa, uint8_t *bufferSize) const
{
    return PICC_ReqaOrWupa(PICC_CMD_WUPA, bufferAtqa, bufferSize);
}

uint8_t Device::PICC_ReqaOrWupa(uint8_t command, uint8_t *bufferAtqa, uint8_t *bufferSize) const
{
    uint8_t validBits;
    uint8_t status;

    if (bufferAtqa == nullptr || *bufferSize < 2) { // The ATQA response is 2 bytes long.
        return STATUS_NO_ROOM;
    }
    PCD_ClearRegisterBitMask(CollReg,
                             0x80); // ValuesAfterColl=1 => Bits received after collision are cleared.
    validBits
        = 7; // For REQA and WUPA we need the short frame format - transmit only 7 bits of the last (and only) uint8_t. TxLastBits = BitFramingReg[2..0]
    status = PCD_TransceiveData(&command, 1, bufferAtqa, bufferSize, &validBits);
    if (status != STATUS_OK) {
        return status;
    }
    if (*bufferSize != 2 || validBits != 0) { // ATQA must be exactly 16 bits.
        return STATUS_ERROR;
    }
    return STATUS_OK;
}

uint8_t Device::PICC_Select(Uid *uid, uint8_t validBits) const
{
    bool uidComplete;
    bool selectDone;
    bool useCascadeTag;
    uint8_t cascadeLevel = 1;
    uint8_t result;
    uint8_t count;
    uint8_t index;
    uint8_t uidIndex; // The first index in uid->uidByte[] that is used in the current Cascade Level.
    signed char currentLevelKnownBits; // The number of known UID bits in the current Cascade Level.
    uint8_t buffer[9]; // The SELECT/ANTICOLLISION commands uses a 7 uint8_t standard frame + 2 bytes CRC_A
    uint8_t bufferUsed; // The number of bytes used in the buffer, ie the number of bytes to transfer to the FIFO.
    uint8_t rxAlign; // Used in BitFramingReg. Defines the bit position for the first bit received.
    uint8_t txLastBits; // Used in BitFramingReg. The number of valid bits in the last transmitted uint8_t.
    uint8_t *responseBuffer;
    uint8_t responseLength;

    // Description of buffer structure:
    //		uint8_t 0: SEL 				Indicates the Cascade Level: PICC_CMD_SEL_CL1, PICC_CMD_SEL_CL2 or PICC_CMD_SEL_CL3
    //		uint8_t 1: NVB					Number of Valid Bits (in complete command, not just the UID): High nibble: complete bytes, Low nibble: Extra bits.
    //		uint8_t 2: UID-data or CT		See explanation below. CT means Cascade Tag.
    //		uint8_t 3: UID-data
    //		uint8_t 4: UID-data
    //		uint8_t 5: UID-data
    //		uint8_t 6: BCC					Block Check Character - XOR of bytes 2-5
    //		uint8_t 7: CRC_A
    //		uint8_t 8: CRC_A
    // The BCC and CRC_A is only transmitted if we know all the UID bits of the current Cascade Level.
    //
    // Description of bytes 2-5: (Section 6.5.4 of the ISO/IEC 14443-3 draft: UID contents and cascade levels)
    //		UID size	Cascade level	uint8_t2	uint8_t3	uint8_t4	uint8_t5
    //		========	=============	=====	=====	=====	=====
    //		 4 bytes		1			uid0	uid1	uid2	uid3
    //		 7 bytes		1			CT		uid0	uid1	uid2
    //						2			uid3	uid4	uid5	uid6
    //		10 bytes		1			CT		uid0	uid1	uid2
    //						2			CT		uid3	uid4	uid5
    //						3			uid6	uid7	uid8	uid9

    // Sanity checks
    if (validBits > 80) {
        return STATUS_INVALID;
    }

    // Prepare MFRC522
    PCD_ClearRegisterBitMask(CollReg,
                             0x80); // ValuesAfterColl=1 => Bits received after collision are cleared.

    // Repeat Cascade Level loop until we have a complete UID.
    uidComplete = false;
    while (!uidComplete) {
        // Set the Cascade Level in the SEL uint8_t, find out if we need to use the Cascade Tag in uint8_t 2.
        switch (cascadeLevel) {
        case 1:
            buffer[0] = PICC_CMD_SEL_CL1;
            uidIndex = 0;
            useCascadeTag = validBits
                            && uid->size > 4; // When we know that the UID has more than 4 bytes
            break;

        case 2:
            buffer[0] = PICC_CMD_SEL_CL2;
            uidIndex = 3;
            useCascadeTag = validBits
                            && uid->size > 7; // When we know that the UID has more than 7 bytes
            break;

        case 3:
            buffer[0] = PICC_CMD_SEL_CL3;
            uidIndex = 6;
            useCascadeTag = false; // Never used in CL3.
            break;

        default:
            return STATUS_INTERNAL_ERROR;
            break;
        }

        // How many UID bits are known in this Cascade Level?
        currentLevelKnownBits = static_cast<char>(validBits - (8 * uidIndex));
        if (currentLevelKnownBits < 0) {
            currentLevelKnownBits = 0;
        }
        // Copy the known bits from uid->uidByte[] to buffer[]
        index = 2; // destination index in buffer[]
        if (useCascadeTag) {
            buffer[index++] = PICC_CMD_CT;
        }
        uint8_t bytesToCopy
            = currentLevelKnownBits / 8
              + (currentLevelKnownBits % 8
                     ? 1
                     : 0); // The number of bytes needed to represent the known bits for this level.
        if (bytesToCopy) {
            uint8_t maxbytes
                = useCascadeTag
                      ? 3
                      : 4; // Max 4 bytes in each Cascade Level. Only 3 left if we use the Cascade Tag
            if (bytesToCopy > maxbytes) {
                bytesToCopy = maxbytes;
            }
            for (count = 0; count < bytesToCopy; count++) {
                buffer[index++] = uid->uidByte[uidIndex + count];
            }
        }
        // Now that the data has been copied we need to include the 8 bits in CT in currentLevelKnownBits
        if (useCascadeTag) {
            currentLevelKnownBits += 8;
        }

        // Repeat anti collision loop until we can transmit all UID bits + BCC and receive a SAK - max 32 iterations.
        selectDone = false;
        while (!selectDone) {
            // Find out how many bits and bytes to send and receive.
            if (currentLevelKnownBits
                >= 32) { // All UID bits in this Cascade Level are known. This is a SELECT.
                //Serial.print(F("SELECT: currentLevelKnownBits=")); Serial.println(currentLevelKnownBits, DEC);
                buffer[1] = 0x70; // NVB - Number of Valid Bits: Seven whole bytes
                // Calculate BCC - Block Check Character
                buffer[6] = buffer[2]
                            ^ static_cast<uint8_t>(
                                buffer[3]
                                ^ static_cast<uint8_t>(buffer[4] ^ static_cast<uint8_t>(buffer[5])));
                // Calculate CRC_A
                result = PCD_CalculateCRC(buffer, 7, &buffer[7]);
                if (result != STATUS_OK) {
                    return result;
                }
                txLastBits = 0; // 0 => All 8 bits are valid.
                bufferUsed = 9;
                // Store response in the last 3 bytes of buffer (BCC and CRC_A - not needed after tx)
                responseBuffer = &buffer[6];
                responseLength = 3;
            } else { // This is an ANTICOLLISION.
                //Serial.print(F("ANTICOLLISION: currentLevelKnownBits=")); Serial.println(currentLevelKnownBits, DEC);
                txLastBits = currentLevelKnownBits % 8;
                count = currentLevelKnownBits / 8;      // Number of whole bytes in the UID part.
                index = 2 + count;                      // Number of whole bytes: SEL + NVB + UIDs
                buffer[1] = (index << 4u) + txLastBits; // NVB - Number of Valid Bits
                bufferUsed = index + (txLastBits ? 1 : 0);
                // Store response in the unused part of buffer
                responseBuffer = &buffer[index];
                responseLength = sizeof(buffer) - index;
            }

            // Set bit adjustments
            rxAlign
                = txLastBits; // Having a seperate variable is overkill. But it makes the next line easier to read.
            PCD_WriteRegister(
                BitFramingReg,
                (rxAlign << 4u)
                    + txLastBits); // RxAlign = BitFramingReg[6..4]. TxLastBits = BitFramingReg[2..0]

            // Transmit the buffer and receive the response.
            result = PCD_TransceiveData(buffer,
                                        bufferUsed,
                                        responseBuffer,
                                        &responseLength,
                                        &txLastBits,
                                        rxAlign);
            if (result == STATUS_COLLISION) { // More than one PICC in the field => collision.
                result = PCD_ReadRegister(
                    CollReg); // CollReg[7..0] bits are: ValuesAfterColl reserved CollPosNotValid CollPos[4:0]
                if (result & 0x20u) {        // CollPosNotValid
                    return STATUS_COLLISION; // Without a valid collision position we cannot continue
                }
                uint8_t collisionPos = result & 0x1Fu; // Values 0-31, 0 means bit 32.
                if (collisionPos == 0) {
                    collisionPos = 32;
                }
                if (collisionPos <= currentLevelKnownBits) { // No progress - should not happen
                    return STATUS_INTERNAL_ERROR;
                }
                // Choose the PICC with the bit set.
                currentLevelKnownBits = collisionPos;
                count = (currentLevelKnownBits - 1) % 8; // The bit to modify
                index = 1 + (currentLevelKnownBits / 8)
                        + (count ? 1 : 0); // First uint8_t is index 0.
                buffer[index] |= (1u << count);
            } else if (result != STATUS_OK) {
                return result;
            } else {                               // STATUS_OK
                if (currentLevelKnownBits >= 32) { // This was a SELECT.
                    selectDone = true;             // No more anticollision
                    // We continue below outside the while.
                } else { // This was an ANTICOLLISION.
                    // We now have all 32 bits of the UID in this Cascade Level
                    currentLevelKnownBits = 32;
                    // Run loop again to do the SELECT.
                }
            }
        } // End of while (!selectDone)

        // We do not check the CBB - it was constructed by us above.

        // Copy the found UID bytes from buffer[] to uid->uidByte[]
        index = (buffer[2] == PICC_CMD_CT) ? 3 : 2; // source index in buffer[]
        bytesToCopy = (buffer[2] == PICC_CMD_CT) ? 3 : 4;
        for (count = 0; count < bytesToCopy; count++) {
            uid->uidByte[uidIndex + count] = buffer[index++];
        }

        // Check response SAK (Select Acknowledge)
        if (responseLength != 3
            || txLastBits != 0) { // SAK must be exactly 24 bits (1 uint8_t + CRC_A).
            return STATUS_ERROR;
        }
        // Verify CRC_A - do our own calculation and store the control in buffer[2..3] - those bytes are not needed anymore.
        result = PCD_CalculateCRC(responseBuffer, 1, &buffer[2]);
        if (result != STATUS_OK) {
            return result;
        }
        if ((buffer[2] != responseBuffer[1]) || (buffer[3] != responseBuffer[2])) {
            return STATUS_CRC_WRONG;
        }
        if (responseBuffer[0] & 0x04u) { // Cascade bit set - UID not complete yes
            cascadeLevel++;
        } else {
            uidComplete = true;
            uid->sak = responseBuffer[0];
        }
    } // End of while (!uidComplete)

    // Set correct uid->size
    uid->size = 3 * cascadeLevel + 1;

    return STATUS_OK;
}

uint8_t Device::PICC_HaltA() const
{
    uint8_t result;
    uint8_t buffer[4];

    // Build command buffer
    buffer[0] = PICC_CMD_HLTA;
    buffer[1] = 0;
    // Calculate CRC_A
    result = PCD_CalculateCRC(buffer, 2, &buffer[2]);
    if (result != STATUS_OK) {
        return result;
    }

    // Send the command.
    // The standard says:
    //		If the PICC responds with any modulation during a period of 1 ms after the end of the frame containing the
    //		HLTA command, this response shall be interpreted as 'not acknowledge'.
    // We interpret that this way: Only STATUS_TIMEOUT is an success.
    result = PCD_TransceiveData(buffer, sizeof(buffer), nullptr, nullptr);
    if (result == STATUS_TIMEOUT) {
        return STATUS_OK;
    }
    if (result == STATUS_OK) { // That is ironically NOT ok in this case ;-)
        return STATUS_ERROR;
    }
    return result;
}

//-----------------------------------------------------------------------------------
// Functions for communicating with MIFARE PICCs
//-----------------------------------------------------------------------------------

uint8_t Device::PCD_Authenticate(uint8_t command, uint8_t blockAddr, MifareKey *key, Uid *uid) const
{
    uint8_t waitIRq = 0x10; // IdleIRq

    // Build command buffer
    uint8_t sendData[12];
    sendData[0] = command;
    sendData[1] = blockAddr;
    for (uint8_t i = 0; i < MF_KEY_SIZE; i++) { // 6 key bytes
        sendData[2 + i] = key->keyByte[i];
    }
    for (uint8_t i = 0; i < 4; i++) { // The first 4 bytes of the UID
        sendData[8 + i] = uid->uidByte[i];
    }

    // Start the authentication.
    return PCD_CommunicateWithPICC(PcdMfAuthent, waitIRq, &sendData[0], sizeof(sendData));
}

void Device::PCD_StopCrypto1() const
{
    // Clear MFCrypto1On bit
    PCD_ClearRegisterBitMask(
        Status2Reg,
        0x08); // Status2Reg[7..0] bits are: TempSensClear I2CForceHS reserved reserved MFCrypto1On ModemState[2:0]
}

uint8_t Device::MIFARE_Read(uint8_t blockAddr, uint8_t *buffer, uint8_t *bufferSize) const
{
    uint8_t result;

    // Sanity check
    if (buffer == nullptr || *bufferSize < 18) {
        return STATUS_NO_ROOM;
    }

    // Build command buffer
    buffer[0] = PICC_CMD_MF_READ;
    buffer[1] = blockAddr;
    // Calculate CRC_A
    result = PCD_CalculateCRC(buffer, 2, &buffer[2]);
    if (result != STATUS_OK) {
        return result;
    }

    // Transmit the buffer and receive the response, validate CRC_A.
    return PCD_TransceiveData(buffer, 4, buffer, bufferSize, nullptr, 0, true);
}

uint8_t Device::MIFARE_Write(uint8_t blockAddr, uint8_t *buffer, uint8_t bufferSize) const
{
    uint8_t result;

    // Sanity check
    if (buffer == nullptr || bufferSize < 16) {
        return STATUS_INVALID;
    }

    // Mifare Classic protocol requires two communications to perform a write.
    // Step 1: Tell the PICC we want to write to block blockAddr.
    uint8_t cmdBuffer[2];
    cmdBuffer[0] = PICC_CMD_MF_WRITE;
    cmdBuffer[1] = blockAddr;
    result = PCD_MIFARE_Transceive(cmdBuffer,
                                   2); // Adds CRC_A and checks that the response is MF_ACK.
    if (result != STATUS_OK) {
        return result;
    }

    // Step 2: Transfer the data
    result = PCD_MIFARE_Transceive(buffer,
                                   bufferSize); // Adds CRC_A and checks that the response is MF_ACK.
    if (result != STATUS_OK) {
        return result;
    }

    return STATUS_OK;
}

[[maybe_unused]] uint8_t Device::MIFARE_Ultralight_Write(uint8_t page,
                                                         uint8_t *buffer,
                                                         uint8_t bufferSize) const
{
    uint8_t result;

    // Sanity check
    if (buffer == nullptr || bufferSize < 4) {
        return STATUS_INVALID;
    }

    // Build commmand buffer
    uint8_t cmdBuffer[6];
    cmdBuffer[0] = PICC_CMD_UL_WRITE;
    cmdBuffer[1] = page;
    memcpy(&cmdBuffer[2], buffer, 4);

    // Perform the write
    result = PCD_MIFARE_Transceive(cmdBuffer,
                                   6); // Adds CRC_A and checks that the response is MF_ACK.
    if (result != STATUS_OK) {
        return result;
    }
    return STATUS_OK;
}

[[maybe_unused]] uint8_t Device::MIFARE_Decrement(uint8_t blockAddr, long delta)
{
    return MIFARE_TwoStepHelper(PICC_CMD_MF_DECREMENT, blockAddr, delta);
}

[[maybe_unused]] uint8_t Device::MIFARE_Increment(uint8_t blockAddr, long delta)
{
    return MIFARE_TwoStepHelper(PICC_CMD_MF_INCREMENT, blockAddr, delta);
}

[[maybe_unused]] uint8_t Device::MIFARE_Restore(uint8_t blockAddr)
{
    // The datasheet describes Restore as a two step operation, but does not explain what data to transfer in step 2.
    // Doing only a single step does not work, so I chose to transfer 0L in step two.
    return MIFARE_TwoStepHelper(PICC_CMD_MF_RESTORE, blockAddr, 0L);
}

uint8_t Device::MIFARE_TwoStepHelper(uint8_t command, uint8_t blockAddr, long data) const
{
    uint8_t result;
    uint8_t cmdBuffer[2]; // We only need room for 2 bytes.

    // Step 1: Tell the PICC the command and block address
    cmdBuffer[0] = command;
    cmdBuffer[1] = blockAddr;
    result = PCD_MIFARE_Transceive(cmdBuffer,
                                   2); // Adds CRC_A and checks that the response is MF_ACK.
    if (result != STATUS_OK) {
        return result;
    }

    // Step 2: Transfer the data
    result = PCD_MIFARE_Transceive((uint8_t *) &data,
                                   4,
                                   true); // Adds CRC_A and accept timeout as success.
    if (result != STATUS_OK) {
        return result;
    }

    return STATUS_OK;
}

[[maybe_unused]] uint8_t Device::MIFARE_Transfer(uint8_t blockAddr) const
{
    uint8_t result;
    uint8_t cmdBuffer[2]; // We only need room for 2 bytes.

    // Tell the PICC we want to transfer the result into block blockAddr.
    cmdBuffer[0] = PICC_CMD_MF_TRANSFER;
    cmdBuffer[1] = blockAddr;
    result = PCD_MIFARE_Transceive(cmdBuffer,
                                   2); // Adds CRC_A and checks that the response is MF_ACK.
    if (result != STATUS_OK) {
        return result;
    }
    return STATUS_OK;
}

uint8_t Device::MIFARE_GetValue(uint8_t blockAddr, long *value) const
{
    uint8_t status;
    uint8_t buffer[18];
    uint8_t size = sizeof(buffer);

    // Read the block
    status = MIFARE_Read(blockAddr, buffer, &size);
    if (status == STATUS_OK) {
        // Extract the value
        *value = (long(buffer[3]) << 24) | (long(buffer[2]) << 16) | (long(buffer[1]) << 8)
                 | long(buffer[0]);
    }
    return status;
}

uint8_t Device::MIFARE_SetValue(uint8_t blockAddr, long value) const
{
    uint8_t buffer[18];

    // Translate the long into 4 bytes; repeated 2x in value block
    buffer[0] = buffer[8] = (static_cast<unsigned long>(value) & 0xFFu);
    buffer[1] = buffer[9] = (static_cast<unsigned long>(value) & 0xFF00u) >> 8u;
    buffer[2] = buffer[10] = (static_cast<unsigned long>(value) & 0xFF0000u) >> 16u;
    buffer[3] = buffer[11] = (static_cast<unsigned long>(value) & 0xFF000000u) >> 24u;
    // Inverse 4 bytes also found in value block
    buffer[4] = ~buffer[0];
    buffer[5] = ~buffer[1];
    buffer[6] = ~buffer[2];
    buffer[7] = ~buffer[3];
    // Address 2x with inverse address 2x
    buffer[12] = buffer[14] = blockAddr;
    buffer[13] = buffer[15] = ~blockAddr;

    // Write the whole data block
    return MIFARE_Write(blockAddr, buffer, 16);
}

//-----------------------------------------------------------------------------------
// Support functions
//-----------------------------------------------------------------------------------

uint8_t Device::PCD_MIFARE_Transceive(uint8_t *sendData, uint8_t sendLen, bool acceptTimeout) const
{
    uint8_t result;
    uint8_t cmdBuffer[18]; // We need room for 16 bytes data and 2 bytes CRC_A.

    // Sanity check
    if (sendData == nullptr || sendLen > 16) {
        return STATUS_INVALID;
    }

    // Copy sendData[] to cmdBuffer[] and add CRC_A
    memcpy(cmdBuffer, sendData, sendLen);
    result = PCD_CalculateCRC(cmdBuffer, sendLen, &cmdBuffer[sendLen]);
    if (result != STATUS_OK) {
        return result;
    }
    sendLen += 2;

    // Transceive the data, store the reply in cmdBuffer[]
    uint8_t waitIRq = 0x30; // RxIRq and IdleIRq
    uint8_t cmdBufferSize = sizeof(cmdBuffer);
    uint8_t validBits = 0;
    result = PCD_CommunicateWithPICC(PcdTransceive,
                                     waitIRq,
                                     cmdBuffer,
                                     sendLen,
                                     cmdBuffer,
                                     &cmdBufferSize,
                                     &validBits);
    if (acceptTimeout && result == STATUS_TIMEOUT) {
        return STATUS_OK;
    }
    if (result != STATUS_OK) {
        return result;
    }
    // The PICC must reply with a 4 bit ACK
    if (cmdBufferSize != 1 || validBits != 4) {
        return STATUS_ERROR;
    }
    if (cmdBuffer[0] != MF_ACK) {
        return STATUS_MIFARE_NACK;
    }
    return STATUS_OK;
}

std::string Device::GetStatusCodeName(uint8_t code)
{
    switch (code) {
    case STATUS_OK:
        return "Success.";
    case STATUS_ERROR:
        return "Error in communication.";
    case STATUS_COLLISION:
        return "Collission detected.";
    case STATUS_TIMEOUT:
        return "Timeout in communication.";
    case STATUS_NO_ROOM:
        return "A buffer is not big enough.";
    case STATUS_INTERNAL_ERROR:
        return "Internal error in the code. Should not happen.";
    case STATUS_INVALID:
        return "Invalid argument.";
    case STATUS_CRC_WRONG:
        return "The CRC_A does not match.";
    case STATUS_MIFARE_NACK:
        return "A MIFARE PICC responded with NAK.";
    default:
        return "Unknown error";
    }
}

uint8_t Device::PICC_GetType(uint8_t sak)
{
    if (sak & 0x04u) { // UID not complete
        return PICC_TYPE_NOT_COMPLETE;
    }

    switch (sak) {
    case 0x09:
        return PICC_TYPE_MIFARE_MINI;
        break;
    case 0x08:
        return PICC_TYPE_MIFARE_1K;
        break;
    case 0x18:
        return PICC_TYPE_MIFARE_4K;
        break;
    case 0x00:
        return PICC_TYPE_MIFARE_UL;
        break;
    case 0x10:
    case 0x11:
        return PICC_TYPE_MIFARE_PLUS;
        break;
    case 0x01:
        return PICC_TYPE_TNP3XXX;
        break;
    default:
        break;
    }

    if (sak & 0x20u) {
        return PICC_TYPE_ISO_14443_4;
    }

    if (sak & 0x40u) {
        return PICC_TYPE_ISO_18092;
    }

    return PICC_TYPE_UNKNOWN;
}

std::string Device::PICC_GetTypeName(uint8_t piccType)
{
    switch (piccType) {
    case PICC_TYPE_ISO_14443_4:
        return "PICC compliant with ISO/IEC 14443-4";
    case PICC_TYPE_ISO_18092:
        return "PICC compliant with ISO/IEC 18092 (NFC)";
    case PICC_TYPE_MIFARE_MINI:
        return "MIFARE Mini, 320 bytes";
    case PICC_TYPE_MIFARE_1K:
        return "MIFARE 1KB";
    case PICC_TYPE_MIFARE_4K:
        return "MIFARE 4KB";
    case PICC_TYPE_MIFARE_UL:
        return "MIFARE Ultralight or Ultralight C";
    case PICC_TYPE_MIFARE_PLUS:
        return "MIFARE Plus";
    case PICC_TYPE_TNP3XXX:
        return "MIFARE TNP3XXX";
    case PICC_TYPE_NOT_COMPLETE:
        return "SAK indicates UID is not complete.";
    case PICC_TYPE_UNKNOWN:
    default:
        return "Unknown type";
    }
}

[[maybe_unused]] void Device::PICC_DumpToSerial(Uid *uid) const
{
    MifareKey key;

    // UID
    printf("Card UID:");
    for (uint8_t i = 0; i < uid->size; i++) {
        if (uid->uidByte[i] < 0x10)
            printf(" 0");
        else
            printf(" ");
        printf("%X", uid->uidByte[i]);
    }
    printf("\n");

    // PICC type
    uint8_t piccType = PICC_GetType(uid->sak);
    printf("PICC type: ");
    //Serial.println(PICC_GetTypeName(piccType));
    printf("%s", PICC_GetTypeName(piccType).c_str());

    // Dump contents
    switch (piccType) {
    case PICC_TYPE_MIFARE_MINI:
    case PICC_TYPE_MIFARE_1K:
    case PICC_TYPE_MIFARE_4K:
        // All keys are set to FFFFFFFFFFFFh at chip delivery from the factory.
        for (unsigned char &i : key.keyByte) {
            i = 0xFF;
        }
        PICC_DumpMifareClassicToSerial(uid, piccType, &key);
        break;

    case PICC_TYPE_MIFARE_UL:
        PICC_DumpMifareUltralightToSerial();
        break;

    case PICC_TYPE_ISO_14443_4:
    case PICC_TYPE_ISO_18092:
    case PICC_TYPE_MIFARE_PLUS:
    case PICC_TYPE_TNP3XXX:
        printf("Dumping memory contents not implemented for that PICC type.");
        break;

    case PICC_TYPE_UNKNOWN:
    case PICC_TYPE_NOT_COMPLETE:
    default:
        break; // No memory dump here
    }

    printf("\n");
    PICC_HaltA(); // Already done if it was a MIFARE Classic PICC.
}

void Device::PICC_DumpMifareClassicToSerial(Uid *uid, uint8_t piccType, MifareKey *key) const
{
    uint8_t amountOfSectors = 0;
    switch (piccType) {
    case PICC_TYPE_MIFARE_MINI:
        // Has 5 sectors * 4 blocks/sector * 16 bytes/block = 320 bytes.
        amountOfSectors = 5;
        break;

    case PICC_TYPE_MIFARE_1K:
        // Has 16 sectors * 4 blocks/sector * 16 bytes/block = 1024 bytes.
        amountOfSectors = 16;
        break;

    case PICC_TYPE_MIFARE_4K:
        // Has (32 sectors * 4 blocks/sector + 8 sectors * 16 blocks/sector) * 16 bytes/block = 4096 bytes.
        amountOfSectors = 40;
        break;

    default: // Should not happen. Ignore.
        break;
    }

    // Dump sectors, highest address first.
    if (amountOfSectors) {
        printf("Sector Block   0  1  2  3   4  5  6  7   8  9 10 11  12 13 14 15  AccessBits\n");
        for (char i = static_cast<char>(amountOfSectors - 1); i >= 0; i--) {
            PICC_DumpMifareClassicSectorToSerial(uid, key, i);
        }
    }
    PICC_HaltA(); // Halt the PICC before stopping the encrypted session.
    PCD_StopCrypto1();
}

void Device::PICC_DumpMifareClassicSectorToSerial(Uid *uid, MifareKey *key, uint8_t sector) const
{
    uint8_t status;
    uint8_t firstBlock;   // Address of lowest address to dump actually last block dumped)
    uint8_t no_of_blocks; // Number of blocks in sector
    bool isSectorTrailer; // Set to true while handling the "last" (ie highest address) in the sector.

    // The access bits are stored in a peculiar fashion.
    // There are four groups:
    //		g[3]	Access bits for the sector trailer, block 3 (for sectors 0-31) or block 15 (for sectors 32-39)
    //		g[2]	Access bits for block 2 (for sectors 0-31) or blocks 10-14 (for sectors 32-39)
    //		g[1]	Access bits for block 1 (for sectors 0-31) or blocks 5-9 (for sectors 32-39)
    //		g[0]	Access bits for block 0 (for sectors 0-31) or blocks 0-4 (for sectors 32-39)
    // Each group has access bits [C1 C2 C3]. In this code C1 is MSB and C3 is LSB.
    // The four CX bits are stored together in a nible cx and an inverted nible cx_.
    uint8_t c1, c2, c3;    // Nibbles
    uint8_t c1_, c2_, c3_; // Inverted nibbles
    bool invertedError;    // True if one of the inverted nibbles did not match
    uint8_t g[4];          // Access bits for each of the four groups.
    uint8_t group;         // 0-3 - active group for access bits
    bool firstInGroup;     // True for the first block dumped in the group

    // Determine position and size of sector.
    if (sector < 32) { // Sectors 0..31 has 4 blocks each
        no_of_blocks = 4;
        firstBlock = sector * no_of_blocks;
    } else if (sector < 40) { // Sectors 32-39 has 16 blocks each
        no_of_blocks = 16;
        firstBlock = 128 + (sector - 32) * no_of_blocks;
    } else { // Illegal input, no MIFARE Classic PICC has more than 40 sectors.
        return;
    }

    // Dump blocks, highest address first.
    uint8_t uint8_tCount;
    uint8_t buffer[18];
    uint8_t blockAddr;
    isSectorTrailer = true;
    for (char blockOffset = static_cast<char>(no_of_blocks - 1); blockOffset >= 0; blockOffset--) {
        blockAddr = firstBlock + blockOffset;
        // Sector number - only on first line
        if (isSectorTrailer) {
            if (sector < 10)
                printf("   "); // Pad with spaces
            else
                printf("  "); // Pad with spaces
            printf("%02X", sector);
            printf("   ");
        } else {
            printf("       ");
        }
        // Block number
        if (blockAddr < 10)
            printf("   "); // Pad with spaces
        else {
            if (blockAddr < 100)
                printf("  "); // Pad with spaces
            else
                printf(" "); // Pad with spaces
        }
        printf("%02X", blockAddr);
        printf("  ");
        // Establish encrypted communications before reading the first block
        if (isSectorTrailer) {
            status = PCD_Authenticate(PICC_CMD_MF_AUTH_KEY_A, firstBlock, key, uid);
            if (status != STATUS_OK) {
                printf("PCD_Authenticate() failed: ");
                printf("%s\n", GetStatusCodeName(status).c_str());
                return;
            }
        }
        // Read block
        uint8_tCount = sizeof(buffer);
        status = MIFARE_Read(blockAddr, buffer, &uint8_tCount);
        if (status != STATUS_OK) {
            printf("MIFARE_Read() failed: ");
            printf("%s\n", GetStatusCodeName(status).c_str());
            continue;
        }
        // Dump data
        for (uint8_t index = 0; index < 16; index++) {
            if (buffer[index] < 0x10)
                printf(" 0");
            else
                printf(" ");
            printf("9x%02X", buffer[index]);
            if ((index % 4) == 3) {
                printf(" ");
            }
        }
        // Parse sector trailer data
        if (isSectorTrailer) {
            c1 = buffer[7] >> 4;
            c2 = buffer[8] & 0xF;
            c3 = buffer[8] >> 4;
            c1_ = buffer[6] & 0xF;
            c2_ = buffer[6] >> 4;
            c3_ = buffer[7] & 0xF;
            invertedError = (c1 != (~c1_ & 0xF)) || (c2 != (~c2_ & 0xF)) || (c3 != (~c3_ & 0xF));
            g[0] = ((c1 & 1) << 2) | ((c2 & 1) << 1) | ((c3 & 1) << 0);
            g[1] = ((c1 & 2) << 1) | ((c2 & 2) << 0) | ((c3 & 2) >> 1);
            g[2] = ((c1 & 4) << 0) | ((c2 & 4) >> 1) | ((c3 & 4) >> 2);
            g[3] = ((c1 & 8) >> 1) | ((c2 & 8) >> 2) | ((c3 & 8) >> 3);
            isSectorTrailer = false;
        }

        // Which access group is this block in?
        if (no_of_blocks == 4) {
            group = blockOffset;
            firstInGroup = true;
        } else {
            group = blockOffset / 5;
            firstInGroup = (group == 3) || (group != (blockOffset + 1) / 5);
        }

        if (firstInGroup) {
            // Print access bits
            printf(" [ ");
            printf("%02X", (g[group] >> 2) & 1);
            printf(" ");
            printf("%02X", (g[group] >> 1) & 1);
            printf(" ");
            printf("%02X", (g[group] >> 0) & 1);
            printf(" ] ");

            if (invertedError) {
                printf(" Inverted access bits did not match! ");
            }
        }

        if (group != 3 && (g[group] == 1 || g[group] == 6)) { // Not a sector trailer, a value block
            long value = (long(buffer[3]) << 24) | (long(buffer[2]) << 16) | (long(buffer[1]) << 8)
                         | long(buffer[0]);
            printf(" Value=");
            printf("0x%02lX", value);
            printf(" Adr=");
            printf("0x%02X", buffer[12]);
        }
        printf("\n");
    }
}

void Device::PICC_DumpMifareUltralightToSerial() const
{
    uint8_t status;
    uint8_t uint8_tCount;
    uint8_t buffer[18];
    uint8_t i;

    printf("Page  0  1  2  3");
    // Try the mpages of the original Ultralight. Ultralight C has more pages.
    for (uint8_t page = 0; page < 16; page += 4) { // Read returns data for 4 pages at a time.
        // Read pages
        uint8_tCount = sizeof(buffer);
        status = MIFARE_Read(page, buffer, &uint8_tCount);
        if (status != STATUS_OK) {
            printf("MIFARE_Read() failed: ");
            printf("%s\n", GetStatusCodeName(status).c_str());
            break;
        }
        // Dump data
        for (uint8_t offset = 0; offset < 4; offset++) {
            i = page + offset;
            if (i < 10)
                printf("  "); // Pad with spaces
            else
                printf(" "); // Pad with spaces
            printf("%02X", i);
            printf("  ");
            for (uint8_t index = 0; index < 4; index++) {
                i = 4 * offset + index;
                if (buffer[i] < 0x10)
                    printf(" 0");
                else
                    printf(" ");
                printf("%02X", buffer[i]);
            }
            printf("\n");
        }
    }
}

[[maybe_unused]] void Device::MIFARE_SetAccessBits(
    uint8_t *accessBitBuffer, uint8_t g0, uint8_t g1, uint8_t g2, uint8_t g3)
{
    uint8_t c1 = ((g3 & 4) << 1) | ((g2 & 4) << 0) | ((g1 & 4) >> 1) | ((g0 & 4) >> 2);
    uint8_t c2 = ((g3 & 2) << 2) | ((g2 & 2) << 1) | ((g1 & 2) << 0) | ((g0 & 2) >> 1);
    uint8_t c3 = ((g3 & 1) << 3) | ((g2 & 1) << 2) | ((g1 & 1) << 1) | ((g0 & 1) << 0);

    accessBitBuffer[0] = (static_cast<uint8_t>(~c2) & 0xFu) << 4u
                         | (static_cast<uint8_t>(~c1) & 0xFu);
    accessBitBuffer[1] = c1 << 4u | (static_cast<uint8_t>(~c3) & 0xFu);
    accessBitBuffer[2] = c3 << 4u | c2;
}

bool Device::MIFARE_OpenUidBackdoor(bool logErrors) const
{
    // Magic sequence:
    // > 50 00 57 CD (HALT + CRC)
    // > 40 (7 bits only)
    // < A (4 bits only)
    // > 43
    // < A (4 bits only)
    // Then you can write to sector 0 without authenticating

    PICC_HaltA(); // 50 00 57 CD

    uint8_t cmd = 0x40;
    uint8_t validBits = 7; /* Our command is only 7 bits. After receiving card response,
			 this will contain amount of valid response bits. */
    uint8_t response[32];  // Card's response is written here
    uint8_t received;
    uint8_t status = PCD_TransceiveData(&cmd,
                                        (uint8_t) 1,
                                        response,
                                        &received,
                                        &validBits,
                                        (uint8_t) 0,
                                        false); // 40
    if (status != STATUS_OK) {
        if (logErrors) {
            printf("Card did not respond to 0x40 after HALT command. Are you sure it is a UID "
                   "changeable one?");
            printf("Error name: ");
            printf("%s", GetStatusCodeName(status).c_str());
        }
        return false;
    }
    if (received != 1 || response[0] != 0x0A) {
        if (logErrors) {
            printf("Got bad response on backdoor 0x40 command: ");
            printf("0x%02X", response[0]);
            printf(" (");
            printf("%02X", validBits);
            printf(" valid bits)\r\n");
        }
        return false;
    }

    cmd = 0x43;
    validBits = 8;
    status = PCD_TransceiveData(&cmd,
                                (uint8_t) 1,
                                response,
                                &received,
                                &validBits,
                                (uint8_t) 0,
                                false); // 43
    if (status != STATUS_OK) {
        if (logErrors) {
            printf("Error in communication at command 0x43, after successfully executing 0x40");
            printf("Error name: ");
            printf("%s\n", GetStatusCodeName(status).c_str());
        }
        return false;
    }
    if (received != 1 || response[0] != 0x0A) {
        if (logErrors) {
            printf("Got bad response on backdoor 0x43 command: ");
            printf("%02X", response[0]);
            printf(" (");
            printf("%02X", validBits);
            printf(" valid bits)\r\n");
        }
        return false;
    }

    // You can now write to sector 0 without authenticating!
    return true;
}

[[maybe_unused]] bool Device::MIFARE_SetUid(const uint8_t *newUid, uint8_t uidSize, bool logErrors)
{
    // UID + BCC uint8_t can not be larger than 16 together
    if (!newUid || !uidSize || uidSize > 15) {
        if (logErrors) {
            printf("New UID buffer empty, size 0, or size > 15 given");
        }
        return false;
    }

    // Authenticate for reading
    MifareKey key = {0xFF, 0xFF, 0xFF, 0xFF, 0xFF, 0xFF};
    uint8_t status = PCD_Authenticate(Device::PICC_CMD_MF_AUTH_KEY_A, (uint8_t) 1, &key, &m_uid);
    if (status != STATUS_OK) {
        if (status == STATUS_TIMEOUT) {
            // We get a read timeout if no card is selected yet, so let's select one

            // Wake the card up again if sleeping
            //			  uint8_t atqa_answer[2];
            //			  uint8_t atqa_size = 2;
            //			  PICC_WakeupA(atqa_answer, &atqa_size);

            if (!PICC_IsNewCardPresent() || !PICC_ReadCardSerial()) {
                printf(
                    "No card was previously selected, and none are available. Failed to set UID.");
                return false;
            }

            status = PCD_Authenticate(Device::PICC_CMD_MF_AUTH_KEY_A, (uint8_t) 1, &key, &m_uid);
            if (status != STATUS_OK) {
                // We tried, time to give up
                if (logErrors) {
                    printf("Failed to authenticate to card for reading, could not set UID: ");
                    printf("%s\n", GetStatusCodeName(status).c_str());
                }
                return false;
            }
        } else {
            if (logErrors) {
                printf("PCD_Authenticate() failed: ");
                printf("%s\n", GetStatusCodeName(status).c_str());
            }
            return false;
        }
    }

    // Read block 0
    uint8_t block0_buffer[18];
    uint8_t uint8_tCount = sizeof(block0_buffer);
    status = MIFARE_Read((uint8_t) 0, block0_buffer, &uint8_tCount);
    if (status != STATUS_OK) {
        if (logErrors) {
            printf("MIFARE_Read() failed: ");
            printf("%s\n", GetStatusCodeName(status).c_str());
            printf("Are you sure your KEY A for sector 0 is 0xFFFFFFFFFFFF?");
        }
        return false;
    }

    // Write new UID to the data we just read, and calculate BCC uint8_t
    uint8_t bcc = 0;
    for (int i = 0; i < uidSize; i++) {
        block0_buffer[i] = newUid[i];
        bcc ^= newUid[i];
    }

    // Write BCC uint8_t to buffer
    block0_buffer[uidSize] = bcc;

    // Stop encrypted traffic so we can send raw bytes
    PCD_StopCrypto1();

    // Activate UID backdoor
    if (!MIFARE_OpenUidBackdoor(logErrors)) {
        if (logErrors) {
            printf("Activating the UID backdoor failed.");
        }
        return false;
    }

    // Write modified block 0 back to card
    status = MIFARE_Write((uint8_t) 0, block0_buffer, (uint8_t) 16);
    if (status != STATUS_OK) {
        if (logErrors) {
            printf("MIFARE_Write() failed: ");
            printf("%s\n", GetStatusCodeName(status).c_str());
        }
        return false;
    }

    // Wake the card up again
    uint8_t atqa_answer[2];
    uint8_t atqa_size = 2;
    PICC_WakeupA(atqa_answer, &atqa_size);

    return true;
}

[[maybe_unused]] bool Device::MIFARE_UnbrickUidSector(bool logErrors) const
{
    MIFARE_OpenUidBackdoor(logErrors);

    uint8_t block0_buffer[] = {0x01,
                               0x02,
                               0x03,
                               0x04,
                               0x04,
                               0x00,
                               0x00,
                               0x00,
                               0x00,
                               0x00,
                               0x00,
                               0x00,
                               0x00,
                               0x00,
                               0x00,
                               0x00};

    // Write modified block 0 back to card
    uint8_t status = MIFARE_Write((uint8_t) 0, block0_buffer, (uint8_t) 16);
    if (status != STATUS_OK) {
        if (logErrors) {
            printf("MIFARE_Write() failed: ");
            printf("%s\n", GetStatusCodeName(status).c_str());
        }
        return false;
    }
    return true;
}

//-----------------------------------------------------------------------------------
// Convenience functions - does not add extra functionality
//-----------------------------------------------------------------------------------

bool Device::PICC_IsNewCardPresent() const
{
    uint8_t bufferATQA[2];
    uint8_t bufferSize = sizeof(bufferATQA);
    uint8_t result = PICC_RequestA(bufferATQA, &bufferSize);
    return (result == STATUS_OK || result == STATUS_COLLISION);
}

bool Device::PICC_ReadCardSerial()
{
    uint8_t result = PICC_Select(&m_uid);
    return (result == STATUS_OK);
}

Device::Uid Device::getUid()
{
    return m_uid;
}

void Device::delayMS(int ms)
{
    timespec ts = {
        .tv_sec = 0,
        .tv_nsec = 50 * 1000000,
    };
    nanosleep(&ts, nullptr);
}

} // namespace Mfrc522