import QtQuick 2.0
import QtQuick.Layouts 1.3
import QtQuick.Controls 2.14
import QtGraphicalEffects 1.0

Item {
    id: header
    width: parent.width
    height: parent.height * 0.07

    Rectangle {
        anchors.fill: parent
        color: "#242424"

        RowLayout {
            anchors.fill: parent
            spacing: 6

            Label {
                color: "#e5e5e5"
                text: "Timeterm"
                fontSizeMode: Text.Fit
                font.pixelSize: 20
                horizontalAlignment: Text.AlignHCenter
            }

            Label {
                id: dateTime
                color: "#e5e5e5"
                anchors.centerIn: parent
                fontSizeMode: Text.Fit
                font.pixelSize: 20
                function setDateTime() {
                    dateTime.text = new Date().toLocaleString(
                                Qt.locale("nl_NL"),
                                "d MMMM yyyy    h:mm:ss") // eg. Donderdag 17 september 2020 13:08
                }
            }

            Label {
                color: "#e5e5e5"
                text: "Wifi"
                fontSizeMode: Text.Fit
                font.pixelSize: 20
                transformOrigin: Item.Center
            }
        }

        layer.enabled: true
        layer.effect: DropShadow {
            transparentBorder: true
            verticalOffset: 8
        }
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
