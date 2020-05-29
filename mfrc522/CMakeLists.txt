cmake_minimum_required(VERSION 3.16)
project(mfrc522)

set(CMAKE_CXX_STANDARD 17)

add_library(mfrc522 library.cpp library.h)