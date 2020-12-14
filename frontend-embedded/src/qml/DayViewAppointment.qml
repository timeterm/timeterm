import QtQuick 2.12

Rectangle {
    property var appointment
    property var startFirstAppointment
    property var secondToPixelRatio

    function setBackgroundColor() { // pick colors with 33% Saturation in HSV
        if (appointment.isCanceled) {
            return "#ffabab"
        } else if (appointment.isOptional) {
            if (!appointment.isStudentEnrolled) {
                return "#ffddab"
            }
            return "#c4ffab"
        }
        return "#e5e5e5"
    }

    function setBorderColor() { // pick corresponding backgroundColor, set Saturation in HSV to 80%
        if (appointment.isCanceled) {
            return "#ff3333"
        } else if (appointment.isOptional) {
            if (!appointment.isStudentEnrolled) {
                return "#ffad33"
            }
            return "#70ff33"
        }
        return "#b5b5b5"
    }

    function setTimeSlot() {
        if (appointment.startTimeSlotName === "0") {
            return ""
        } else if (appointment.startTimeSlotName !== appointment.endTimeSlotName) {
            return appointment.startTimeSlotName + "-" + appointment.endTimeSlotName
        }
        return appointment.startTimeSlotName
    }

    y: (appointment.startTime.getMillisecondsInDay()
       - startFirstAppointment) / 1000 // time in seconds calculated from the start of the first appointment of the day
       * secondToPixelRatio
    height: (appointment.endTime.getMillisecondsInDay() - appointment.startTime.getMillisecondsInDay())
            / 1000 * secondToPixelRatio - 5 // - 5 because of the spacing between appointments

    width: dayHeader.width
    anchors.right: parent.right

    color: setBackgroundColor()
    border.width: 1
    border.color: setBorderColor()
    radius: 5

    Text {
        id: timeSlot
        anchors.left: parent.left
        anchors.leftMargin: dayPage.customMargin
        anchors.verticalCenter: parent.verticalCenter
        font.pixelSize: dayPage.textSize
        text: setTimeSlot()
    }

    Text {
        anchors.left: parent.left
        anchors.leftMargin: parent.width * 0.12
        anchors.right: locations.left
        anchors.rightMargin: dayPage.customMargin
        anchors.verticalCenter: parent.verticalCenter
        font.pixelSize: dayPage.textSize
        elide: Text.ElideRight
        text: appointment.subjects.join(", ")

        property bool isPressed: false

        TapHandler {
            onPressedChanged: parent.isPressed = pressed
        }

        TapToolTip {
            visible: parent.truncated && parent.isPressed
        }
    }

    Text {
        id: locations
        width: (parent.width - dayPage.customMargin * 2) * 0.32
        anchors.horizontalCenter: parent.horizontalCenter
        anchors.verticalCenter: parent.verticalCenter
        horizontalAlignment: Text.AlignHCenter
        font.pixelSize: dayPage.textSize
        elide: Text.ElideRight
        text: appointment.locations.join(", ")

        property bool isPressed: false

        TapHandler {
            onPressedChanged: parent.isPressed = pressed
        }

        TapToolTip {
            visible: parent.truncated && parent.isPressed
        }
    }
    
    Text {
        anchors.left: locations.right
        anchors.leftMargin: dayPage.customMargin
        anchors.right: parent.right
        anchors.rightMargin: dayPage.customMargin
        anchors.verticalCenter: parent.verticalCenter
        horizontalAlignment: Text.AlignRight
        font.pixelSize: dayPage.textSize
        color: "#666666"
        elide: Text.ElideRight
        text: appointment.teachers.join(", ")

        property bool isPressed: false

        TapHandler {
            onPressedChanged: parent.isPressed = pressed
        }

        TapToolTip {
            visible: parent.truncated && parent.isPressed
        }
    }

    TapHandler {
        enabled: appointment.isOptional
        onTapped: function() {
            popup.appointment = appointment
            popup.open()
        }
    }
}
