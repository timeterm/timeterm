cmake_minimum_required(VERSION 3.17)

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

find_package(Qt5 COMPONENTS Core Quick QuickControls2 LinguistTools REQUIRED)

set(TS_FILES frontend-embedded_nl_NL.ts)
set(SOURCE_FILES src/cpp/main.cpp qml.qrc src/cpp/cardreader/mfrc522cardreader.cpp src/cpp/cardreader/mfrc522cardreader.h src/cpp/cardreader/cardreader.h src/cpp/cardreader/cardreader.cpp src/cpp/cardreader/cardreadercontroller.cpp src/cpp/cardreader/cardreadercontroller.h src/cpp/cardreader/fakecardreader.cpp src/cpp/cardreader/fakecardreader.h src/cpp/api/apiclient.cpp src/cpp/api/apiclient.h src/cpp/messagequeue/enums.h src/cpp/api/timetermuser.h src/cpp/api/timetermuser.cpp src/cpp/api/zermeloappointments.h src/cpp/api/zermeloappointment.h src/cpp/api/zermeloappointment.cpp src/cpp/api/zermeloappointments.cpp src/cpp/messagequeue/strings.h src/cpp/messagequeue/strings.cpp src/cpp/messagequeue/messages/disowntokenmessage.h src/cpp/messagequeue/messages/disowntokenmessage.cpp src/cpp/messagequeue/messages/retrievenewtokenmessage.cpp src/cpp/messagequeue/messages/retrievenewtokenmessage.h src/cpp/messagequeue/natsoptions.cpp src/cpp/messagequeue/natsoptions.h src/cpp/messagequeue/natsstatusstringer.cpp src/cpp/messagequeue/natsstatusstringer.h src/cpp/util/scopeguard.h src/cpp/util/scopedpointer.h src/cpp/util/unixsignalhandler.cpp src/cpp/util/unixsignalhandler.h src/cpp/util/teardown.h src/cpp/messagequeue/natsconnection.cpp src/cpp/messagequeue/natsconnection.h src/cpp/messagequeue/jetstreamconsumer.cpp src/cpp/messagequeue/jetstreamconsumer.h src/cpp/messagequeue/natscallbackhandlersingleton.cpp src/cpp/messagequeue/natscallbackhandlersingleton.h)

if(ANDROID)
    add_library(frontend-embedded SHARED
      ${SOURCE_FILES}
      ${TS_FILES}
    )
else()
    add_executable(frontend-embedded
      ${SOURCE_FILES}
      ${TS_FILES}
    )
endif()

set(BUILD_TESTING NO CACHE BOOL "Build the NATS tests")
set(NATS_BUILD_EXAMPLES NO CACHE BOOL "Build the NATS examples")
set(NATS_BUILD_LIB_SHARED NO CACHE BOOL "Build the NATS library shared")
set(NATS_BUILD_TYPE ${CMAKE_BUILD_TYPE} CACHE STRING "Build type: Release, Debug, RelWithDebInfo, MinRelSize")
add_subdirectory(3rdparty/nats)
add_subdirectory(../proto/src ${CMAKE_BINARY_DIR}/timeterm_proto)

target_compile_definitions(frontend-embedded
  PRIVATE $<$<OR:$<CONFIG:Debug>,$<CONFIG:RelWithDebInfo>>:QT_QML_DEBUG> NATS_HAS_STREAMING)
target_link_libraries(frontend-embedded
  PRIVATE Qt5::Core Qt5::Quick Qt5::QuickControls2 nats_static timeterm_proto)
target_include_directories(frontend-embedded
  PRIVATE ${CMAKE_CURRENT_SOURCE_DIR}/3rdparty/nats/src src/cpp)

# Make sure we only compile with the MFRC522 module if we're on the Raspberry Pi.
if(RASPBERRYPI)
    add_subdirectory(../mfrc522 ${CMAKE_BINARY_DIR}/mfrc522)

    target_compile_definitions(frontend-embedded PRIVATE RASPBERRYPI)
    target_link_libraries(frontend-embedded PRIVATE mfrc522)
endif()

qt5_create_translation(QM_FILES ${CMAKE_SOURCE_DIR} ${TS_FILES})
