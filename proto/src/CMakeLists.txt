cmake_minimum_required(VERSION 3.5)
project(timeterm_proto)

set(CMAKE_VERBOSE_MAKEFILE TRUE)

if (NOT MSVC)
    set(CMAKE_CXX_FLAGS "${CMAKE_CXX_FLAGS} -std=c++11")
else ()
    add_definitions(-D_WIN32_WINNT=0x600)
endif ()

# This branch assumes that Protobuf is already installed
# on this system, so it can be located by find_package().

# Find Protobuf installation
# Looks for protobuf-config.cmake file installed by Protobuf's cmake installation.
find_package(Protobuf REQUIRED)
message(STATUS "Using protobuf ${protobuf_VERSION}")

protobuf_generate_cpp(PROTO_SRCS PROTO_HDRS timeterm_proto/messages.proto)

file(MAKE_DIRECTORY ${CMAKE_CURRENT_BINARY_DIR}/timeterm_proto)
file(GLOB GENERATED_FILES ${CMAKE_CURRENT_BINARY_DIR}/*.pb.*)
foreach(GENERATED_FILE IN LISTS GENERATED_FILES)
	string(REGEX REPLACE "^${CMAKE_CURRENT_BINARY_DIR}" "" GENERATED_FILE_RELATIVE ${GENERATED_FILE})
	file(RENAME ${GENERATED_FILE} ${CMAKE_CURRENT_BINARY_DIR}/timeterm_proto${GENERATED_FILE_RELATIVE})
endforeach()

add_library(timeterm_proto ${PROTO_SRCS} ${PROTO_HDRS})
target_link_libraries(timeterm_proto PUBLIC ${Protobuf_LIBRARIES})
target_include_directories(timeterm_proto PUBLIC ${CMAKE_CURRENT_BINARY_DIR})
target_include_directories(timeterm_proto PRIVATE ${Protobuf_INCLUDE_DIRS})
