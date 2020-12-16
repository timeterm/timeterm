import QtQuick 2.12
import QtQuick.Controls 2.12
import QtQml.Models 2.12
import QtQuick.Layouts 1.12
import QtGraphicalEffects 1.0
import Timeterm.Logging 1.0

import "../js/TimeFunctions.js" as TimeFunctions

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
            text: LogManager.messages.join("\n")
        }
    }

    Connections {
        target: LogManager

        function onMessagesChanged() {
            blabla.text = LogManager.messages.join("\n")
        }
    }
}
