import QtQuick 2.0
import QtQuick.Controls 2.5
import Timeterm.Logging 1.0

Page {
    anchors.fill: parent
    padding: 32

    background: Rectangle {
        color: "#FFFFFF"
    }

    Flickable {
        anchors.fill: parent
        contentHeight: blabla.height

        Text {
            width: parent.width

            id: blabla
            font.family: "Fira Code"
            wrapMode: Text.Wrap
            text: TtLogManager.messages.join("\n")
        }
    }

    Connections {
        target: TtLogManager

        function onMessagesChanged() {
            blabla.text = TtLogManager.messages.join("\n")
        }
    }
}
