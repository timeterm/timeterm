import QtQuick 2.12

Rectangle {
    property var appointment
    property var startFirstAppointment
    property var secondToPixelRatio
    property var weekAppointmentWidth

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

    x: weekTimeLine.width + weekPage.height * 0.02 + (appointment.startTime.getDay() - 1) * (weekAppointmentWidth + weekPage.height * 0.02); // timeLineWidth + margin + day * (weekAppointmentWidth + margin)
    y: (appointment.startTime.getMillisecondsInDay()
       - startFirstAppointment) / 1000 // time in seconds calculated from the start of the first appointment of the day
       * secondToPixelRatio
       + weekPage.height * 0.08
    height: (appointment.endTime.getMillisecondsInDay() - appointment.startTime.getMillisecondsInDay())
            / 1000 * secondToPixelRatio - 5 // - 5 because of the spacing between appointments

    width: weekAppointmentWidth

    color: setBackgroundColor()
    border.width: 1
    border.color: setBorderColor()
    radius: 5

    Text {
        id: timeSlot
        anchors.left: parent.left
        anchors.leftMargin: weekPage.customMargin
        anchors.verticalCenter: parent.verticalCenter
        font.pixelSize: weekPage.textSize * 0.7
        text: setTimeSlot()
    }

    Text {
        anchors.left: parent.left
        anchors.leftMargin: parent.width * 0.2
        anchors.right: locations.left
        anchors.rightMargin: weekPage.customMargin * 0.5
        anchors.verticalCenter: parent.verticalCenter
        font.pixelSize: weekPage.textSize * 0.7
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
        width: (parent.width - weekPage.customMargin * 2) * 0.3
        anchors.horizontalCenter: parent.horizontalCenter
        anchors.horizontalCenterOffset: parent.width * 0.05
        anchors.verticalCenter: parent.verticalCenter
        horizontalAlignment: Text.AlignHCenter
        font.pixelSize: weekPage.textSize * 0.7
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
        anchors.leftMargin: weekPage.customMargin * 0.5
        anchors.right: parent.right
        anchors.rightMargin: weekPage.customMargin
        anchors.verticalCenter: parent.verticalCenter
        horizontalAlignment: Text.AlignRight
        font.pixelSize: weekPage.textSize * 0.7
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

    Text {
        visible: appointment.participationId === 0
        anchors.centerIn: parent
        font.pixelSize: weekPage.textSize * 0.7
        text: "Inschrijven"
    }

    TapHandler {
        enabled: appointment.isOptional
        onTapped: function() {
            popup.appointment = appointment
            popup.open()
            console.log("TapHandler tapped")
        }
    }
}
