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

    function openChoosableAppointments() {
        if (appointment.isOptional) {
            stackView.push("ChoosableAppointmentView.qml", {"appointment": appointment})
        }
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
        anchors.left: parent.left
        anchors.leftMargin: weekPage.customMargin
        anchors.verticalCenter: parent.verticalCenter
        font.pixelSize: weekPage.textSize * 0.75
        text: setTimeSlot()
    }

    Text {
        anchors.horizontalCenter: parent.horizontalCenter
        anchors.horizontalCenterOffset: -parent.width / 6 + weekPage.customMargin / 2
        anchors.verticalCenter: parent.verticalCenter
        font.pixelSize: weekPage.textSize * 0.75
        text: appointment.subjects.join(", ")
    }

    Text {
        anchors.horizontalCenter: parent.horizontalCenter
        anchors.horizontalCenterOffset: parent.width / 6 - weekPage.customMargin / 2
        anchors.verticalCenter: parent.verticalCenter
        font.pixelSize: weekPage.textSize * 0.75
        text: appointment.locations.join(", ")
    }
    
    Text {
        anchors.right: parent.right
        anchors.verticalCenter: parent.verticalCenter
        anchors.rightMargin: weekPage.customMargin
        font.pixelSize: weekPage.textSize * 0.75
        color: "#666666"
        text: appointment.teachers.join(", ")
    }

    MouseArea {
        anchors.fill: parent
        onClicked: openChoosableAppointments()
    }
}
