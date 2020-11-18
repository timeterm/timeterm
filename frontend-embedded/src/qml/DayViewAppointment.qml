import QtQuick 2.12

Rectangle {
    property var appointment
    property var startFirstAppointment
    property var secondToPixelRatio

    y:   (appointment.startTime.getHours() * 3600
        + appointment.startTime.getMinutes() * 60
        + appointment.startTime.getSeconds()
        - (startFirstAppointment.getHours() * 3600
         + startFirstAppointment.getMinutes()* 60
         + startFirstAppointment.getSeconds())) // time in seconds calculated from the first appointment of the day
       * secondToPixelRatio
    height: (appointment.endTime - appointment.startTime)/ 1000 * secondToPixelRatio - 5 // - 5 because of the spacing between appointments

    width: dayHeader.width
    anchors.right: parent.right

    color: "#e5e5e5"
    border.width: 1
    border.color: "#b5b5b5"
    radius: 5

    Text {
        anchors.left: parent.left
        anchors.leftMargin: dayPage.customMargin
        anchors.verticalCenter: parent.verticalCenter
        font.pixelSize: dayPage.textSize
        text: appointment.subjects.join(", ")
    }
    Text {
        anchors.horizontalCenter: parent.horizontalCenter
        anchors.horizontalCenterOffset: -parent.width * 0.125
        anchors.verticalCenter: parent.verticalCenter
        font.pixelSize: dayPage.textSize
        color: "#666666"
        text: appointment.teachers.join(", ")
    }
    Text {
        anchors.horizontalCenter: parent.horizontalCenter
        anchors.horizontalCenterOffset: parent.width * 0.125
        anchors.verticalCenter: parent.verticalCenter
        font.pixelSize: dayPage.textSize
        color: "#666666"
        text: appointment.groups.join(", ")
    }
    Text {
        anchors.right: parent.right
        anchors.rightMargin: dayPage.customMargin
        anchors.verticalCenter: parent.verticalCenter
        font.pixelSize: dayPage.textSize
        text: appointment.locations.join(", ")
    }
}
