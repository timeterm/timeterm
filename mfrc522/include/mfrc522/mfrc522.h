#pragma once

#include "spi.h"
#include <cstdint>
#include <cstdio>
#include <string>

namespace Mfrc522
{

/// Firmware data for self-test
/// Reference values based on firmware version; taken from 16.1.1 in spec.
/// Version 1.0
const std::uint8_t firmwareReferenceV1_0[] = {0x00, 0xC6, 0x37, 0xD5, 0x32, 0xB7, 0x57, 0x5C,
                                              0xC2, 0xD8, 0x7C, 0x4D, 0xD9, 0x70, 0xC7, 0x73,
                                              0x10, 0xE6, 0xD2, 0xAA, 0x5E, 0xA1, 0x3E, 0x5A,
                                              0x14, 0xAF, 0x30, 0x61, 0xC9, 0x70, 0xDB, 0x2E,
                                              0x64, 0x22, 0x72, 0xB5, 0xBD, 0x65, 0xF4, 0xEC,
                                              0x22, 0xBC, 0xD3, 0x72, 0x35, 0xCD, 0xAA, 0x41,
                                              0x1F, 0xA7, 0xF3, 0x53, 0x14, 0xDE, 0x7E, 0x02,
                                              0xD9, 0x0F, 0xB5, 0x5E, 0x25, 0x1D, 0x29, 0x79};
/// Version 2.0
const std::uint8_t firmwareReverenceV2_0[] = {0x00, 0xEB, 0x66, 0xBA, 0x57, 0xBF, 0x23, 0x95,
                                              0xD0, 0xE3, 0x0D, 0x3D, 0x27, 0x89, 0x5C, 0xDE,
                                              0x9D, 0x3B, 0xA7, 0x00, 0x21, 0x5B, 0x89, 0x82,
                                              0x51, 0x3A, 0xEB, 0x02, 0x0C, 0xA5, 0x00, 0x49,
                                              0x7C, 0x84, 0x4D, 0xB3, 0xCC, 0xD2, 0x1B, 0x81,
                                              0x5D, 0x48, 0x76, 0xD5, 0x71, 0x61, 0x21, 0xA9,
                                              0x86, 0x96, 0x83, 0x38, 0xCF, 0x9D, 0x5B, 0x6D,
                                              0xDC, 0x15, 0xBA, 0x3E, 0x7D, 0x95, 0x3B, 0x2F};

class Device
{
public:
    /// MFRC522 registers. Described in chapter 9 of the datasheet.
    /// When using SPI all addresses are shifted one bit left in the "SPI address uint8_t" (section 8.1.2.3)
    enum PcdRegister
    {
        // Page 0: Command and status
        CommandReg = 0x01u << 1u,   ///< starts and stops command execution
        ComIEnReg = 0x02u << 1u,    ///< enable and disable interrupt request control bits
        DivIEnReg = 0x03u << 1u,    ///< enable and disable interrupt request control bits
        ComIrqReg = 0x04u << 1u,    ///< interrupt request bits
        DivIrqReg = 0x05u << 1u,    ///< interrupt request bits
        ErrorReg = 0x06u << 1u,     ///< error bits showing the error status of the last command executed
        Status1Reg = 0x07u << 1u,   ///< communication status bits
        Status2Reg = 0x08u << 1u,   ///< receiver and transmitter status bits
        FIFODataReg = 0x09u << 1u,  ///< input and output of 64 uint8_t FIFO buffer
        FIFOLevelReg = 0x0Au << 1u, ///< number of bytes stored in the FIFO buffer
        WaterLevelReg = 0x0Bu << 1u,///< level for FIFO underflow and overflow warning
        ControlReg = 0x0Cu << 1u,   ///< miscellaneous control registers
        BitFramingReg = 0x0Du << 1u,///< adjustments for bit-oriented frames
        CollReg = 0x0Eu << 1u,      ///< bit position of the first bit-collision detected on the RF interface

        // Page 1: Command
        ModeReg = 0x11u << 1u,       ///< defines general modes for transmitting and receiving
        TxModeReg = 0x12u << 1u,     ///< defines transmission data rate and framing
        RxModeReg = 0x13u << 1u,     ///< defines reception data rate and framing
        TxControlReg = 0x14u << 1u,  ///< controls the logical behavior of the antenna driver pins TX1 and TX2
        TxASKReg = 0x15u << 1u,      ///< controls the setting of the transmission modulation
        TxSelReg = 0x16u << 1u,      ///< selects the internal sources for the antenna driver
        RxSelReg = 0x17u << 1u,      ///< selects internal receiver settings
        RxThresholdReg = 0x18u << 1u,///< selects thresholds for the bit decoder
        DemodReg = 0x19u << 1u,      ///< defines demodulator settings
        MfTxReg = 0x1Cu << 1u,       ///< controls some MIFARE communication transmit parameters
        MfRxReg = 0x1Du << 1u,       ///< controls some MIFARE communication receive parameters
        SerialSpeedReg = 0x1Fu << 1u,///< selects the speed of the serial UART interface

        // Page 2: Configuration
        CrcResultRegH = 0x21u << 1u,///< shows the MSB and LSB values of the CRC calculation
        CrcResultRegL = 0x22u << 1u,
        ModWidthReg = 0x24u << 1u,  ///< controls the ModWidth setting?
        RFCfgReg = 0x26u << 1u,     ///< configures the receiver gain
        GsNReg = 0x27u << 1u,       ///< selects the conductance of the antenna driver pins TX1 and TX2 for modulation
        CWGsPReg = 0x28u << 1u,     ///< defines the conductance of the p-driver output during periods of no modulation
        ModGsPReg = 0x29u << 1u,    ///< defines the conductance of the p-driver output during periods of modulation
        TModeReg = 0x2Au << 1u,     ///< defines settings for the internal timer
        TPrescalerReg = 0x2Bu << 1u,///< the lower 8 bits of the TPrescaler value. The 4 high bits are in TModeReg.
        TReloadRegH = 0x2Cu << 1u,  ///< defines the 16-bit timer reload value
        TReloadRegL = 0x2Du << 1u,
        TCounterValueRegH = 0x2Eu << 1u,///< shows the 16-bit timer value
        TCounterValueRegL = 0x2Fu << 1u,

        // Page 3: Test Registers
        TestSel1Reg = 0x31u << 1u,    ///< general test signal configuration
        TestSel2Reg = 0x32u << 1u,    ///< general test signal configuration
        TestPinEnReg = 0x33u << 1u,   ///< enables pin output driver on pins D1 to D7
        TestPinValueReg = 0x34u << 1u,///< defines the values for D1 to D7 when it is used as an I/O bus
        TestBusReg = 0x35u << 1u,     ///< shows the status of the internal test bus
        AutoTestReg = 0x36u << 1u,    ///< controls the digital self test
        VersionReg = 0x37u << 1u,     ///< shows the software version
        AnalogTestReg = 0x38u << 1u,  ///< controls the pins AUX1 and AUX2
        TestDAC1Reg = 0x39u << 1u,    ///< defines the test value for TestDAC1
        TestDAC2Reg = 0x3Au << 1u,    ///< defines the test value for TestDAC2
        TestADCReg = 0x3Bu << 1u      ///< shows the value of ADC I and Q channels
    };

    /// MFRC522 commands. Described in chapter 10 of the datasheet.
    enum PcdCommand
    {
        PcdIdle = 0x00,            ///< no action, cancels current command execution
        PcdMem = 0x01,             ///< stores 25 bytes into the internal buffer
        PcdGenerateRandomId = 0x02,///< generates a 10-uint8_t random ID number
        PcdCalcCrc = 0x03,         ///< activates the CRC coprocessor or performs a self test
        PcdTransmit = 0x04,        ///< transmits data from the FIFO buffer
        PcdNoCmdChange = 0x07,     ///< no command change, can be used to modify the CommandReg register bits without affecting the command, for example, the PowerDown bit
        PcdReceive = 0x08,         ///< activates the receiver circuits
        PcdTransceive = 0x0C,      ///< transmits data from FIFO buffer to antenna and automatically activates the receiver after transmission
        PcdMfAuthent = 0x0E,       ///< performs the MIFARE standard authentication as a reader
        PcdSoftReset = 0x0F        ///< resets the MFRC522
    };

    /// MFRC522 RxGain[2:0] masks, defines the receiver's signal voltage gain factor (on the PCD).
    /// Described in 9.3.3.6 / table 98 of the datasheet at http://www.nxp.com/documents/data_sheet/MFRC522.pdf
    enum PcdRxGain
    {
        RxGain18dB = 0x00u << 4u,  ///< 000b - 18 dB, minimum
        RxGain23dB = 0x01u << 4u,  ///< 001b - 23 dB
        RxGain18dB_2 = 0x02u << 4u,///< 010b - 18 dB, it seems 010b is a duplicate for 000b
        RxGain23dB_2 = 0x03u << 4u,///< 011b - 23 dB, it seems 011b is a duplicate for 001b
        RxGain33dB = 0x04u << 4u,  ///< 100b - 33 dB, average, and typical default
        RxGain38dB = 0x05u << 4u,  ///< 101b - 38 dB
        RxGain43dB = 0x06u << 4u,  ///< 110b - 43 dB
        RxGain48dB = 0x07u << 4u,  ///< 111b - 48 dB, maximum
        RxGainMin = 0x00u << 4u,   ///< 000b - 18 dB, minimum, convenience for RxGain18dB
        RxGainAvg = 0x04u << 4u,   ///< 100b - 33 dB, average, convenience for RxGain_33dB
        RxGainMax = 0x07u << 4u    ///< 111b - 48 dB, maximum, convenience for RxGain_48dB
    };

    /// Commands sent to the PICC.
    enum PiccCommand
    {
        // The commands used by the PCD to manage communication with several PICCs (ISO 14443-3, Type A, section 6.4)
        PiccCmdReqa = 0x26,  ///< REQuest command, Type A. Invites PICCs in state IDLE to go to READY and prepare for anticollision or selection. 7 bit frame.
        PiccCmdWupa = 0x52,  ///< Wake-UP command, Type A. Invites PICCs in state IDLE and HALT to go to READY(*) and prepare for anticollision or selection. 7 bit frame.
        PiccCmdCt = 0x88,    ///< Cascade Tag. Not really a command, but used during anti collision.
        PiccCmdSelCl1 = 0x93,///< Anti collision/Select, Cascade Level 1
        PiccCmdSelCl2 = 0x95,///< Anti collision/Select, Cascade Level 2
        PiccCmdSelCl3 = 0x97,///< Anti collision/Select, Cascade Level 3
        PiccCmdHlta = 0x50,  ///< HaLT command, Type A. Instructs an ACTIVE PICC to go to state HALT.

        // The commands used for MIFARE Classic (from http://www.nxp.com/documents/data_sheet/MF1S503x.pdf, Section 9)
        // Use PcdMfAuthent to authenticate access to a sector, then use these commands to read/write/modify the blocks on the sector.
        // The read/write commands can also be used for MIFARE Ultralight.
        PiccCmdMfAuthKeyA = 0x60, ///< Perform authentication with Key A
        PiccCmdMfAuthKeyB = 0x61, ///< Perform authentication with Key B
        PiccCmdMfRead = 0x30,     ///< Reads one 16 uint8_t block from the authenticated sector of the PICC. Also used for MIFARE Ultralight.
        PiccCmdMfWrite = 0xA0,    ///< Writes one 16 uint8_t block to the authenticated sector of the PICC. Called "COMPATIBILITY WRITE" for MIFARE Ultralight.
        PiccCmdMfDecrement = 0xC0,///< Decrements the contents of a block and stores the result in the internal data register.
        PiccCmdMfIncrement = 0xC1,///< Increments the contents of a block and stores the result in the internal data register.
        PiccCmdMfRestore = 0xC2,  ///< Reads the contents of a block into the internal data register.
        PiccCmdMfTransfer = 0xB0, ///< Writes the contents of the internal data register to a block.

        // The commands used for MIFARE Ultralight (from http://www.nxp.com/documents/data_sheet/MF0ICU1.pdf, Section 8.6)
        // The PiccCmdMfRead and PiccCmdMfWrite can also be used for MIFARE Ultralight.
        PiccCmdUlWrite = 0xA2///< Writes one 4 uint8_t page to the PICC.
    };

    /// MIFARE constants that does not fit anywhere else
    enum MifareMisc
    {
        MfAck = 0xA, ///< The MIFARE Classic uses a 4 bit ACK/NAK. Any other value than 0xA is NAK.
        MfKeySize = 6///< A Mifare Crypto1 key is 6 bytes.
    };

    /// PICC types we can detect. Remember to update piccGetTypeName() if you add more.
    enum PiccType
    {
        PiccTypeUnknown = 0,
        PiccTypeIso14443_4 = 1,  ///< PICC compliant with ISO/IEC 14443-4
        PiccTypeIso18092 = 2,    ///< PICC compliant with ISO/IEC 18092 (NFC)
        PiccTypeMifareMini = 3,  ///< MIFARE Classic protocol, 320 bytes
        PiccTypeMifare1K = 4,    ///< MIFARE Classic protocol, 1KB
        PiccTypeMifare4K = 5,    ///< MIFARE Classic protocol, 4KB
        PiccTypeMifareUl = 6,    ///< MIFARE Ultralight or Ultralight C
        PiccTypeMifarePlus = 7,  ///< MIFARE Plus
        PiccTypeTnp3xxx = 8,     ///< Only mentioned in NXP AN 10833 MIFARE Type Identification Procedure
        PiccTypeNotComplete = 255///< SAK indicates UID is not complete.
    };

    /// Return codes from the functions in this class. Remember to update getStatusCodeName() if you add more.
    enum StatusCode
    {
        StatusOk = 1,           ///< Success
        StatusError = 2,        ///< Error in communication
        StatusCollision = 3,    ///< Collission detected
        StatusTimeout = 4,      ///< Timeout in communication.
        StatusNoRoom = 5,       ///< A buffer is not big enough.
        StatusInternalError = 6,///< Internal error in the code. Should not happen ;-)
        StatusInvalid = 7,      ///< Invalid argument.
        StatusCrcWrong = 8,     ///< The CRC_A does not match
        StatusMifareNack = 9    ///< A MIFARE PICC responded with NAK.
    };

    /// A struct used for passing the UID of a PICC.
    typedef struct
    {
        /// Number of bytes in the UID. 4, 7 or 10.
        uint8_t size;

        /// UID bytes
        uint8_t uidByte[10];

        /// The SAK (Select acknowledge) uint8_t returned from the PICC after successful selection.
        uint8_t sak;
    } Uid;

    /// A struct used for passing a MIFARE Crypto1 key
    typedef struct
    {
        uint8_t keyByte[MfKeySize];
    } MifareKey;

    /// Size of the MFRC522 FIFO
    [[maybe_unused]] static const uint8_t kFifoSize = 64;// The FIFO is 64 bytes.

    //-----------------------------------------------------------------------------------
    // Functions for setting up the Raspberry Pi
    //-----------------------------------------------------------------------------------

    /**
     * Constructor.
     * Prepares the output pins.
     */
    Device(std::initializer_list<Spi::DeviceOpenOption> spiOptions = {
               Spi::withSpeed(4000000),
           });

    //-----------------------------------------------------------------------------------
    // Basic interface functions for communicating with the MFRC522
    //-----------------------------------------------------------------------------------

    /**
     * Writes a uint8_t to the specified register in the MFRC522 chip.
     * The interface is described in the datasheet section 8.1.2.
     * \param reg The register to write to. One of the PcdRegister enums.
     * \param value The value to write.
     */
    void pcdWriteRegister(uint8_t reg, uint8_t value) const;

    /**
     * Writes a number of bytes to the specified register in the MFRC522 chip.
     * The interface is described in the datasheet section 8.1.2.
     * \param reg The register to write to. One of the PcdRegister enums.
     * \param count The number of bytes to write to the register
     * \param values The values to write. uint8_t array.
     */
    void pcdWriteRegister(uint8_t reg, uint8_t count, uint8_t *values) const;

    /**
     * Reads a uint8_t from the specified register in the MFRC522 chip.
     * The interface is described in the datasheet section 8.1.2.
     * \param reg The register to read from. One of the PcdRegister enums.
     */
    [[nodiscard]] uint8_t pcdReadRegister(uint8_t reg) const;

    /**
     * Reads a number of bytes from the specified register in the MFRC522 chip.
     * The interface is described in the datasheet section 8.1.2.
     * \param reg The register to read from. One of the PcdRegister enums.
     * \param count The number of bytes to read
     * \param values uint8_t array to store the values in.
     * \param rxAlign Only bit positions rxAlign..7 in values[0] are updated.
     */
    void pcdReadRegister(uint8_t reg, uint8_t count, uint8_t *values, uint8_t rxAlign = 0) const;

    /**
     * Sets the bits given in mask in register reg.
     * \param reg The register to update. One of the PcdRegister enums.
     * \param mask The bits to set.
     */
    void pcdSetRegisterBitMask(uint8_t reg, uint8_t mask) const;

    /**
     * Clears the bits given in mask from register reg.
     * \param reg The register to update. One of the PcdRegister enums.
     * \param mask The bits to clear.
     */
    void pcdClearRegisterBitMask(uint8_t reg, uint8_t mask) const;

    /**
     * Use the CRC coprocessor in the MFRC522 to calculate a CRC_A.
     *
     * \param data In: Pointer to the data to transfer to the FIFO for CRC calculation.
     * \param length In: The number of bytes to transfer.
     * \param result Out: Pointer to result buffer. Result is written to result[0..1], low uint8_t first.
     *
     * \return StatusOk on success, STATUS_??? otherwise.
     */
    uint8_t pcdCalculateCrc(uint8_t *data, uint8_t length, uint8_t *result) const;

    //-----------------------------------------------------------------------------------
    // Functions for manipulating the MFRC522
    //-----------------------------------------------------------------------------------

    /**
     * Initializes the MFRC522 chip.
     */
    void pcdInit() const;

    /**
     * Performs a soft reset on the MFRC522 chip and waits for it to be ready again.
     */
    void pcdReset() const;

    /**
     * Turns the antenna on by enabling pins TX1 and TX2.
     * After a reset these pins are disabled.
     */
    void pcdAntennaOn() const;

    /**
     * Turns the antenna off by disabling pins TX1 and TX2.
     */
    [[maybe_unused]] void pcdAntennaOff() const;

    /**
     * Get the current MFRC522 Receiver Gain (RxGain[2:0]) value.
     * See 9.3.3.6 / table 98 in http://www.nxp.com/documents/data_sheet/MFRC522.pdf
     * NOTE: Return value scrubbed with (0x07<<4)=01110000b as RCFfgReg may use reserved bits.
     *
     * \return Value of the RxGain, scrubbed to the 3 bits used.
     */
    [[nodiscard]] uint8_t pcdGetAntennaGain() const;

    /**
     * Set the MFRC522 Receiver Gain (RxGain) to value specified by given mask.
     * See 9.3.3.6 / table 98 in http://www.nxp.com/documents/data_sheet/MFRC522.pdf
     * NOTE: Given mask is scrubbed with (0x07<<4)=01110000b as RCFfgReg may use reserved bits.
     */
    [[maybe_unused]] void pcdSetAntennaGain(uint8_t mask) const;

    /**
     * Performs a self-test of the MFRC522
     * See 16.1.1 in http://www.nxp.com/documents/data_sheet/MFRC522.pdf
     *
     * \return Whether or not the test passed.
     */
    [[maybe_unused]] [[nodiscard]] bool pcdPerformSelfTest() const;

    //-----------------------------------------------------------------------------------
    // Functions for communicating with PICCs
    //-----------------------------------------------------------------------------------

    /**
     * Executes the Transceive command.
     * CRC validation can only be done if backData and backLen are specified.
     *
     * \param sendData Pointer to the data to transfer to the FIFO.
     * \param sendLen Number of bytes to transfer to the FIFO.
     * \param backData NULL or pointer to buffer if data should be read back after executing the command.
     * \param backLen In: Max number of bytes to write to *backData. Out: The number of bytes returned.
     * \param validBits In/Out: The number of valid bits in the last uint8_t. 0 for 8 valid bits. Default NULL.
     * \param rxAlign In: Defines the bit position in backData[0] for the first bit received. Default 0.
     * \param checkCrc In: True => The last two bytes of the response is assumed to be a CRC_A that must be validated.
     *
     * \return StatusOk on success, STATUS_??? otherwise.
     */
    uint8_t pcdTransceiveData(uint8_t *sendData,
                              uint8_t sendLen,
                              uint8_t *backData,
                              uint8_t *backLen,
                              uint8_t *validBits = nullptr,
                              uint8_t rxAlign = 0,
                              bool checkCrc = false) const;

    /**
     * Transfers data to the MFRC522 FIFO, executes a command, waits for completion and transfers data back from the FIFO.
     * CRC validation can only be done if backData and backLen are specified.
     *
     * \param command The command to execute. One of the PcdCommand enums.
     * \param waitIRq The bits in the ComIrqReg register that signals successful completion of the command.
     * \param sendData Pointer to the data to transfer to the FIFO.
     * \param sendLen Number of bytes to transfer to the FIFO.
     * \param backData NULL or pointer to buffer if data should be read back after executing the command.
     * \param backLen In: Max number of bytes to write to *backData. Out: The number of bytes returned.
     * \param validBits In/Out: The number of valid bits in the last uint8_t. 0 for 8 valid bits.
     * \param rxAlign In: Defines the bit position in backData[0] for the first bit received. Default 0.
     * \param checkCrc In: True => The last two bytes of the response is assumed to be a CRC_A that must be validated.
     *
     * \return StatusOk on success, STATUS_??? otherwise.
     */
    uint8_t pcdCommunicateWithPicc(uint8_t command,
                                   uint8_t waitIRq,
                                   uint8_t *sendData,
                                   uint8_t sendLen,
                                   uint8_t *backData = nullptr,
                                   uint8_t *backLen = nullptr,
                                   uint8_t *validBits = nullptr,
                                   uint8_t rxAlign = 0,
                                   bool checkCrc = false) const;

    /**
     * Transmits a REQuest command, Type A. Invites PICCs in state IDLE to go to READY and prepare for anticollision or selection. 7 bit frame.
     * Beware: When two PICCs are in the field at the same time I often get StatusTimeout - probably due do bad antenna design.
     *
     * \param bufferAtqa The buffer to store the ATQA (Answer to request) in
     * \param bufferSize Buffer size, at least two bytes. Also number of bytes returned if StatusOk.
     *
     * \return StatusOk on success, STATUS_??? otherwise.
     */
    uint8_t piccRequestA(uint8_t *bufferAtqa, uint8_t *bufferSize) const;

    /**
     * Transmits a Wake-UP command, Type A. Invites PICCs in state IDLE and HALT to go to READY(*) and prepare for anticollision or selection. 7 bit frame.
     * Beware: When two PICCs are in the field at the same time I often get StatusTimeout - probably due do bad antenna design.
     *
     * \param bufferAtqa The buffer to store the ATQA (Answer to request) in
     * \param bufferSize Buffer size, at least two bytes. Also number of bytes returned if StatusOk.
     *
     * @return StatusOk on success, STATUS_??? otherwise.
     */
    uint8_t piccWakeupA(uint8_t *bufferAtqa, uint8_t *bufferSize) const;

    /**
     * Transmits REQA or WUPA commands.
     * Beware: When two PICCs are in the field at the same time I often get StatusTimeout - probably due do bad antenna design.
     *
     * \param command The command to send - PiccCmdReqa or PiccCmdWupa
     * \param bufferAtqa The buffer to store the ATQA (Answer to request) in
     * \param bufferSize Buffer size, at least two bytes. Also number of bytes returned if StatusOk.
     *
     * @return StatusOk on success, STATUS_??? otherwise.
     */
    uint8_t piccReqaOrWupa(uint8_t command, uint8_t *bufferAtqa, uint8_t *bufferSize) const;

    /**
     * Transmits SELECT/ANTICOLLISION commands to select a single PICC.
     * Before calling this function the PICCs must be placed in the READY(*) state by calling piccRequestA() or piccWakeupA().
     * On success:
     * 		- The chosen PICC is in state ACTIVE(*) and all other PICCs have returned to state IDLE/HALT. (Figure 7 of the ISO/IEC 14443-3 draft.)
     * 		- The UID size and value of the chosen PICC is returned in *uid along with the SAK.
     *
     * A PICC UID consists of 4, 7 or 10 bytes.
     * Only 4 bytes can be specified in a SELECT command, so for the longer UIDs two or three iterations are used:
     * 		UID size	Number of UID bytes		Cascade levels		Example of PICC
     * 		========	===================		==============		===============
     * 		single				 4						1				MIFARE Classic
     * 		double				 7						2				MIFARE Ultralight
     * 		triple				10						3				Not currently in use?
     *
     * \param uid Pointer to Uid struct. Normally output, but can also be used to supply a known UID.
     * \param validBits The number of known UID bits supplied in *uid. Normally 0. If set you must also supply uid->size.
     *
     * \return StatusOk on success, STATUS_??? otherwise.
     */
    uint8_t piccSelect(Uid *uid, uint8_t validBits = 0) const;

#pragma clang diagnostic push
#pragma ide diagnostic ignored "modernize-use-nodiscard"
    /**
     * Instructs a PICC in state ACTIVE(*) to go to state HALT.
     *
     * \return StatusOk on success, STATUS_??? otherwise.
     */
    uint8_t piccHaltA() const;
#pragma clang diagnostic pop

    //-----------------------------------------------------------------------------------
    // Functions for communicating with MIFARE PICCs
    //-----------------------------------------------------------------------------------

    /**
     * Executes the MFRC522 MFAuthent command.
     * This command manages MIFARE authentication to enable a secure communication to any MIFARE Mini, MIFARE 1K and MIFARE 4K card.
     * The authentication is described in the MFRC522 datasheet section 10.3.1.9 and http://www.nxp.com/documents/data_sheet/MF1S503x.pdf section 10.1.
     * For use with MIFARE Classic PICCs.
     * The PICC must be selected - ie in state ACTIVE(*) - before calling this function.
     * Remember to call pcdStopCrypto1() after communicating with the authenticated PICC - otherwise no new communications can start.
     *
     * All keys are set to FFFFFFFFFFFFh at chip delivery.
     *
     * \param command PiccCmdMfAuthKeyA or PiccCmdMfAuthKeyB
     * \param blockAddr The block number. See numbering in the comments in the .h file.
     * \param key Pointer to the Crypteo1 key to use (6 bytes)
     * \param uid Pointer to Uid struct. The first 4 bytes of the UID is used.
     *
     * \return StatusOk on success, STATUS_??? otherwise. Probably StatusTimeout if you supply the wrong key.
     */
    uint8_t pcdAuthenticate(uint8_t command, uint8_t blockAddr, MifareKey *key, Uid *uid) const;

    /**
     * Used to exit the PCD from its authenticated state.
     * Remember to call this function after communicating with an authenticated PICC - otherwise no new communications can start.
     */
    void pcdStopCrypto1() const;

    /**
     * Reads 16 bytes (+ 2 bytes CRC_A) from the active PICC.
     *
     * For MIFARE Classic the sector containing the block must be authenticated before calling this function.
     *
     * For MIFARE Ultralight only addresses 00h to 0Fh are decoded.
     * The MF0ICU1 returns a NAK for higher addresses.
     * The MF0ICU1 responds to the READ command by sending 16 bytes starting from the page address defined by the command argument.
     * For example; if blockAddr is 03h then pages 03h, 04h, 05h, 06h are returned.
     * A roll-back is implemented: If blockAddr is 0Eh, then the contents of pages 0Eh, 0Fh, 00h and 01h are returned.
     *
     * The buffer must be at least 18 bytes because a CRC_A is also returned.
     * Checks the CRC_A before returning StatusOk.
     *
     * \param blockAddr MIFARE Classic: The block (0-0xff) number. MIFARE Ultralight: The first page to return data from.
     * \param buffer The buffer to store the data in
     * \param bufferSize Buffer size, at least 18 bytes. Also number of bytes returned if StatusOk.
     *
     * \return StatusOk on success, STATUS_??? otherwise.
     */
    uint8_t mifareRead(uint8_t blockAddr, uint8_t *buffer, uint8_t *bufferSize) const;

    /**
     * Writes 16 bytes to the active PICC.
     *
     * For MIFARE Classic the sector containing the block must be authenticated before calling this function.
     *
     * For MIFARE Ultralight the operation is called "COMPATIBILITY WRITE".
     * Even though 16 bytes are transferred to the Ultralight PICC, only the least significant 4 bytes (bytes 0 to 3)
     * are written to the specified address. It is recommended to set the remaining bytes 04h to 0Fh to all logic 0.
     *
     * \param blockAddr MIFARE Classic: The block (0-0xff) number. MIFARE Ultralight: The page (2-15) to write to.
     * \param buffer The 16 bytes to write to the PICC
     * \param bufferSize Buffer size, must be at least 16 bytes. Exactly 16 bytes are written.
     *
     * \return StatusOk on success, STATUS_??? otherwise.
     */
    uint8_t mifareWrite(uint8_t blockAddr, uint8_t *buffer, uint8_t bufferSize) const;

    /**
     * MIFARE Decrement subtracts the delta from the value of the addressed block, and stores the result in a volatile memory.
     * For MIFARE Classic only. The sector containing the block must be authenticated before calling this function.
     * Only for blocks in "value block" mode, ie with access bits [C1 C2 C3] = [110] or [001].
     * Use mifareTransfer() to store the result in a block.
     *
     * \param blockAddr The block (0-0xff) number.
     * \param delta This number is subtracted from the value of block blockAddr.
     *
     * \return StatusOk on success, STATUS_??? otherwise.
     */
    [[maybe_unused]] uint8_t mifareDecrement(uint8_t blockAddr, long delta);

    /**
     * MIFARE Increment adds the delta to the value of the addressed block, and stores the result in a volatile memory.
     * For MIFARE Classic only. The sector containing the block must be authenticated before calling this function.
     * Only for blocks in "value block" mode, ie with access bits [C1 C2 C3] = [110] or [001].
     * Use mifareTransfer() to store the result in a block.
     *
     * \param blockAddr The block (0-0xff) number.
     * \param delta This number is added to the value of block blockAddr.
     *
     * @return StatusOk on success, STATUS_??? otherwise.
     */
    [[maybe_unused]] uint8_t mifareIncrement(uint8_t blockAddr, long delta);

    /**
     * MIFARE Restore copies the value of the addressed block into a volatile memory.
     * For MIFARE Classic only. The sector containing the block must be authenticated before calling this function.
     * Only for blocks in "value block" mode, ie with access bits [C1 C2 C3] = [110] or [001].
     * Use mifareTransfer() to store the result in a block.
     *
     * \param blockAddr The block (0-0xff) number.
     *
     * \return StatusOk on success, STATUS_??? otherwise.
     */
    [[maybe_unused]] uint8_t mifareRestore(uint8_t blockAddr);

    /**
     * MIFARE Transfer writes the value stored in the volatile memory into one MIFARE Classic block.
     * For MIFARE Classic only. The sector containing the block must be authenticated before calling this function.
     * Only for blocks in "value block" mode, ie with access bits [C1 C2 C3] = [110] or [001].
     *
     * \param blockAddr The block (0-0xff) number.
     *
     * \return StatusOk on success, STATUS_??? otherwise.
     */
    [[maybe_unused]] [[nodiscard]] uint8_t mifareTransfer(uint8_t blockAddr) const;

    /**
     * Writes a 4 uint8_t page to the active MIFARE Ultralight PICC.
     *
     * \param page The page (2-15) to write to.
     * \param buffer The 4 bytes to write to the PICC
     * \param bufferSize Buffer size, must be at least 4 bytes. Exactly 4 bytes are written.
     *
     * \return StatusOk on success, STATUS_??? otherwise.
     */
    [[maybe_unused]] uint8_t mifareUltralightWrite(uint8_t page,
                                                   uint8_t *buffer,
                                                   uint8_t bufferSize) const;

    /**
     * Helper routine to read the current value from a Value Block.
     *
     * Only for MIFARE Classic and only for blocks in "value block" mode, that
     * is: with access bits [C1 C2 C3] = [110] or [001]. The sector containing
     * the block must be authenticated before calling this function.
     *
     * \param[in]   blockAddr   The block (0x00-0xff) number.
     * \param[out]  value       Current value of the Value Block.
     *
     * \return StatusOk on success, STATUS_??? otherwise.
     */
    [[maybe_unused]] [[nodiscard]] uint8_t mifareGetValue(uint8_t blockAddr, long *value) const;

    /**
     * Helper routine to write a specific value into a Value Block.
     *
     * Only for MIFARE Classic and only for blocks in "value block" mode, that
     * is: with access bits [C1 C2 C3] = [110] or [001]. The sector containing
     * the block must be authenticated before calling this function.
     *
     * \param[in]   blockAddr   The block (0x00-0xff) number.
     * \param[in]   value       New value of the Value Block.
     *
     * \return StatusOk on success, STATUS_??? otherwise.
     */
    [[maybe_unused]] [[nodiscard]] uint8_t mifareSetValue(uint8_t blockAddr, long value) const;

    //-----------------------------------------------------------------------------------
    // Support functions
    //-----------------------------------------------------------------------------------

    /**
     * Wrapper for MIFARE protocol communication.
     * Adds CRC_A, executes the Transceive command and checks that the response is MfAck or a timeout.
     *
     * \param sendData Pointer to the data to transfer to the FIFO. Do NOT include the CRC_A.
     * \param sendLen Number of bytes in sendData.
     * \param acceptTimeout True => A timeout is also success
     *
     * @return StatusOk on success, STATUS_??? otherwise.
     */
    uint8_t pcdMifareTransceive(uint8_t *sendData,
                                uint8_t sendLen,
                                bool acceptTimeout = false) const;

    /**
     * Returns a __FlashStringHelper pointer to a status code name.
     *
     * \param code One of the StatusCode enums.
     */
    static std::string getStatusCodeName(uint8_t code);

    /**
     * Translates the SAK (Select Acknowledge) to a PICC type.
     *
     * \param sak The SAK uint8_t returned from piccSelect().
     *
     * \return PiccType
     */
    static uint8_t piccGetType(uint8_t sak);

    /**
     * Returns a String pointer to the PICC type name.
     *
     * \param type One of the PiccType enums.
     */
    static std::string piccGetTypeName(uint8_t type);

    /**
     * Dumps debug info about the selected PICC to Serial.
     * On success the PICC is halted after dumping the data.
     * For MIFARE Classic the factory default key of 0xFFFFFFFFFFFF is tried.
     *
     * \param uid Pointer to Uid struct returned from a successful piccSelect().
     */
    [[maybe_unused]] void piccDumpToSerial(Uid *uid) const;

    /**
     * Dumps memory contents of a MIFARE Classic PICC.
     * On success the PICC is halted after dumping the data.
     *
     * \param uid Pointer to Uid struct returned from a successful piccSelect().
     * \param piccType One of the PiccType enums.
     * \param key Key A used for all sectors.
     */
    void piccDumpMifareClassicToSerial(Uid *uid, uint8_t piccType, MifareKey *key) const;

    /**
     * Dumps memory contents of a sector of a MIFARE Classic PICC.
     * Uses pcdAuthenticate(), mifareRead() and pcdStopCrypto1.
     * Always uses PiccCmdMfAuthKeyA because only Key A can always read the sector trailer access bits.
     *
     * \param uid Pointer to Uid struct returned from a successful piccSelect().
     * \param key Key A for the sector.
     * \param sector The sector to dump, 0..39.
     */
    void piccDumpMifareClassicSectorToSerial(Uid *uid, MifareKey *key, uint8_t sector) const;

    /**
     * Dumps memory contents of a MIFARE Ultralight PICC.
     */
    void piccDumpMifareUltralightToSerial() const;

    /**
     * Calculates the bit pattern needed for the specified access bits. In the [C1 C2 C3] tupples C1 is MSB (=4) and C3 is LSB (=1).
     *
     * \param accessBitBuffer Pointer to uint8_t 6, 7 and 8 in the sector trailer. bytes [0..2] will be set.
     * \param g0 Access bits [C1 C2 C3] for block 0 (for sectors 0-31) or blocks 0-4 (for sectors 32-39)
     * \param g1 Access bits C1 C2 C3] for block 1 (for sectors 0-31) or blocks 5-9 (for sectors 32-39)
     * \param g2 Access bits C1 C2 C3] for block 2 (for sectors 0-31) or blocks 10-14 (for sectors 32-39)
     * \param g3 Access bits C1 C2 C3] for the sector trailer, block 3 (for sectors 0-31) or block 15 (for sectors 32-39)
     */
    [[maybe_unused]] static void mifareSetAccessBits(
        uint8_t *accessBitBuffer, uint8_t g0, uint8_t g1, uint8_t g2, uint8_t g3);

#pragma clang diagnostic push
#pragma ide diagnostic ignored "modernize-use-nodiscard"
    /**
     * Performs the "magic sequence" needed to get Chinese UID changeable
     * Mifare cards to allow writing to sector 0, where the card UID is stored.
     *
     * Note that you do not need to have selected the card through REQA or WUPA,
     * this sequence works immediately when the card is in the reader vicinity.
     * This means you can use this method even on "bricked" cards that your reader does
     * not recognise anymore (see Device::mifareUnbrickUidSector).
     *
     * Of course with non-bricked devices, you're free to select them before calling this function.
     */
    bool mifareOpenUidBackdoor(bool logErrors) const;
#pragma clang diagnostic pop

    /**
     * Reads entire block 0, including all manufacturer data, and overwrites
     * that block with the new UID, a freshly calculated BCC, and the original
     * manufacturer data.
     *
     * It assumes a default KEY A of 0xFFFFFFFFFFFF.
     * Make sure to have selected the card before this function is called.
     */
    [[maybe_unused]] bool mifareSetUid(const uint8_t *newUid, uint8_t uidSize, bool logErrors);

    /**
     * Resets entire sector 0 to zeroes, so the card can be read again by readers.
     */
    [[maybe_unused]] [[nodiscard]] bool mifareUnbrickUidSector(bool logErrors) const;

    //-----------------------------------------------------------------------------------
    // Convenience functions - do not add extra functionality
    //-----------------------------------------------------------------------------------

    /**
     * Returns true if a PICC responds to PiccCmdReqa.
     * Only "new" cards in state IDLE are invited. Sleeping cards in state HALT are ignored.
     *
     * @return Whether the new card is present or not.
     */
    [[nodiscard]] bool piccIsNewCardPresent() const;

    /**
     * Simple wrapper around piccSelect.
     * Returns true if a UID could be read.
     * Remember to call piccIsNewCardPresent(), piccRequestA() or piccWakeupA() first.
     * The read UID is available in the class member m_uid.
     *
     * @return True if the UID could be read.
     */
    bool piccReadCardSerial();

    Uid getUid();

private:
    /**
     * Helper function for the two-step MIFARE Classic protocol operations Decrement, Increment and Restore.
     *
     * \param command The command to use
     * \param blockAddr The block (0-0xff) number.
     * \param data The data to transfer in step 2
     *
     * \return StatusOk on success, STATUS_??? otherwise.
     */
    [[nodiscard]] uint8_t mifareTwoStepHelper(uint8_t command, uint8_t blockAddr, long data) const;

    /**
     * Sleep for ms milliseconds
     */
    static void delayMS(int ms);

    /// Used by piccReadCardSerial().
    Uid m_uid{};

    /// The SPI device, used for communicating with the MFRC522.
    Spi::Device m_spiDev;
};

}// namespace Mfrc522
