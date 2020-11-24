import QtQuick 2.0

Item {
    width: parent.width
    property var time
    property var textSize

    Rectangle {
        anchors.horizontalCenter: parent.horizontalCenter
        width: parent.width * 0.9
        height: 2
        color: "#666666"
        border.color: "#666666"
        border.width: 2
    }
    Text {
        x: parent.width * 0.05
        text: time
        font.pixelSize: textSize
    }
}
