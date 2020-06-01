cmake_minimum_required(VERSION 3.13)
project(mfrc522_2)

set(CMAKE_CXX_STANDARD 17)

add_executable(mfrc522_2 mfrc522.cpp mfrc522.h main.cpp spi.h spi.cpp gpio.cpp gpio.h)