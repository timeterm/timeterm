#pragma once

#include <fcntl.h>
#include <mutex>
#include <unistd.h>
#include <unordered_set>

namespace Mfrc522
{// NOLINT

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
void unexportAllPins();
void unexportPin(uint8_t pin);
uint8_t readPin(uint8_t pin);

}// namespace Mfrc522::Gpio