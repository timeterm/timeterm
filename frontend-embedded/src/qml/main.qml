import QtQuick 2.14
import QtQuick.Window 2.14
import QtQuick.VirtualKeyboard 2.14
import QtQuick.Controls 2.14
import QtQuick.Layouts 1.3
import QtGraphicalEffects 1.0
import QtQml 2.3

ApplicationWindow {
    id: mainWindow
    visible: true
    visibility: Qt.WindowFullScreen
    width: 640
    height: 480
    title: qsTr("Timeterm")

    HeaderComponent {}

    Button {
        y: 200
        text: "blabla"
    }

    Internals {
        id: internals

        onCardRead: function (uid) {
            mainWindow.title = uid
        }

        onTimetableReceived: function (timetable) {
            console.log("Timetable received")
            console.log(timetable.data[0].locations[0])
        }

        Component.onCompleted: {
            internals.getAppointments(new Date(), new Date())
        }
    }
}
