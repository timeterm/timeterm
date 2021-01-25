import QtQuick 2.12
import QtQuick.Window 2.12
import QtQuick.Controls 2.12

import "../js/TimeFunctions.js" as TimeFunction

ApplicationWindow {
    id: mainWindow
    visible: true
    //visibility: Qt.WindowFullScreen
    width: 1280
    height: 800
    title: qsTr("Timeterm")

    header: HeaderComponent {
        id: header
        z: 2
    }

    StackView {
        id: stackView
        initialItem: Login {}
        anchors.fill: parent
    }

    Internals {
        id: internals

        onCardRead: function (uid) {
            if (internals.getApiClientCardUid() === "") {
                internals.setApiClientCardUid(uid)
                logoutTimer.restart()

                const startOfWeek = new Date().startOfWeek()
                const endOfWeek = new Date().endOfWeek()
                internals.getAppointments(startOfWeek, endOfWeek)
            }
        }

        onTimetableReceived: function (timetable) {
            let routerItem = stackView.find(function(item, index) {
                return item instanceof Router})
            if (routerItem) {
                routerItem.redirectTimetable(timetable)
            } else {
                stackView.push("Router.qml", {"id": "routerView"})
                internals.timetableReceived(timetable)
            }
        }

        onTimetableRequestFailed: function () {
            errorPopup.open()
        }

        onChoiceUpdateSucceeded: function () {
            const startOfWeek = new Date().startOfWeek()
            const endOfWeek = new Date().endOfWeek()
            internals.getAppointments(startOfWeek, endOfWeek)
        }

        onChoiceUpdateFailed: function () {
            errorPopup.open()
        }

        onNetworkStateChanged: function (state) {
            header.networkStateChanged(state)
        }
    }

    ErrorPopup {
        id: errorPopup
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

    Timer {
        id: logoutTimer
        interval: 1000 * 10 // 10 seconds
        running: false
        onTriggered: {
            stackView.pop(null) // logout
            internals.setApiClientCardUid("")
        }
        //onRunningChanged: console.log("Timer running? " + running)
    }
}
