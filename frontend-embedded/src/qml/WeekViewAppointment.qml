import QtQuick 2.12

Rectangle {
    property var appointment
    property var startFirstAppointment
    property var secondToPixelRatio
    property var weekAppointmentWidth

    x: weekTimeLine.width + weekPage.height * 0.02 + (appointment.startTime.getDay() - 1) * (weekAppointmentWidth + weekPage.height * 0.02); // timeLineWidth + margin + day * (weekAppointmentWidth + margin)
    y: (appointment.startTime.getMillisecondsInDay()
       - startFirstAppointment) / 1000 // time in seconds calculated from the start of the first appointment of the day
       * secondToPixelRatio
       + weekPage.height * 0.08
    height: (appointment.endTime.getMillisecondsInDay() - appointment.startTime.getMillisecondsInDay())
            / 1000 * secondToPixelRatio - 5 // - 5 because of the spacing between appointments

    width: weekAppointmentWidth

    color: appointment.isCanceled ? "#FFB5AB" : "#e5e5e5"
    border.width: 1
    border.color: appointment.isCanceled ? "#ff4229" :"#b5b5b5"
    radius: 5

    Text {
        anchors.left: parent.left
        anchors.leftMargin: weekPage.customMargin
        anchors.verticalCenter: parent.verticalCenter
        font.pixelSize: weekPage.textSize * 0.75
        text: (appointment.startTimeSlotName === appointment.endTimeSlotName ? appointment.startTimeSlotName : appointment.startTimeSlotName + " - " + appointment.endTimeSlotName)
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
}
