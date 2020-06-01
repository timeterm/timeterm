#pragma once

#include <fcntl.h>
#include <mutex>
#include <unistd.h>
#include <unordered_set>

namespace Mfrc522// NOLINT
{

//! The Gpio namespace.
namespace Gpio
{
}

}// namespace Mfrc522

namespace Mfrc522::Gpio
{

enum class PinDirection
{
    Out,
    In,
};

std::string_view pinDirectionToStringView(PinDirection direction);

class GlobalManager
{
public:
    static GlobalManager &singleton();
    ~GlobalManager();

    void exportPin(uint8_t pin, PinDirection direction);
    void writePin(uint8_t pin, uint8_t value);
    uint8_t readPin(uint8_t pin);
    void unexportPin(uint8_t pin);
    void unexportAllPins();

private:
    GlobalManager() = default;

public:
    GlobalManager(GlobalManager const &) = delete;
    void operator=(GlobalManager const &) = delete;

private:
    std::mutex m_mtx;
    std::unordered_set<uint8_t> m_exportedPins;
};

void exportPin(uint8_t pin, PinDirection direction);
void writePin(uint8_t pin, uint8_t value);
[[maybe_unused]] void unexportAllPins();
[[maybe_unused]] void unexportPin(uint8_t pin);
uint8_t readPin(uint8_t pin);

class InvalidPinDirectionException: public std::runtime_error
{
public:
    InvalidPinDirectionException()
        : std::runtime_error("invalid PinDirection")
    {}
};

class UnexportedPinWriteException: public std::runtime_error
{
public:
    UnexportedPinWriteException()
        : std::runtime_error("write to unexported pin")
    {}
};

class UnexportedPinReadException: public std::runtime_error
{
public:
    UnexportedPinReadException()
        : std::runtime_error("read of unexported pin")
    {}
};

class InvalidPinValueException: public std::runtime_error
{
public:
    InvalidPinValueException()
        : std::runtime_error("invalid pin value")
    {}
};

class PinOpenException: public std::runtime_error
{
public:
    PinOpenException()
        : std::runtime_error("could not open pin")
    {}
};

class DirectionOpenException: public std::runtime_error
{
public:
    DirectionOpenException()
        : std::runtime_error("could not open direction")
    {}
};

class PinDirectionSetException: public std::runtime_error
{
public:
    PinDirectionSetException()
        : std::runtime_error("could not set pin direction")
    {}
};

class PinValueSetException: public std::runtime_error
{
public:
    PinValueSetException()
        : std::runtime_error("could not set pin value")
    {}
};

class PinExportException: public std::runtime_error
{
public:
    PinExportException()
        : std::runtime_error("could not export pin")
    {}
};

class PinUnexportException: public std::runtime_error
{
public:
    PinUnexportException()
        : std::runtime_error("could not unexport pin")
    {}
};

class ExportOpenException: public std::runtime_error
{
public:
    ExportOpenException()
        : std::runtime_error("could not open pin export")
    {}
};

class UnexportOpenException: public std::runtime_error
{
public:
    UnexportOpenException()
        : std::runtime_error("could not open pin unexport")
    {}
};

}// namespace Mfrc522::Gpio