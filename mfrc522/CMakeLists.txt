cmake_minimum_required(VERSION 3.16)
project(mfrc522)

set(CMAKE_CXX_STANDARD 17)

add_library(mfrc522 src/spi.cpp include/mfrc522/spi.h src/mfrc522.cpp include/mfrc522/mfrc522.h include/mfrc522/gpio.h src/gpio.cpp)
target_include_directories(mfrc522 PUBLIC include)

add_executable(mfrc522_read bin/main.cpp)
target_link_libraries(mfrc522_read mfrc522)