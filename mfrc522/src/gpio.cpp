#include <iostream>
#include <mfrc522/gpio.h>
#include <vector>

namespace Mfrc522::Gpio {

void _exportPin(uint8_t pin)
{
    int fd = open("/sys/class/gpio/export", O_WRONLY);
    if (fd == -1) {
        throw ExportOpenException();
    }

    auto pinStr = std::to_string(pin);
    if (write(fd, pinStr.c_str(), pinStr.length()) != pinStr.length()) {
        throw PinExportException();
    }

    close(fd);
}

void _setPinDirection(uint8_t pin, PinDirection direction)
{
    auto pinStr = std::to_string(pin);
    auto path = "/sys/class/gpio/gpio" + pinStr + "/direction";

    int fd = open(path.c_str(), O_WRONLY);
    if (fd == -1) {
        throw DirectionOpenException();
    }

    auto directionStr = pinDirectionToStringView(direction);
    if (write(fd, directionStr.data(), directionStr.length()) != directionStr.length()) {
        throw PinDirectionSetException();
    }
}

void _unexportPin(uint8_t pin)
{
    int fd = open("/sys/class/gpio/unexport", O_WRONLY);
    if (fd == -1) {
        throw UnexportOpenException();
    }

    auto pinStr = std::to_string(pin);
    if (write(fd, pinStr.c_str(), pinStr.length()) != pinStr.length()) {
        throw PinUnexportException();
    }

    close(fd);
}

void _writePin(uint8_t pin, uint8_t value)
{
    auto pinStr = std::to_string(pin);
    auto path = "/sys/class/gpio/gpio" + pinStr + "/value";

    int fd = open(path.c_str(), O_WRONLY);
    if (fd == -1) {
        throw PinOpenException();
    }

    auto valueStr = std::to_string(value);
    if (write(fd, valueStr.c_str(), valueStr.length()) != valueStr.length()) {
        throw PinValueSetException();
    }

    close(fd);
}

uint8_t _readPin(uint8_t pin)
{
    auto pinStr = std::to_string(pin);
    auto path = "/sys/class/gpio/gpio" + pinStr + "/value";

    int fd = open(path.c_str(), O_WRONLY);
    if (fd == -1) {
        throw PinOpenException();
    }

    char bytes[4] = {0};
    read(fd, bytes, 4);
    auto byte = strtoul(bytes, nullptr, 10);
    if (byte > UINT8_MAX) {
        throw InvalidPinValueException();
    }

    close(fd);

    return byte;
}

GlobalManager &GlobalManager::singleton()
{
    static GlobalManager instance;
    return instance;
}

void GlobalManager::exportPin(uint8_t pin, PinDirection direction)
{
    auto guard = std::lock_guard{m_mtx};

    if (m_exportedPins.find(pin) != m_exportedPins.end()) {
        // Pin is already exported.
        return;
    }

    _exportPin(pin);
    _setPinDirection(pin, direction);

    m_exportedPins.insert(pin);
}

GlobalManager::~GlobalManager()
{
    if (std::current_exception()) {
        // An exception is currently propagating.
        try {
            unexportAllPins();
        } catch (...) {
            // We're currently in the destructor. In the case of an exception already propagating
            // we don't want the program to completely shut down due to another exception being
            // thrown, hence the catch-all.
        }

        return;
    }

    unexportAllPins();
}

void GlobalManager::unexportAllPins()
{
    auto guard = std::lock_guard{m_mtx};

    auto it = m_exportedPins.begin();
    while (it != m_exportedPins.end()) {
        _unexportPin(*it);
        it = m_exportedPins.erase(it);
    }
}

void GlobalManager::unexportPin(uint8_t pin)
{
    auto guard = std::lock_guard{m_mtx};

    auto it = m_exportedPins.find(pin);
    if (it == m_exportedPins.end()) {
        return;
    }

    _unexportPin(*it);
    m_exportedPins.erase(it);
}

void GlobalManager::writePin(uint8_t pin, uint8_t value)
{
    auto guard = std::lock_guard{m_mtx};

    if (m_exportedPins.find(pin) == m_exportedPins.end()) {
        throw UnexportedPinWriteException();
    }

    _writePin(pin, value);
}

uint8_t GlobalManager::readPin(uint8_t pin)
{
    auto guard = std::lock_guard{m_mtx};

    if (m_exportedPins.find(pin) == m_exportedPins.end()) {
        throw UnexportedPinReadException();
    }

    return _readPin(pin);
}

std::string_view pinDirectionToStringView(PinDirection direction)
{
    switch (direction) {
    case PinDirection::Out:
        return "out";
    case PinDirection::In:
        return "in";
    default:
        throw InvalidPinDirectionException();
    }
}

void exportPin(uint8_t pin, PinDirection direction)
{
    GlobalManager::singleton().exportPin(pin, direction);
}

void writePin(uint8_t pin, uint8_t value)
{
    GlobalManager::singleton().writePin(pin, value);
}

void unexportAllPins()
{
    GlobalManager::singleton().unexportAllPins();
}

void unexportPin(uint8_t pin)
{
    GlobalManager::singleton().unexportPin(pin);
}

uint8_t readPin(uint8_t pin)
{
    return GlobalManager::singleton().readPin(pin);
}

} // namespace Mfrc522::Gpio
