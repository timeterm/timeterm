#!/usr/bin/env sh

# Folders and/or files that need to be ignored when counting:
# 
#  3rdparty/
#  build*/
#  cmake-build-*/
#  CMakeLists.txt.*
#  node_modules/
#  *.log
#
#  /api/example-liveschedule.json
#  /design/js/CCapture.js
#  /design/js/download.js
#  /design/js/p5.js
#  /design/js/webm-writer-0.2.0.js
#  /docs/public/
#  /docs/themes/
#  /frontend-embedded/dbus
#  /mfrc522/docs/_build/
#  /mfrc522/docs/api/
#  /mfrc522/docs/doxyoutput/
#  /mfrc522/docs/ext/
#  /os/sources/
# !/os/sources/meta-timeterm/
#  /os/sstate-cache/
#  /os/setup-environment.sh
#  /proto/build
#  /proto/go

tokei                                   \
  -e '3rdparty/'                        \
  -e 'build*/'                          \
  -e 'cmake-build-*/'                   \
  -e 'CMakeLists.txt.*'                 \
  -e 'node_modules/'                    \
  -e '*.log'                            \
  -e '/api/example-liveschedule.json'   \
  -e '/design/js/CCapture.js'           \
  -e '/design/js/download.js'           \
  -e '/design/js/p5.js'                 \
  -e '/design/js/webm-writer-0.2.0.js'  \
  -e '/docs/public/'                    \
  -e '/docs/themes/'                    \
  -e '/frontend-embedded/dbus'          \
  -e '/mfrc522/docs/_build/'            \
  -e '/mfrc522/docs/api/'               \
  -e '/mfrc522/docs/doxyoutput/'        \
  -e '/mfrc522/docs/ext/'               \
  -e '/os/sources/'                     \
  -e '!/os/sources/meta-timeterm/'      \
  -e '/os/sstate-cache/'                \
  -e '/os/setup-environment.sh'         \
  -e '/proto/build'                     \
  -e '/proto/go'                        \
  $@
