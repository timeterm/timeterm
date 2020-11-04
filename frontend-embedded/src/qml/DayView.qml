import QtQuick 2.9
import QtQuick.Controls 2.5
import QtQml.Models 2.3
import QtQuick.Layouts 1.3

Page {
    id: dayPage
    anchors.fill: parent
    padding: 32

    property int textSize: dayPage.height * 0.04
    property int customMargin: dayPage.height * 0.05

    background: Rectangle {
        color: "#FFFFFF"
    }

    function setTimetable(timetable) {
        // Pretty-print the timetable as JSON
        console.log()

        console.log(timetable.data[0].startTime.toString())
        console.log((timetable.data[0].endTime - timetable.data[0].startTime) / 1000)

        // Use it
        dayViewList.model = timetable.data
    }

    Component {
        id: dayHeader

        Item {
            width: ListView.view.width
            height: ListView.view.height * 0.08 + dayPage.height * 0.02

            Rectangle {
                anchors.fill: parent
                anchors.bottomMargin: dayPage.height * 0.02
                color: "#b5b5b5"
                radius: 5
                z: 1
                Text {
                    text: "Dinsdag"
                    anchors.verticalCenter: parent.verticalCenter
                    anchors.centerIn: parent
                    font.pixelSize: textSize
                }
            }
        }
    }

    Component {
        id: dayItem

        Rectangle {
            z: -1

            width: ListView.view.width
            //height: ListView.view.height * 0.09
            //Layout.minimumHeight: ListView.view.height * 0.08
            height: ListView.view.height * (modelData.endTime - modelData.startTime) / 1000 / 22500
            color: "#e5e5e5"
            radius: 5

            Text {
                anchors.left: parent.left
                anchors.leftMargin: customMargin
                anchors.verticalCenter: parent.verticalCenter
                font.pixelSize: textSize
                text: modelData.subjects.join(", ")
            }
            Text {
                anchors.horizontalCenter: parent.horizontalCenter
                anchors.horizontalCenterOffset: -parent.width * 0.125
                anchors.verticalCenter: parent.verticalCenter
                font.pixelSize: textSize
                color: "#666666"
                text: modelData.teachers.join(", ")
            }
            Text {
                anchors.horizontalCenter: parent.horizontalCenter
                anchors.horizontalCenterOffset: parent.width * 0.125
                anchors.verticalCenter: parent.verticalCenter
                font.pixelSize: textSize
                color: "#666666"
                text: modelData.groups.join(", ")
            }
            Text {
                anchors.right: parent.right
                anchors.rightMargin: customMargin
                anchors.verticalCenter: parent.verticalCenter
                font.pixelSize: textSize
                text: modelData.locations.join(", ")
            }
        }
    }

    ListView {
        id: dayViewList
        anchors.top: parent.top
        anchors.bottom: parent.bottom
        anchors.right: parent.right
        anchors.margins: parent.height * 0.02

        width: parent.width * 0.8
        header: dayHeader
        headerPositioning: ListView.OverlayHeader
        delegate: dayItem
        spacing: parent.height * 0.02
    }
}
