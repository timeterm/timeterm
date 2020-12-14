import QtQuick 2.12
import QtQuick.Window 2.12
import QtQuick.Controls 2.12

import "../js/TimeFunctions.js" as TimeFunction

ApplicationWindow {
    id: mainWindow
    visible: true
    visibility: Qt.WindowFullScreen
    width: 1024
    height: 600
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
            stackView.push("Router.qml", {"id": "routerView"})

            header.title = uid

            const startOfWeek = new Date().startOfWeek()
            const endOfWeek = new Date().endOfWeek()
            internals.getAppointments(startOfWeek, endOfWeek)
        }

        onTimetableReceived: function (timetable) {
            console.log("Timetable received")
            console.log(timetable.data[0].locations[0])

            let routerItem = stackView.find(function(item, index) {
                return item instanceof Router})
            if (routerItem) {
                routerItem.redirectTimetable(timetable)
            } else {
                console.log("No routerItem available")
            }
        }

        onNetworkStateChanged: function (state) {
            header.networkStateChanged(state)
        }
    }
}
