import QtQuick 2.15
import QtQuick.Controls 2.15

Popup {
    anchors.centerIn: parent

    padding: 10
    modal: true
    focus: true

    background: Rectangle {
        border.color: "#ffabab"
        radius: 5
    }

    contentItem: Text {
        anchors.centerIn: parent
        font.pixelSize: mainWindow.height * 0.05
        text: "Er is een fout opgetreden"
    }

    // Use this to check if there was some action
    MouseArea {
        anchors.fill: parent
        propagateComposedEvents: true
        onPressed: {
            logoutTimer.restart()
            mouse.accepted = false
        }
    }

    onOpened: errorTimer.restart()

    Timer {
        id: errorTimer
        interval: 2000
        running: false
        onTriggered: errorPopup.close()
    }
}
