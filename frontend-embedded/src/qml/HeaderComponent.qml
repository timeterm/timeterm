import QtQuick 2.0
import QtQuick.Layouts 1.3
import QtQuick.Controls 2.14
import QtGraphicalEffects 1.0

Item {
    id: header
    width: parent.width
    height: parent.height * 0.07

    property int textSize: height * 0.5
    property var textColor: "#e5e5e5"

    Rectangle {
        anchors.fill: parent
        color: "#242424"

        Label {
            id: title
            anchors.left: parent.left
            anchors.leftMargin: parent.height * 0.5
            anchors.top: parent.top
            anchors.bottom: parent.bottom
            verticalAlignment: "AlignVCenter"
            color: textColor
            text: "Timeterm"
            fontSizeMode: Text.Fit
            font.pixelSize: textSize
        }

        Label {
            id: dateTime
            anchors.centerIn: parent
            color: textColor
            fontSizeMode: Text.Fit
            font.pixelSize: textSize
            function setDateTime() {
                dateTime.text = new Date().toLocaleString(
                            Qt.locale("nl_NL"),
                            "d MMMM yyyy    h:mm:ss") // eg. 17 september 2020  13:08:22
            }
        }

        Label {
            id: wifi
            color: textColor
            anchors.right: parent.right
            anchors.rightMargin: parent.height * 0.5
            anchors.top: parent.top
            anchors.bottom: parent.bottom
            verticalAlignment: "AlignVCenter"
            text: "Wifi"
            fontSizeMode: Text.Fit
            font.pixelSize: textSize
        }
    }

    layer.enabled: true
    layer.effect: DropShadow {
        color: "#40000000"
        horizontalOffset: 0
        verticalOffset: 3
        radius: 8
        samples: 17
        spread: 0
    }

    Timer {
        id: dateTimeTimer
        interval: 1000
        repeat: true
        running: true
        triggeredOnStart: true
        onTriggered: dateTime.setDateTime()
    }
}
