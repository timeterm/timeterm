cmake_minimum_required(VERSION 3.5)
project(timeterm_proto)

set(CMAKE_VERBOSE_MAKEFILE TRUE)
set(CMAKE_CXX_STANDARD 11)

cmake_policy(SET CMP0071 NEW)

if(NOT MSVC)
#    set(CMAKE_CXX_FLAGS "${CMAKE_CXX_FLAGS} -std=c++11")
else()
    add_definitions(-D_WIN32_WINNT=0x600)
endif()

set(Protobuf_DEBUG ON)
if(DEFINED VCPKG_TARGET_TRIPLET)
	message(STATUS "timeterm/proto: Vcpkg detected, using static libprotobuf")
	set(Protobuf_USE_STATIC_LIBS ON)
	find_package(protobuf CONFIG REQUIRED)
else()
	find_package(Protobuf REQUIRED)
endif()

message(STATUS "timeterm/proto: using libprotobuf + protoc ${protobuf_VERSION}")

if(MSVC AND protobuf_MSVC_STATIC_RUNTIME)
	foreach(flag_var
			CMAKE_CXX_FLAGS CMAKE_CXX_FLAGS_DEBUG CMAKE_CXX_FLAGS_RELEASE
			CMAKE_CXX_FLAGS_MINSIZEREL CMAKE_CXX_FLAGS_RELWITHDEBINFO)
		if(${flag_var} MATCHES "/MD")
			string(REGEX REPLACE "/MD" "/MT" ${flag_var} "${${flag_var}}")
		endif(${flag_var} MATCHES "/MD")
	endforeach()
endif()

add_library(timeterm_proto STATIC)

protobuf_generate(
		LANGUAGE cpp
		TARGET timeterm_proto
		PROTOS timeterm_proto/devcfg/devcfg.proto timeterm_proto/mq/mq.proto
)

target_link_libraries(timeterm_proto PUBLIC protobuf::libprotobuf)
target_include_directories(timeterm_proto PUBLIC ${CMAKE_CURRENT_BINARY_DIR})
