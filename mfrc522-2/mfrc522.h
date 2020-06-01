/**
 * MFRC522.h - Library to use ARDUINO RFID MODULE KIT 13.56 MHZ WITH TAGS SPI W AND R BY COOQROBOT.
 * Based on code Dr.Leong   ( WWW.B2CQSHOP.COM )
 * Created by Miguel Balboa (circuitito.com), Jan, 2012.
 * Rewritten by Søren Thing Andersen (access.thing.dk), fall of 2013 (Translation to English, refactored, comments, anti collision, cascade levels.)
 * Extended by Tom Clement with functionality to write to sector 0 of UID changeable Mifare cards.
 * Released into the public domain.
 *


-- Repurposed to fit Raspberry Pi ---

 */
#ifndef MFRC522_h
#define MFRC522_h

#include <cstdint>
#include <cstdio>
#include <string>
#include "spi.h"

using namespace std;
using namespace Mfrc522;

typedef uint16_t word;

// Firmware data for self-test
// Reference values based on firmware version; taken from 16.1.1 in spec.
// Version 1.0

const std::uint8_t MFRC522_firmware_referenceV1_0[]  = {
    0x00, 0xC6, 0x37, 0xD5, 0x32, 0xB7, 0x57, 0x5C,
    0xC2, 0xD8, 0x7C, 0x4D, 0xD9, 0x70, 0xC7, 0x73,
    0x10, 0xE6, 0xD2, 0xAA, 0x5E, 0xA1, 0x3E, 0x5A,
    0x14, 0xAF, 0x30, 0x61, 0xC9, 0x70, 0xDB, 0x2E,
    0x64, 0x22, 0x72, 0xB5, 0xBD, 0x65, 0xF4, 0xEC,
    0x22, 0xBC, 0xD3, 0x72, 0x35, 0xCD, 0xAA, 0x41,
    0x1F, 0xA7, 0xF3, 0x53, 0x14, 0xDE, 0x7E, 0x02,
    0xD9, 0x0F, 0xB5, 0x5E, 0x25, 0x1D, 0x29, 0x79
};

// Version 2.0
const std::uint8_t MFRC522_firmware_referenceV2_0[] = {
    0x00, 0xEB, 0x66, 0xBA, 0x57, 0xBF, 0x23, 0x95,
    0xD0, 0xE3, 0x0D, 0x3D, 0x27, 0x89, 0x5C, 0xDE,
    0x9D, 0x3B, 0xA7, 0x00, 0x21, 0x5B, 0x89, 0x82,
    0x51, 0x3A, 0xEB, 0x02, 0x0C, 0xA5, 0x00, 0x49,
    0x7C, 0x84, 0x4D, 0xB3, 0xCC, 0xD2, 0x1B, 0x81,
    0x5D, 0x48, 0x76, 0xD5, 0x71, 0x61, 0x21, 0xA9,
    0x86, 0x96, 0x83, 0x38, 0xCF, 0x9D, 0x5B, 0x6D,
    0xDC, 0x15, 0xBA, 0x3E, 0x7D, 0x95, 0x3B, 0x2F
};



class MFRC522 {
public:
    // MFRC522 registers. Described in chapter 9 of the datasheet.
    // When using SPI all addresses are shifted one bit left in the "SPI address uint8_t" (section 8.1.2.3)
    enum PCD_Register {
        // Page 0: Command and status
        //						  0x00			// reserved for future use
        CommandReg				= 0x01u << 1u,	// starts and stops command execution
        ComIEnReg				= 0x02u << 1u,	// enable and disable interrupt request control bits
        DivIEnReg				= 0x03u << 1u,	// enable and disable interrupt request control bits
        ComIrqReg				= 0x04u << 1u,	// interrupt request bits
        DivIrqReg				= 0x05u << 1u,	// interrupt request bits
        ErrorReg				= 0x06u << 1u,	// error bits showing the error status of the last command executed
        Status1Reg				= 0x07u << 1u,	// communication status bits
        Status2Reg				= 0x08u << 1u,	// receiver and transmitter status bits
        FIFODataReg				= 0x09u << 1u,	// input and output of 64 uint8_t FIFO buffer
        FIFOLevelReg			= 0x0Au << 1u,	// number of uint8_ts stored in the FIFO buffer
        WaterLevelReg			= 0x0Bu << 1u,	// level for FIFO underflow and overflow warning
        ControlReg				= 0x0Cu << 1u,	// miscellaneous control registers
        BitFramingReg			= 0x0Du << 1u,	// adjustments for bit-oriented frames
        CollReg					= 0x0Eu << 1u,	// bit position of the first bit-collision detected on the RF interface
        //						  0x0F			// reserved for future use

        // Page 1: Command
        // 						  0x10			// reserved for future use
        ModeReg					= 0x11u << 1u,	// defines general modes for transmitting and receiving
        TxModeReg				= 0x12u << 1u,	// defines transmission data rate and framing
        RxModeReg				= 0x13u << 1u,	// defines reception data rate and framing
        TxControlReg			= 0x14u << 1u,	// controls the logical behavior of the antenna driver pins TX1 and TX2
        TxASKReg				= 0x15u << 1u,	// controls the setting of the transmission modulation
        TxSelReg				= 0x16u << 1u,	// selects the internal sources for the antenna driver
        RxSelReg				= 0x17u << 1u,	// selects internal receiver settings
        RxThresholdReg			= 0x18u << 1u,	// selects thresholds for the bit decoder
        DemodReg				= 0x19u << 1u,	// defines demodulator settings
        // 						  0x1A			// reserved for future use
        // 						  0x1B			// reserved for future use
        MfTxReg					= 0x1Cu << 1u,	// controls some MIFARE communication transmit parameters
        MfRxReg					= 0x1Du << 1u,	// controls some MIFARE communication receive parameters
        // 						  0x1E			// reserved for future use
        SerialSpeedReg			= 0x1Fu << 1u,	// selects the speed of the serial UART interface

        // Page 2: Configuration
        // 						  0x20			// reserved for future use
        CRCResultRegH			= 0x21u << 1u,	// shows the MSB and LSB values of the CRC calculation
        CRCResultRegL			= 0x22u << 1u,
        // 						  0x23			// reserved for future use
        ModWidthReg				= 0x24u << 1u,	// controls the ModWidth setting?
        // 						  0x25			// reserved for future use
        RFCfgReg				= 0x26u << 1u,	// configures the receiver gain
        GsNReg					= 0x27u << 1u,	// selects the conductance of the antenna driver pins TX1 and TX2 for modulation
        CWGsPReg				= 0x28u << 1u,	// defines the conductance of the p-driver output during periods of no modulation
        ModGsPReg				= 0x29u << 1u,	// defines the conductance of the p-driver output during periods of modulation
        TModeReg				= 0x2Au << 1u,	// defines settings for the internal timer
        TPrescalerReg			= 0x2Bu << 1u,	// the lower 8 bits of the TPrescaler value. The 4 high bits are in TModeReg.
        TReloadRegH				= 0x2Cu << 1u,	// defines the 16-bit timer reload value
        TReloadRegL				= 0x2Du << 1u,
        TCounterValueRegH		= 0x2Eu << 1u,	// shows the 16-bit timer value
        TCounterValueRegL		= 0x2Fu << 1u,

        // Page 3: Test Registers
        // 						  0x30			// reserved for future use
        TestSel1Reg				= 0x31u << 1u,	// general test signal configuration
        TestSel2Reg				= 0x32u << 1u,	// general test signal configuration
        TestPinEnReg			= 0x33u << 1u,	// enables pin output driver on pins D1 to D7
        TestPinValueReg			= 0x34u << 1u,	// defines the values for D1 to D7 when it is used as an I/O bus
        TestBusReg				= 0x35u << 1u,	// shows the status of the internal test bus
        AutoTestReg				= 0x36u << 1u,	// controls the digital self test
        VersionReg				= 0x37u << 1u,	// shows the software version
        AnalogTestReg			= 0x38u << 1u,	// controls the pins AUX1 and AUX2
        TestDAC1Reg				= 0x39u << 1u,	// defines the test value for TestDAC1
        TestDAC2Reg				= 0x3Au << 1u,	// defines the test value for TestDAC2
        TestADCReg				= 0x3Bu << 1u		// shows the value of ADC I and Q channels
        // 						  0x3C			// reserved for production tests
        // 						  0x3D			// reserved for production tests
        // 						  0x3E			// reserved for production tests
        // 						  0x3F			// reserved for production tests
    };

    // MFRC522 commands. Described in chapter 10 of the datasheet.
    enum PCD_Command {
        PCD_Idle				= 0x00,		// no action, cancels current command execution
        PCD_Mem					= 0x01,		// stores 25 uint8_ts into the internal buffer
        PCD_GenerateRandomID	= 0x02,		// generates a 10-uint8_t random ID number
        PCD_CalcCRC				= 0x03,		// activates the CRC coprocessor or performs a self test
        PCD_Transmit			= 0x04,		// transmits data from the FIFO buffer
        PCD_NoCmdChange			= 0x07,		// no command change, can be used to modify the CommandReg register bits without affecting the command, for example, the PowerDown bit
        PCD_Receive				= 0x08,		// activates the receiver circuits
        PCD_Transceive 			= 0x0C,		// transmits data from FIFO buffer to antenna and automatically activates the receiver after transmission
        PCD_MFAuthent 			= 0x0E,		// performs the MIFARE standard authentication as a reader
        PCD_SoftReset			= 0x0F		// resets the MFRC522
    };

    // MFRC522 RxGain[2:0] masks, defines the receiver's signal voltage gain factor (on the PCD).
    // Described in 9.3.3.6 / table 98 of the datasheet at http://www.nxp.com/documents/data_sheet/MFRC522.pdf
    enum PCD_RxGain {
        RxGain_18dB				= 0x00u << 4u,	// 000b - 18 dB, minimum
        RxGain_23dB				= 0x01u << 4u,	// 001b - 23 dB
        RxGain_18dB_2			= 0x02u << 4u,	// 010b - 18 dB, it seems 010b is a duplicate for 000b
        RxGain_23dB_2			= 0x03u << 4u,	// 011b - 23 dB, it seems 011b is a duplicate for 001b
        RxGain_33dB				= 0x04u << 4u,	// 100b - 33 dB, average, and typical default
        RxGain_38dB				= 0x05u << 4u,	// 101b - 38 dB
        RxGain_43dB				= 0x06u << 4u,	// 110b - 43 dB
        RxGain_48dB				= 0x07u << 4u,	// 111b - 48 dB, maximum
        RxGain_min				= 0x00u << 4u,	// 000b - 18 dB, minimum, convenience for RxGain_18dB
        RxGain_avg				= 0x04u << 4u,	// 100b - 33 dB, average, convenience for RxGain_33dB
        RxGain_max				= 0x07u << 4u		// 111b - 48 dB, maximum, convenience for RxGain_48dB
    };

    // Commands sent to the PICC.
    enum PICC_Command {
        // The commands used by the PCD to manage communication with several PICCs (ISO 14443-3, Type A, section 6.4)
        PICC_CMD_REQA			= 0x26,		// REQuest command, Type A. Invites PICCs in state IDLE to go to READY and prepare for anticollision or selection. 7 bit frame.
        PICC_CMD_WUPA			= 0x52,		// Wake-UP command, Type A. Invites PICCs in state IDLE and HALT to go to READY(*) and prepare for anticollision or selection. 7 bit frame.
        PICC_CMD_CT				= 0x88,		// Cascade Tag. Not really a command, but used during anti collision.
        PICC_CMD_SEL_CL1		= 0x93,		// Anti collision/Select, Cascade Level 1
        PICC_CMD_SEL_CL2		= 0x95,		// Anti collision/Select, Cascade Level 2
        PICC_CMD_SEL_CL3		= 0x97,		// Anti collision/Select, Cascade Level 3
        PICC_CMD_HLTA			= 0x50,		// HaLT command, Type A. Instructs an ACTIVE PICC to go to state HALT.
        // The commands used for MIFARE Classic (from http://www.nxp.com/documents/data_sheet/MF1S503x.pdf, Section 9)
        // Use PCD_MFAuthent to authenticate access to a sector, then use these commands to read/write/modify the blocks on the sector.
        // The read/write commands can also be used for MIFARE Ultralight.
        PICC_CMD_MF_AUTH_KEY_A	= 0x60,		// Perform authentication with Key A
        PICC_CMD_MF_AUTH_KEY_B	= 0x61,		// Perform authentication with Key B
        PICC_CMD_MF_READ		= 0x30,		// Reads one 16 uint8_t block from the authenticated sector of the PICC. Also used for MIFARE Ultralight.
        PICC_CMD_MF_WRITE		= 0xA0,		// Writes one 16 uint8_t block to the authenticated sector of the PICC. Called "COMPATIBILITY WRITE" for MIFARE Ultralight.
        PICC_CMD_MF_DECREMENT	= 0xC0,		// Decrements the contents of a block and stores the result in the internal data register.
        PICC_CMD_MF_INCREMENT	= 0xC1,		// Increments the contents of a block and stores the result in the internal data register.
        PICC_CMD_MF_RESTORE		= 0xC2,		// Reads the contents of a block into the internal data register.
        PICC_CMD_MF_TRANSFER	= 0xB0,		// Writes the contents of the internal data register to a block.
        // The commands used for MIFARE Ultralight (from http://www.nxp.com/documents/data_sheet/MF0ICU1.pdf, Section 8.6)
        // The PICC_CMD_MF_READ and PICC_CMD_MF_WRITE can also be used for MIFARE Ultralight.
        PICC_CMD_UL_WRITE		= 0xA2		// Writes one 4 uint8_t page to the PICC.
    };

    // MIFARE constants that does not fit anywhere else
    enum MIFARE_Misc {
        MF_ACK					= 0xA,		// The MIFARE Classic uses a 4 bit ACK/NAK. Any other value than 0xA is NAK.
        MF_KEY_SIZE				= 6			// A Mifare Crypto1 key is 6 uint8_ts.
    };

    // PICC types we can detect. Remember to update PICC_GetTypeName() if you add more.
    enum PICC_Type {
        PICC_TYPE_UNKNOWN		= 0,
        PICC_TYPE_ISO_14443_4	= 1,	// PICC compliant with ISO/IEC 14443-4
        PICC_TYPE_ISO_18092		= 2, 	// PICC compliant with ISO/IEC 18092 (NFC)
        PICC_TYPE_MIFARE_MINI	= 3,	// MIFARE Classic protocol, 320 uint8_ts
        PICC_TYPE_MIFARE_1K		= 4,	// MIFARE Classic protocol, 1KB
        PICC_TYPE_MIFARE_4K		= 5,	// MIFARE Classic protocol, 4KB
        PICC_TYPE_MIFARE_UL		= 6,	// MIFARE Ultralight or Ultralight C
        PICC_TYPE_MIFARE_PLUS	= 7,	// MIFARE Plus
        PICC_TYPE_TNP3XXX		= 8,	// Only mentioned in NXP AN 10833 MIFARE Type Identification Procedure
        PICC_TYPE_NOT_COMPLETE	= 255	// SAK indicates UID is not complete.
    };

    // Return codes from the functions in this class. Remember to update GetStatusCodeName() if you add more.
    enum StatusCode {
        STATUS_OK				= 1,	// Success
        STATUS_ERROR			= 2,	// Error in communication
        STATUS_COLLISION		= 3,	// Collission detected
        STATUS_TIMEOUT			= 4,	// Timeout in communication.
        STATUS_NO_ROOM			= 5,	// A buffer is not big enough.
        STATUS_INTERNAL_ERROR	= 6,	// Internal error in the code. Should not happen ;-)
        STATUS_INVALID			= 7,	// Invalid argument.
        STATUS_CRC_WRONG		= 8,	// The CRC_A does not match
        STATUS_MIFARE_NACK		= 9		// A MIFARE PICC responded with NAK.
    };

    // A struct used for passing the UID of a PICC.
    typedef struct {
        uint8_t		size;			// Number of uint8_ts in the UID. 4, 7 or 10.
        uint8_t		uidByte[10];
        uint8_t		sak;			// The SAK (Select acknowledge) uint8_t returned from the PICC after successful selection.
    } Uid;

    // A struct used for passing a MIFARE Crypto1 key
    typedef struct {
        uint8_t		keyuint8_t[MF_KEY_SIZE];
    } MIFARE_Key;

    // Member variables
    Uid uid;								// Used by PICC_ReadCardSerial().
    Spi::Device spiDev;

    // Size of the MFRC522 FIFO
    static const uint8_t FIFO_SIZE = 64;		// The FIFO is 64 uint8_ts.

    /////////////////////////////////////////////////////////////////////////////////////
    // Functions for setting up the Raspberry Pi
    /////////////////////////////////////////////////////////////////////////////////////
    MFRC522();
    /////////////////////////////////////////////////////////////////////////////////////
    // Basic interface functions for communicating with the MFRC522
    /////////////////////////////////////////////////////////////////////////////////////
    void PCD_WriteRegister(uint8_t reg, uint8_t value);
    void PCD_WriteRegister(uint8_t reg, uint8_t count, uint8_t *values);
    uint8_t PCD_ReadRegister(uint8_t reg);
    void PCD_ReadRegister(uint8_t reg, uint8_t count, uint8_t *values, uint8_t rxAlign = 0);
    void PCD_SetRegisterBitMask(uint8_t reg, uint8_t mask);
    void PCD_ClearRegisterBitMask(uint8_t reg, uint8_t mask);
    uint8_t PCD_CalculateCRC(uint8_t *data, uint8_t length, uint8_t *result);

    /////////////////////////////////////////////////////////////////////////////////////
    // Functions for manipulating the MFRC522
    /////////////////////////////////////////////////////////////////////////////////////
    void PCD_Init();
    void PCD_Reset();
    void PCD_AntennaOn();
    void PCD_AntennaOff();
    uint8_t PCD_GetAntennaGain();
    void PCD_SetAntennaGain(uint8_t mask);
    bool PCD_PerformSelfTest();

    /////////////////////////////////////////////////////////////////////////////////////
    // Functions for communicating with PICCs
    /////////////////////////////////////////////////////////////////////////////////////
    uint8_t PCD_TransceiveData(uint8_t *sendData, uint8_t sendLen, uint8_t *backData, uint8_t *backLen, uint8_t *validBits = nullptr, uint8_t rxAlign = 0, bool checkCRC = false);
    uint8_t PCD_CommunicateWithPICC(uint8_t command, uint8_t waitIRq, uint8_t *sendData, uint8_t sendLen, uint8_t *backData = nullptr, uint8_t *backLen = nullptr, uint8_t *validBits = nullptr, uint8_t rxAlign = 0, bool checkCRC = false);
    uint8_t PICC_RequestA(uint8_t *bufferATQA, uint8_t *bufferSize);
    uint8_t PICC_WakeupA(uint8_t *bufferATQA, uint8_t *bufferSize);
    uint8_t PICC_REQA_or_WUPA(uint8_t command, uint8_t *bufferATQA, uint8_t *bufferSize);
    uint8_t PICC_Select(Uid *uid, uint8_t validBits = 0);
    uint8_t PICC_HaltA();

    /////////////////////////////////////////////////////////////////////////////////////
    // Functions for communicating with MIFARE PICCs
    /////////////////////////////////////////////////////////////////////////////////////
    uint8_t PCD_Authenticate(uint8_t command, uint8_t blockAddr, MIFARE_Key *key, Uid *uid);
    void PCD_StopCrypto1();
    uint8_t MIFARE_Read(uint8_t blockAddr, uint8_t *buffer, uint8_t *bufferSize);
    uint8_t MIFARE_Write(uint8_t blockAddr, uint8_t *buffer, uint8_t bufferSize);
    uint8_t MIFARE_Decrement(uint8_t blockAddr, long delta);
    uint8_t MIFARE_Increment(uint8_t blockAddr, long delta);
    uint8_t MIFARE_Restore(uint8_t blockAddr);
    uint8_t MIFARE_Transfer(uint8_t blockAddr);
    uint8_t MIFARE_Ultralight_Write(uint8_t page, uint8_t *buffer, uint8_t bufferSize);
    uint8_t MIFARE_GetValue(uint8_t blockAddr, long *value);
    uint8_t MIFARE_SetValue(uint8_t blockAddr, long value);

    /////////////////////////////////////////////////////////////////////////////////////
    // Support functions
    /////////////////////////////////////////////////////////////////////////////////////
    uint8_t PCD_MIFARE_Transceive(uint8_t *sendData, uint8_t sendLen, bool acceptTimeout = false);
    // old function used too much memory, now name moved to flash; if you need char, copy from flash to memory
    //const char *GetStatusCodeName(uint8_t code);
    const string GetStatusCodeName(uint8_t code);
    uint8_t PICC_GetType(uint8_t sak);
    // old function used too much memory, now name moved to flash; if you need char, copy from flash to memory
    //const char *PICC_GetTypeName(uint8_t type);
    const string PICC_GetTypeName(uint8_t type);
    void PICC_DumpToSerial(Uid *uid);
    void PICC_DumpMifareClassicToSerial(Uid *uid, uint8_t piccType, MIFARE_Key *key);
    void PICC_DumpMifareClassicSectorToSerial(Uid *uid, MIFARE_Key *key, uint8_t sector);
    void PICC_DumpMifareUltralightToSerial();
    void MIFARE_SetAccessBits(uint8_t *accessBitBuffer, uint8_t g0, uint8_t g1, uint8_t g2, uint8_t g3);
    bool MIFARE_OpenUidBackdoor(bool logErrors);
    bool MIFARE_SetUid(uint8_t *newUid, uint8_t uidSize, bool logErrors);
    bool MIFARE_UnbrickUidSector(bool logErrors);

    /////////////////////////////////////////////////////////////////////////////////////
    // Convenience functions - does not add extra functionality
    /////////////////////////////////////////////////////////////////////////////////////
    bool PICC_IsNewCardPresent();
    bool PICC_ReadCardSerial();

private:
    uint8_t MIFARE_TwoStepHelper(uint8_t command, uint8_t blockAddr, long data);
};

#endif
