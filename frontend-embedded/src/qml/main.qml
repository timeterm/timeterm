import QtQuick 2.12
import QtQuick.Window 2.12
import QtQuick.VirtualKeyboard 2.12
import QtQuick.Controls 2.12
import QtQuick.Layouts 1.12
import QtGraphicalEffects 1.0
import QtQml 2.12

ApplicationWindow {
    id: mainWindow
    visible: true
    visibility: Qt.WindowFullScreen
    width: 640
    height: 480
    title: qsTr("Timeterm")

    header: HeaderComponent {
        id: header
        z: 2
    }

    TabBar {
        id: menuBar
        anchors.left: parent.left
        anchors.top: parent.top
        anchors.bottom: parent.bottom
        //width: 300
        leftPadding: width * 0.125
        width: parent.width * 0.15
        z: 1

        background: Rectangle {
            color: "#e5e5e5"
        }

        TabButton {
            id: dayViewButton
            width: menuBar.width * 0.75
            //width: 150
            height: width
            anchors.left: parent.left
            anchors.top: parent.top
            anchors.topMargin: menuBar.width * 0.125

            background: Rectangle {
                id: firstRect
                color: "#E5E5E5"
                radius: parent.width * 0.10

                visible: false
            }

            DropShadow {
                anchors.fill: firstRect
                transparentBorder: true
                horizontalOffset: parent.width * 0.10
                verticalOffset: parent.width * 0.10
                radius: parent.width * 0.20
                samples: 32
                color: "#BEBEC0"
                source: firstRect
                visible: !dayViewButton.checked
            }

            DropShadow {
                anchors.fill: firstRect
                transparentBorder: true
                horizontalOffset: -parent.width * 0.10
                verticalOffset: -parent.width * 0.10
                radius: parent.width * 0.20
                samples: 32
                color: "#FFFFFF"
                source: firstRect
                visible: !dayViewButton.checked
            }

            InnerShadow {
                anchors.fill: firstRect
                horizontalOffset: -parent.width * 0.033
                verticalOffset: -parent.width * 0.033
                radius: parent.width * 0.20
                samples: 32
                color: "#FFFFFF"
                source: firstRect
                visible: dayViewButton.checked
            }

            InnerShadow {
                anchors.fill: dayViewButton
                horizontalOffset: parent.width * 0.033
                verticalOffset: parent.width * 0.033
                radius: parent.width * 0.20
                samples: 32
                color: "#BEBEC0"
                source: dayViewButton
                visible: dayViewButton.checked
            }

            icon.color: "#424242"
            icon.width: width * (dayViewButton.checked ? 0.375 : 0.40)
            icon.height: width * (dayViewButton.checked ? 0.375 : 0.40)
            icon.source: "../../assets/icons/calendar-today.svg"
            display: AbstractButton.TextUnderIcon

            font.pixelSize: height * (dayViewButton.checked ? 0.15 : 0.16)

            text: "<font color=\"#424242\">Vandaag</font>"
        }

        TabButton {
            id: weekViewButton
            width: menuBar.width * 0.75
            height: width
            anchors.left: parent.left
            anchors.top: dayViewButton.bottom
            anchors.topMargin: menuBar.width * 0.125

            //width: 150
            background: Rectangle {
                id: secondRect
                color: "#E5E5E5"
                radius: parent.width * 0.10

                visible: false
            }

            DropShadow {
                anchors.fill: secondRect
                transparentBorder: true
                horizontalOffset: parent.width * 0.10
                verticalOffset: parent.width * 0.10
                radius: parent.width * 0.20
                samples: 32
                color: "#BEBEC0"
                source: secondRect
                visible: !weekViewButton.checked
            }

            DropShadow {
                anchors.fill: secondRect
                transparentBorder: true
                horizontalOffset: -parent.width * 0.10
                verticalOffset: -parent.width * 0.10
                radius: parent.width * 0.20
                samples: 32
                color: "#FFFFFF"
                source: secondRect
                visible: !weekViewButton.checked
            }

            InnerShadow {
                anchors.fill: secondRect
                horizontalOffset: -parent.width * 0.033
                verticalOffset: -parent.width * 0.033
                radius: parent.width * 0.20
                samples: 32
                color: "#FFFFFF"
                source: secondRect
                visible: weekViewButton.checked
            }

            InnerShadow {
                anchors.fill: weekViewButton
                horizontalOffset: parent.width * 0.033
                verticalOffset: parent.width * 0.033
                radius: parent.width * 0.20
                samples: 32
                color: "#BEBEC0"
                source: weekViewButton
                visible: weekViewButton.checked
            }

            icon.color: "#424242"
            icon.width: width * (weekViewButton.checked ? 0.375 : 0.40)
            icon.height: width * (weekViewButton.checked ? 0.375 : 0.40)
            icon.source: "../../assets/icons/calendar-week.svg"
            display: AbstractButton.TextUnderIcon

            font.pixelSize: height * (weekViewButton.checked ? 0.15 : 0.16)

            text: "<font color=\"#424242\">Week</font>"
        }

        layer.enabled: true
        layer.effect: DropShadow {
            color: "#40000000"
            horizontalOffset: width * 0.02
            verticalOffset: 0
            radius: width * 0.10
            samples: 32
        }
    }

    StackLayout {
        id: stack
        currentIndex: menuBar.currentIndex

        anchors.left: menuBar.right
        anchors.right: parent.right
        height: parent.height

        function redirectTimetable(timetable) {
            //if (currentIndex == 0) {
                dayView.setTimetable(timetable)
            //} else if (currentIndex == 1) {
                weekView.setTimetable(timetable)
            //}
        }

        DayView {
            id: dayView
            startOfDay: new Date().setHours(0, 0, 0, 0)
            endOfDay: new Date().setHours(24, 0, 0, 0)
        }
        WeekView {
            id: weekView
            startOfWeek: new Date().setHours(0, 0, 0, 0)
            endOfWeek: new Date().setHours(24, 0, 0, 0)
        }
    }

    Internals {
        id: internals

        onCardRead: function (uid) {
            header.title = uid
        }

        onTimetableReceived: function (timetable) {
            console.log("Timetable received")
            console.log(timetable.data[0].locations[0])
            stack.redirectTimetable(timetable)
        }

        Component.onCompleted: {
            internals.getAppointments(new Date(), new Date())
        }
    }
}
