SUMMARY = "Timeterm frontend-embedded"
SECTION = "ui"
LICENSE = "CLOSED"

DEPENDS += "qtdeclarative qtquickcontrols2 qttools qttools-native protobuf protobuf-native libgpiod"

# We're in a monorepo, this should be the root
FILESEXTRAPATHS_prepend := "${THISDIR}/../../../../../:"
# frontend-embedded and proto are in the monorepo root
SRC_URI = "file://frontend-embedded/ \
           file://proto/ \
	   file://mfrc522"

S = "${WORKDIR}/frontend-embedded"

do_install_append () {
	install -d ${D}/opt/frontend-embedded/
	install -m 0755 frontend-embedded ${D}/opt/frontend-embedded/
}

pkg_postinst_ontarget_${PN} () {
#!/bin/sh -e
appcontroller --make-default /opt/frontend-embedded/frontend-embedded
}

FILES_${PN} += "/opt/frontend-embedded/frontend-embedded"

inherit cmake_qt5

EXTRA_OECMAKE += "-DRASPBERRYPI:BOOL=TRUE"
KERNEL_MODULE_AUTOLOAD += "spidev"

