import QtQuick 2.12

Rectangle {
    property var appointment
    property var startFirstAppointment
    property var secondToPixelRatio

    y: (appointment.startTime.getMillisecondsInDay()
       - startFirstAppointment.getMillisecondsInDay()) / 1000 // time in seconds calculated from the start of the first appointment of the day
       * secondToPixelRatio
    height: (appointment.endTime.getMillisecondsInDay() - appointment.startTime.getMillisecondsInDay())
            / 1000 * secondToPixelRatio - 5 // - 5 because of the spacing between appointments

    width: dayHeader.width
    anchors.right: parent.right

    color: appointment.isCanceled ? "#FFB5AB" : "#e5e5e5"
    border.width: 1
    border.color: appointment.isCanceled ? "#ff4229" :"#b5b5b5"
    radius: 5

    Text {
        anchors.left: parent.left
        anchors.leftMargin: dayPage.customMargin
        anchors.verticalCenter: parent.verticalCenter
        font.pixelSize: dayPage.textSize
        text: (appointment.startTimeSlot === appointment.endTimeSlot ? appointment.startTimeSlot : appointment.startTimeSlot + " - " + appointment.endTimeSlot)
    }

    Text {
        anchors.horizontalCenter: parent.horizontalCenter
        anchors.horizontalCenterOffset: -parent.width * 0.125
        anchors.verticalCenter: parent.verticalCenter
        font.pixelSize: dayPage.textSize
        text: appointment.subjects.join(", ")
    }

    Text {
        anchors.horizontalCenter: parent.horizontalCenter
        anchors.horizontalCenterOffset: parent.width * 0.125
        anchors.verticalCenter: parent.verticalCenter
        font.pixelSize: dayPage.textSize
        text: appointment.locations.join(", ")
    }
    
    Text {
        anchors.right: parent.right
        anchors.verticalCenter: parent.verticalCenter
        anchors.rightMargin: dayPage.customMargin
        font.pixelSize: dayPage.textSize
        color: "#666666"
        text: appointment.teachers.join(", ")
    }
}
