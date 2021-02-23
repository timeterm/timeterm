cmake_minimum_required(VERSION 3.13)
project(mfrc522)

set(CMAKE_CXX_STANDARD 17)

find_package(PkgConfig REQUIRED)
pkg_check_modules(libgpiod REQUIRED IMPORTED_TARGET libgpiod libgpiodcxx)

add_library(mfrc522
	src/spi.cpp
	include/mfrc522/spi.h
	src/mfrc522.cpp
	include/mfrc522/mfrc522.h)
target_include_directories(mfrc522 PUBLIC include)
target_link_libraries(mfrc522 PUBLIC PkgConfig::libgpiod)

add_executable(mfrc522_read bin/main.cpp)
target_link_libraries(mfrc522_read mfrc522)
