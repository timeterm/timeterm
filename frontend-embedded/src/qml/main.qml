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
        initialItem: loginView
        anchors.fill: parent
    }

    Login {
        id: loginView
        visible: false
    }

    Component {
        id: routerComponent

        Router {
            id: routerView
            visible: false
        }
    }

    Internals {
        id: internals

        onCardRead: function (uid) {
            header.title = uid

            stackView.push(routerComponent)

            const startOfWeek = new Date().startOfWeek()
            const endOfWeek = new Date().endOfWeek()
            internals.getAppointments(startOfWeek, endOfWeek)
            stackView.push(routerComponent)
            //stackView.push("Router.qml", {"id": "routerView"})
        }

        onTimetableReceived: function (timetable) {
            console.log("Timetable received")
            console.log(timetable.data[0].locations[0])
            routerView.redirectTimetable(timetable)
        }
    }
}
