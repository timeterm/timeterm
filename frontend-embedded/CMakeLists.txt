cmake_minimum_required(VERSION 3.15)

project(frontend-embedded LANGUAGES CXX)

set(CMAKE_INCLUDE_CURRENT_DIR ON)

set(CMAKE_AUTOUIC ON)
set(CMAKE_AUTOMOC ON)
set(CMAKE_AUTORCC ON)

set(CMAKE_CXX_STANDARD 17)
set(CMAKE_CXX_STANDARD_REQUIRED ON)

# QtCreator supports the following variables for Android, which are identical to qmake Android variables.
# Check http://doc.qt.io/qt-5/deployment-android.html for more information.
# They need to be set before the find_package(Qt5 ...) call.

#if(ANDROID)
#    set(ANDROID_PACKAGE_SOURCE_DIR "${CMAKE_CURRENT_SOURCE_DIR}/android")
#    if (ANDROID_ABI STREQUAL "armeabi-v7a")
#        set(ANDROID_EXTRA_LIBS
#            ${CMAKE_CURRENT_SOURCE_DIR}/path/to/libcrypto.so
#            ${CMAKE_CURRENT_SOURCE_DIR}/path/to/libssl.so)
#    endif()
#endif()

set(QT_COMPONENTS Core Quick QuickControls2 LinguistTools DBus)
if(TIMETERMOS)
  list(APPEND QT_COMPONENTS DBus NetworkSettings)
endif()

set(QT_LINK_TARGETS ${QT_COMPONENTS})
list(REMOVE_ITEM QT_LINK_TARGETS LinguistTools)
list(TRANSFORM QT_LINK_TARGETS PREPEND "Qt5::")

find_package(Qt5 COMPONENTS ${QT_COMPONENTS} REQUIRED)

if(TIMETERMOS)
    qt5_add_dbus_interface(GENERATED_SOURCES dbus/org.freedesktop.systemd1.xml ttsystemd)
endif()

if(CMAKE_BUILD_TYPE STREQUAL "Debug")
  list(APPEND CMAKE_CXX_FLAGS "-fno-omit-frame-pointer")
endif()

if(ENABLE_ASAN)
  list(APPEND CMAKE_CXX_FLAGS "-fsanitize=address")
endif()

if(ENABLE_TSAN)
  list(APPEND CMAKE_CXX_FLAGS "-fsanitize=address -fsanitize=thread")
endif()

set(TS_FILES frontend-embedded_nl_NL.ts)
set(SOURCE_FILES
  qml.qrc
  src/cpp/api/apiclient.cpp
  src/cpp/api/apiclient.h
  src/cpp/api/fakeapiclient.cpp
  src/cpp/api/fakeapiclient.h
  src/cpp/api/timetermuser.cpp
  src/cpp/api/timetermuser.h
  src/cpp/api/zermeloappointment.cpp
  src/cpp/api/zermeloappointment.h
  src/cpp/api/zermeloappointments.cpp
  src/cpp/api/zermeloappointments.h
  src/cpp/cardreader/cardreader.cpp
  src/cpp/cardreader/cardreader.h
  src/cpp/cardreader/cardreadercontroller.cpp
  src/cpp/cardreader/cardreadercontroller.h
  src/cpp/cardreader/fakecardreader.cpp
  src/cpp/cardreader/fakecardreader.h
  src/cpp/cardreader/mfrc522cardreader.cpp
  src/cpp/cardreader/mfrc522cardreader.h
  src/cpp/devcfg/configmanager.cpp
  src/cpp/devcfg/configmanager.h
  src/cpp/devcfg/connmanserviceconfig.cpp
  src/cpp/devcfg/connmanserviceconfig.h
  src/cpp/devcfg/usbmount.cpp
  src/cpp/devcfg/usbmount.h
  src/cpp/devcfg/usbmount_ttos.cpp
  src/cpp/logs/logmanager.cpp
  src/cpp/logs/logmanager.h
  src/cpp/main.cpp
  src/cpp/messagequeue/enums.h
  src/cpp/messagequeue/jetstreamconsumer.cpp
  src/cpp/messagequeue/jetstreamconsumer.h
  src/cpp/messagequeue/natsconnection.cpp
  src/cpp/messagequeue/natsconnection.h
  src/cpp/messagequeue/natsoptions.cpp
  src/cpp/messagequeue/natsoptions.h
  src/cpp/messagequeue/natsstatusstringer.cpp
  src/cpp/messagequeue/natsstatusstringer.h
  src/cpp/messagequeue/strings.cpp
  src/cpp/messagequeue/strings.h
  src/cpp/networking/networkmanager.cpp
  src/cpp/networking/networkmanager.h
  src/cpp/util/scopeguard.h
  src/cpp/util/teardown.h
  src/cpp/util/unixsignalhandler.cpp
  src/cpp/util/unixsignalhandler.h
  src/cpp/messagequeue/natssubscription.h
  src/cpp/messagequeue/natssubscription.cpp
  src/cpp/messagequeue/messages/decoders.h
  src/cpp/messagequeue/messages/decoders.cpp
  src/cpp/api/createdevice.cpp
  src/cpp/api/createdevice.h
  src/cpp/api/device.cpp
  src/cpp/api/device.h
  src/cpp/api/natscreds.cpp
  src/cpp/api/natscreds.h
  src/cpp/devcfg/deviceconfig.cpp
  src/cpp/devcfg/deviceconfig.h
  src/cpp/systemd/systemd.cpp
  src/cpp/systemd/systemd.h
  src/cpp/api/servicesresponse.cpp
  src/cpp/api/servicesresponse.h)

if(ANDROID)
    add_library(frontend-embedded SHARED ${SOURCE_FILES} ${GENERATED_SOURCES} ${TS_FILES})
else()
    add_executable(frontend-embedded ${SOURCE_FILES} ${GENERATED_SOURCES} ${TS_FILES})
endif()

set(BUILD_TESTING NO CACHE BOOL "Build the NATS tests")
set(NATS_BUILD_EXAMPLES NO CACHE BOOL "Build the NATS examples")
set(NATS_BUILD_LIB_SHARED NO CACHE BOOL "Build the NATS library shared")
set(NATS_BUILD_TYPE ${CMAKE_BUILD_TYPE} CACHE STRING "Build type: Release, Debug, RelWithDebInfo, MinRelSize")
set(NATS_BUILD_STREAMING NO CACHE BOOL "Build the NATS library with NATS streaming support")
set(NATS_BUILD_TLS_USE_OPENSSL_1_1_API YES CACHE BOOL "Build the NATS library with OpenSSL 1.1+ support")
add_subdirectory(3rdparty/nats)
add_subdirectory(../proto/src ${CMAKE_BINARY_DIR}/timeterm_proto)

target_compile_definitions(frontend-embedded
  PRIVATE $<$<OR:$<CONFIG:Debug>,$<CONFIG:RelWithDebInfo>>:QT_QML_DEBUG>)
target_link_libraries(frontend-embedded
  PRIVATE ${QT_LINK_TARGETS} nats_static timeterm_proto)
target_include_directories(frontend-embedded
  PRIVATE ${CMAKE_CURRENT_SOURCE_DIR}/3rdparty/nats/src src/cpp)

# Make sure we only compile with the MFRC522 module if we're on the Raspberry Pi.
if(TIMETERMOS)
    add_subdirectory(../mfrc522 ${CMAKE_BINARY_DIR}/mfrc522)

    target_compile_definitions(frontend-embedded PRIVATE TIMETERMOS)
    target_link_libraries(frontend-embedded PRIVATE mfrc522)
endif(TIMETERMOS)

# Required for nats.c
if(WIN32)
    target_link_libraries(frontend-embedded PRIVATE ws2_32)
endif(WIN32)

qt5_create_translation(QM_FILES ${CMAKE_SOURCE_DIR} ${TS_FILES})
