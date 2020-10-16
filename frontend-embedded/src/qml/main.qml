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

    header: HeaderComponent {
        z: 2
    }

    TabBar {
        id: menuBar
        anchors.left: parent.left
        anchors.top: parent.top
        anchors.bottom: parent.bottom
        leftPadding: width * 0.125
        width: parent.width * 0.125
        z: 1

        background: Rectangle {
            color: "#e5e5e5"
        }

        TabButton {
            id: dayViewButton
            width: menuBar.width * 0.75
            height: width
            anchors.left: parent.left
            anchors.top: parent.top
            anchors.topMargin: menuBar.width * 0.125

            background: Rectangle {
                id: firstRect
                color: "#E5E5E5"
                radius: 10

                visible: false
            }

            DropShadow {
                anchors.fill: firstRect
                transparentBorder: true
                horizontalOffset: 15
                verticalOffset: 15
                radius: 30
                samples: 61
                color: "#BEBEC0"
                source: firstRect
                visible: !dayViewButton.checked
            }

            DropShadow {
                anchors.fill: firstRect
                transparentBorder: true
                horizontalOffset: -15
                verticalOffset: -15
                radius: 30
                samples: 61
                color: "#FFFFFF"
                source: firstRect
                visible: !dayViewButton.checked
            }

            InnerShadow {
                anchors.fill: firstRect
                horizontalOffset: 15
                verticalOffset: 15
                radius: 30
                samples: 61
                color: "#BEBEC0"
                source: firstRect
                visible: dayViewButton.checked
            }

            InnerShadow {
                anchors.fill: firstRect
                horizontalOffset: -15
                verticalOffset: -15
                radius: 30
                samples: 61
                color: "#FFFFFF"
                source: firstRect
                visible: dayViewButton.checked
            }
            text: qsTr("Dagweergave")
        }

        TabButton {
            id: weekViewButton
            width: menuBar.width * 0.75
            height: width
            anchors.left: parent.left
            anchors.top: dayViewButton.bottom
            anchors.topMargin: menuBar.width * 0.125

            background: Rectangle {
                id: secondRect
                color: "#E5E5E5"
                radius: 10

                visible: false
            }

            DropShadow {
                anchors.fill: secondRect
                transparentBorder: true
                horizontalOffset: 15
                verticalOffset: 15
                radius: 30
                samples: 61
                color: "#BEBEC0"
                source: secondRect
                visible: !weekViewButton.checked
            }

            DropShadow {
                anchors.fill: secondRect
                transparentBorder: true
                horizontalOffset: -15
                verticalOffset: -15
                radius: 30
                samples: 61
                color: "#FFFFFF"
                source: secondRect
                visible: !weekViewButton.checked
            }

            InnerShadow {
                anchors.fill: secondRect
                horizontalOffset: 15
                verticalOffset: 15
                radius: 30
                samples: 61
                color: "#BEBEC0"
                source: secondRect
                visible: weekViewButton.checked
            }

            InnerShadow {
                anchors.fill: secondRect
                horizontalOffset: -15
                verticalOffset: -15
                radius: 30
                samples: 61
                color: "#FFFFFF"
                source: secondRect
                visible: weekViewButton.checked
            }

            text: qsTr("Weekweergave")
        }

        layer.enabled: true
        layer.effect: DropShadow {
            color: "#40000000"
            horizontalOffset: 4
            verticalOffset: 0
            radius: 15
            samples: 31
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
