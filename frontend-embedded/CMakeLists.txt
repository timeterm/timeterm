cmake_minimum_required(VERSION 3.5)

project(frontend-embedded LANGUAGES CXX)

set(CMAKE_INCLUDE_CURRENT_DIR ON)

set(CMAKE_AUTOUIC ON)
set(CMAKE_AUTOMOC ON)
set(CMAKE_AUTORCC ON)

set(CMAKE_CXX_STANDARD 11)
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

set(SOURCE_FILES
  main.cpp
  qml.qrc
  mfrc522cardreader.cpp
  mfrc522cardreader.h
  cardreader.h
  cardreader.cpp
  cardreadercontroller.cpp
  cardreadercontroller.h
  fakecardreader.cpp
  fakecardreader.h
  apiclient.cpp
  apiclient.h
)

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

target_compile_definitions(frontend-embedded
  PRIVATE $<$<OR:$<CONFIG:Debug>,$<CONFIG:RelWithDebInfo>>:QT_QML_DEBUG>)
target_link_libraries(frontend-embedded
  PRIVATE Qt5::Core Qt5::Quick Qt5::QuickControls2)

# Make sure we only compile with the MFRC522 module if we're on the Raspberry Pi.
if(RASPBERRYPI)
    add_subdirectory(../mfrc522 ${CMAKE_BINARY_DIR}/mfrc522)

    target_compile_definitions(frontend-embedded PRIVATE RASPBERRYPI)
    target_link_libraries(frontend-embedded PRIVATE mfrc522)
endif()

qt5_create_translation(QM_FILES ${CMAKE_SOURCE_DIR} ${TS_FILES})
