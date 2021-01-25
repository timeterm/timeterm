import QtQuick 2.12
import QtQuick.Controls 2.12
import QtQml.Models 2.12
import QtQuick.Layouts 1.12
import QtGraphicalEffects 1.0
import Timeterm.Logging 1.0

Popup {
    anchors.centerIn: parent
    width: parent.width * 0.9
    height: parent.height * 0.9
    padding: 0
    modal: true
    focus: true

    property var textSize: height * 0.05
    property var headerTextColor: "#e5e5e5"

    function networkStateChanged(state) {
        if (state.ip === "") {
            ip.text = ""
        } else {
            ip.text = "Ip: " + state.ip
        }
    }

    function cardUidChanged(uid) {
        if (uid === "") {
            cardUid.text = ""
        } else {
            cardUid.text = "CardUid: " + uid
        }
    }

    background: Rectangle {
        color: "#FFFFFF"
        border.color: "#399cf8"
        radius: 5
    }

    Rectangle {
        id: logsHeader
        width: parent.width
        height: parent.height * 0.1
        color: "#242424"

        Label {
            id: ip
            anchors.horizontalCenter: parent.horizontalCenter
            anchors.horizontalCenterOffset: parent.width / -3
            anchors.top: parent.top
            anchors.bottom: parent.bottom
            verticalAlignment: "AlignVCenter"
            color: headerTextColor
            fontSizeMode: Text.Fit
            font.pixelSize: textSize
        }

        Label {
            anchors.horizontalCenter: parent.horizontalCenter
            anchors.top: parent.top
            anchors.bottom: parent.bottom
            verticalAlignment: "AlignVCenter"
            color: headerTextColor
            text: "Logs"
            fontSizeMode: Text.Fit
            font.pixelSize: textSize
        }

        Label {
            id: cardUid
            anchors.horizontalCenter: parent.horizontalCenter
            anchors.horizontalCenterOffset: parent.width / 3
            anchors.top: parent.top
            anchors.bottom: parent.bottom
            verticalAlignment: "AlignVCenter"
            color: headerTextColor
            fontSizeMode: Text.Fit
            font.pixelSize: textSize
        }
    }

    Flickable {
        id: logsFlickable
        anchors.top: logsHeader.bottom
        anchors.left: parent.left
        anchors.right: parent.right
        anchors.bottom: parent.bottom
        anchors.margins: parent.height * 0.05
        contentHeight: messageList.height

        clip: true

        Text {
            width: parent.width

            id: messageList
            font.family: "Inconsolata"
            wrapMode: Text.Wrap
            font.pixelSize: mainWindow.height * 0.03
            text: LogManager.messages.join("\n")
        }

        // Only show the scrollbars when the view is moving.
        states: State {
            name: "ShowScrollBar"
            when: logsFlickable.movingVertically
            PropertyChanges { target: verticalScrollBar; opacity: 1 }
        }

        transitions: Transition {
            NumberAnimation { properties: "opacity"; duration: 400 }
        }
    }

    Rectangle {
        id: verticalScrollBar
        anchors.top: logsFlickable.top
        anchors.right: logsFlickable.right
        anchors.bottom: logsFlickable.bottom
        width: 10
        opacity: 0

        color: "gray"
        radius: width / 2 - 1

        Rectangle {
            id: scrollbar
            y: logsFlickable.visibleArea.yPosition * logsFlickable.height
            anchors.left: parent.left
            anchors.right: parent.right
            height: logsFlickable.visibleArea.heightRatio * logsFlickable.height
            color: "blue"
            radius: width / 2 - 1
        }
    }

    Connections {
        target: LogManager

        function onMessagesChanged() {
            messageList.text = LogManager.messages.join("\n")
        }
    }

    onOpened: logoutTimer.stop()
    onClosed: logoutTimer.restart()
}
