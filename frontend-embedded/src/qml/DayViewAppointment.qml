import QtQuick 2.12

Rectangle {
    property var appointment
    property int textSize: 10
    property int customMargin: 20

    width: parent.width
    height: parent.height * (appointment.endTime - appointment.startTime) / 1000 / 22500
    color: "#e5e5e5"
    radius: 5

    Text {
        anchors.left: parent.left
        anchors.leftMargin: customMargin
        anchors.verticalCenter: parent.verticalCenter
        font.pixelSize: textSize
        text: appointment.subjects.join(", ")
    }
    Text {
        anchors.horizontalCenter: parent.horizontalCenter
        anchors.horizontalCenterOffset: -parent.width * 0.125
        anchors.verticalCenter: parent.verticalCenter
        font.pixelSize: textSize
        color: "#666666"
        text: appointment.teachers.join(", ")
    }
    Text {
        anchors.horizontalCenter: parent.horizontalCenter
        anchors.horizontalCenterOffset: parent.width * 0.125
        anchors.verticalCenter: parent.verticalCenter
        font.pixelSize: textSize
        color: "#666666"
        text: appointment.groups.join(", ")
    }
    Text {
        anchors.right: parent.right
        anchors.rightMargin: customMargin
        anchors.verticalCenter: parent.verticalCenter
        font.pixelSize: textSize
        text: appointment.locations.join(", ")
    }
}
