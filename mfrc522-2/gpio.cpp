#include "gpio.h"
#include <iostream>
#include <vector>

namespace Mfrc522::Gpio
{

void _exportPin(uint8_t pin)
{
    int fd = open("/sys/class/gpio/export", O_WRONLY);
    if (fd == -1) {
        // TODO(rutgerbrf): custom exception
        throw std::runtime_error("could not export pin (not enough permissions?)");
    }

    auto pinStr = std::to_string(pin);
    if (write(fd, pinStr.c_str(), pinStr.length()) != pinStr.length()) {
        // TODO(rutgerbrf): custom exception
        throw std::runtime_error("could not export pin (error writing pin export)");
    }

    close(fd);
}

void _setPinDirection(uint8_t pin, PinDirection direction)
{
    auto pinStr = std::to_string(pin);
    auto path = "/sys/class/gpio/gpio" + pinStr + "/direction";

    int fd = open(path.c_str(), O_WRONLY);
    if (fd == -1) {
        // TODO(rutgerbrf): custom exception
        throw std::runtime_error("could not set pin direction (not enough permissions?)");
    }

    auto directionStr = pinDirectionToStringView(direction);
    if (write(fd, directionStr.data(), directionStr.length()) != directionStr.length()) {
        // TODO(rutgerbrf): custom exception
        throw std::runtime_error("could not set pin direction (error writing pin direction)");
    }
}

void _unexportPin(uint8_t pin)
{
    int fd = open("/sys/class/gpio/unexport", O_WRONLY);
    if (fd == -1) {
        // TODO(rutgerbrf): custom exception
        throw std::runtime_error("could not unexport pin (not enough permissions?)");
    }

    auto pinStr = std::to_string(pin);
    if (write(fd, pinStr.c_str(), pinStr.length()) != pinStr.length()) {
        // TODO(rutgerbrf): custom exception
        throw std::runtime_error("could not unexport pin (error writing pin unexport)");
    }

    close(fd);
}

void _writePin(uint8_t pin, uint8_t value)
{
    auto pinStr = std::to_string(pin);
    auto path = "/sys/class/gpio/gpio" + pinStr + "/value";

    int fd = open(path.c_str(), O_WRONLY);
    if (fd == -1) {
        // TODO(rutgerbrf): custom exception
        throw std::runtime_error("could not set pin value (not enough permissions?)");
    }

    auto valueStr = std::to_string(value);
    if (write(fd, valueStr.c_str(), valueStr.length()) != valueStr.length()) {
        // TODO(rutgerbrf): custom exception
        throw std::runtime_error("could not set pin value (error writing pin value)");
    }

    close(fd);
}

uint8_t _readPin(uint8_t pin)
{
    auto pinStr = std::to_string(pin);
    auto path = "/sys/class/gpio/gpio" + pinStr + "/value";

    int fd = open(path.c_str(), O_WRONLY);
    if (fd == -1) {
        // TODO(rutgerbrf): custom exception
        throw std::runtime_error("could not set pin value (not enough permissions?)");
    }

    char bytes[4] = {0};
    read(fd, bytes, 4);
    auto byte = atoi(bytes);
    if (byte > UINT8_MAX) {
        throw std::runtime_error("invalid pin value");
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
        }
        catch (...) {
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
        // TODO(rutgerbrf): custom exception
        throw std::runtime_error("write to unexported pin");
    }

    _writePin(pin, value);
}

uint8_t GlobalManager::readPin(uint8_t pin) {
    auto guard = std::lock_guard{m_mtx};

    if (m_exportedPins.find(pin) == m_exportedPins.end()) {
        // TODO(rutgerbrf): custom exception
        throw std::runtime_error("write to unexported pin");
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
            // TODO(rutgerbrf): custom exception
            throw std::runtime_error("invalid PinDirection");
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

void unexportAllPins() {
    GlobalManager::singleton().unexportAllPins();
}

void unexportPin(uint8_t pin) {
    GlobalManager::singleton().unexportPin(pin);
}

uint8_t readPin(uint8_t pin) {
    return GlobalManager::singleton().readPin(pin);
}

}// namespace Mfrc522::Gpio