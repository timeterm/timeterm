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

    header: HeaderComponent {}

    TabBar {
        id: menuBar
        anchors.left: parent.left
        anchors.top: parent.top
        anchors.bottom: parent.bottom
        width: parent.width * 0.125

        background: Rectangle {
            color: "#e5e5e5"
        }

        TabButton {
            id: dayViewButton
            width: menuBar.width * 0.75
            height: width
            anchors.top: parent.top
            anchors.horizontalCenter: parent.horizontalCenter
            anchors.topMargin: menuBar.width * 0.125
            anchors.bottom: undefined
            text: qsTr("Dagweergave")
        }
        TabButton {
            id: weekViewButton
            width: menuBar.width * 0.75
            height: width
            anchors.top: dayViewButton.bottom
            anchors.horizontalCenter: parent.horizontalCenter
            anchors.topMargin: menuBar.width * 0.125
            anchors.bottom: undefined
            text: qsTr("Weekweergave")
        }
    }

    StackLayout {
        currentIndex: menuBar.currentIndex

        anchors.left: menuBar.right
        anchors.right: parent.right
        height: parent.height

        Item {
            Layout.fillWidth: true
            Layout.fillHeight: true
            DayView {}
        }
        Item {
            Layout.fillWidth: true
            Layout.fillHeight: true
            WeekView {}
        }
    }

    //    Button {
    //        y: 200
    //        text: "blabla"
    //    }
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
