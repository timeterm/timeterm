import QtQuick 2.12
import QtQuick.Layouts 1.12
import QtQuick.Controls 2.12
import QtGraphicalEffects 1.0

Item {
    id: header
    width: parent.parent.width
    height: parent.parent.height * 0.07

    property int textSize: height * 0.5
    property var textColor: "#e5e5e5"
    property var title: ""

    Connections {
        target: header

        function onTitleChanged() {
            titleLabel.text = "Timeterm-" + title
        }
    }

    Rectangle {
        anchors.fill: parent
        color: "#242424"

        Label {
            id: titleLabel
            anchors.left: parent.left
            anchors.leftMargin: parent.height * 0.5
            anchors.top: parent.top
            anchors.bottom: parent.bottom
            verticalAlignment: "AlignVCenter"
            color: textColor
            text: "Timeterm" + title
            fontSizeMode: Text.Fit
            font.pixelSize: textSize
        }

        Label {
            id: dateTime
            anchors.centerIn: parent
            color: textColor
            fontSizeMode: Text.Fit
            font.pixelSize: textSize
            antialiasing: true
            function setDateTime() {
                dateTime.text = new Date().toLocaleString(
                            Qt.locale("nl_NL"),
                            "d MMMM yyyy    h:mm:ss") // eg. 17 september 2020  13:08:22
            }
        }

        Image {
            id: wifi
            anchors.right: parent.right
            anchors.rightMargin: parent.height * 0.5
            anchors.top: parent.top
            anchors.bottom: parent.bottom
            sourceSize.width: parent.height * 0.6
            sourceSize.height: parent.height * 0.6
            fillMode: Image.PreserveAspectFit
            antialiasing: true
            function setWiFiIcon(rssi, connected) {
                if (connected === undefined)
                    connected = true
                if (!connected) {
                    wifi.source = "../../assets/icons/wifi-strength-off-outline.svg"
                } else {
                    if (rssi <= -80) {
                        wifi.source = "../../assets/icons/wifi-strength-outline.svg"
                    } else if (rssi <= -60) {
                        wifi.source = "../../assets/icons/wifi-strength-1.svg"
                    } else if (rssi <= -40) {
                        wifi.source = "../../assets/icons/wifi-strength-2.svg"
                    } else if (rssi <= -20) {
                        wifi.source = "../../assets/icons/wifi-strength-3.svg"
                    } else {
                        wifi.source = "../../assets/icons/wifi-strength-4.svg"
                    }
                }
            }
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
        interval: 1000 // 1 second
        repeat: true
        running: true
        triggeredOnStart: true
        onTriggered: dateTime.setDateTime()
    }

    Component.onCompleted: wifi.setWiFiIcon(
                               -30) // Only for devellopment without true RSSI values
}
