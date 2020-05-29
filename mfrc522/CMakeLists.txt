cmake_minimum_required(VERSION 3.16)
project(mfrc522)

set(CMAKE_CXX_STANDARD 17)

add_library(mfrc522 src/spi.cpp include/mfrc522/spi.h src/mfrc522.cpp include/mfrc522/mfrc522.h include/mfrc522/gpio.h)
target_include_directories(mfrc522 PUBLIC include)