import QtQuick 2.0
import QtQuick.Controls 2.12
import QtQuick.Layouts 1.12
import QtGraphicalEffects 1.0

Page {
    id: router
    width: stackView.width
    height: stackView.height

    function redirectTimetable(timetable) {
        dayView.setTimetable(timetable)
        weekView.setTimetable(timetable)
    }

    TabBar {
        id: menuBar
        anchors.left: parent.left
        anchors.top: parent.top
        anchors.bottom: parent.bottom
        leftPadding: width * 0.125
        width: parent.width * 0.15
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
                radius: parent.width * 0.10

                visible: false
            }

            DropShadow {
                id: rightBottomOuterShadow
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
                id: leftUpperOuterShadow
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
                id: rightBottomInnerShadow
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
                id: leftUpperInnerShadow
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

    Button {
        id: logoutButton
        width: menuBar.width * 0.75
        height: width
        z: 2
        anchors.left: menuBar.left
        anchors.bottom: menuBar.bottom
        anchors.leftMargin: menuBar.width * 0.125
        anchors.bottomMargin: menuBar.width * 0.125

        background: Rectangle {
            id: thirdRect
            color: "#E5E5E5"
            radius: parent.width * 0.10

            visible: false
        }

        DropShadow {
            anchors.fill: thirdRect
            transparentBorder: true
            horizontalOffset: parent.width * 0.10
            verticalOffset: parent.width * 0.10
            radius: parent.width * 0.20
            samples: 32
            color: "#BEBEC0"
            source: thirdRect
        }

        DropShadow {
            anchors.fill: thirdRect
            transparentBorder: true
            horizontalOffset: -parent.width * 0.10
            verticalOffset: -parent.width * 0.10
            radius: parent.width * 0.20
            samples: 32
            color: "#FFFFFF"
            source: thirdRect
        }

        icon.color: "#424242"
        icon.width: width * 0.40
        icon.height: width * 0.40
        icon.source: "../../assets/icons/logout.svg"
        display: AbstractButton.TextUnderIcon

        font.pixelSize: height * 0.16

        text: "<font color=\"#424242\">Log uit</font>"

        onClicked: stackView.pop()
    }

    StackLayout {
        id: stack
        currentIndex: menuBar.currentIndex

        anchors.left: menuBar.right
        anchors.right: parent.right
        height: parent.height

        DayView {
            id: dayView
            startOfDay: new Date().setHours(0, 0, 0, 0)
            endOfDay: new Date().setHours(24, 0, 0, 0)
        }

        WeekView {
            id: weekView
            startOfWeek: new Date().startOfWeek()
            endOfWeek: new Date().endOfWeek()
        }
    }

    ChoosableAppointmentView {
        id: popup
    }
}
