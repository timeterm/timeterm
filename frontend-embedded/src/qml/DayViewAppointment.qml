import QtQuick 2.12

Rectangle {
    property QString startTimeSlot
    property QString endTimeSlot
    property QDateTime startTime
    property QDateTime endTime
    property QStringList subjects
    property QStringList groups
    property QStringList locations
    property QStringList teachers
    property int textSize: 10
    property int customMargin: 20

    width: parent.width
    height: parent.height * (endTime - startTime) / 1000 / 22500
    color: "#e5e5e5"
    radius: 5

    Text {
        anchors.left: parent.left
        anchors.leftMargin: customMargin
        anchors.verticalCenter: parent.verticalCenter
        font.pixelSize: textSize
        text: subjects.join(", ")
    }
    Text {
        anchors.horizontalCenter: parent.horizontalCenter
        anchors.horizontalCenterOffset: -parent.width * 0.125
        anchors.verticalCenter: parent.verticalCenter
        font.pixelSize: textSize
        color: "#666666"
        text: teachers.join(", ")
    }
    Text {
        anchors.horizontalCenter: parent.horizontalCenter
        anchors.horizontalCenterOffset: parent.width * 0.125
        anchors.verticalCenter: parent.verticalCenter
        font.pixelSize: textSize
        color: "#666666"
        text: groups.join(", ")
    }
    Text {
        anchors.right: parent.right
        anchors.rightMargin: customMargin
        anchors.verticalCenter: parent.verticalCenter
        font.pixelSize: textSize
        text: locations.join(", ")
    }
}
