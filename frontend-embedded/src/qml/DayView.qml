import QtQuick 2.9
import QtQuick.Controls 2.5
import QtQml.Models 2.3
import QtQuick.Layouts 1.3

Page {
    id: dayPage
    anchors.fill: parent
    padding: 32

    background: Rectangle {
        color: "#FFFFFF"
    }

    function setTimetable(timetable) {
        // Pretty-print the timetable as JSON
        console.log()

        // Use it
        dayViewList.model = timetable.data
    }

    Component {
        id: dayHeader

        Item {
            width: ListView.view.width
            height: ListView.view.height * 0.06 + dayPage.height * 0.02

            Rectangle {
                anchors.fill: parent
                anchors.bottomMargin: dayPage.height * 0.02
                color: "#52AEAEAE"
                radius: 5
                Text {
                    text: "Dinsdag"
                    anchors.verticalCenter: parent.verticalCenter
                    anchors.centerIn: parent
                    font.pixelSize: parent.height * 0.75
                }
            }
        }
    }

    Component {
        id: dayItem

        Rectangle {
            width: ListView.view.width
            height: ListView.view.height * 0.09
            //Layout.preferredHeight: 50
            color: "#52AEAEAE"
            radius: 5

            Text {
                anchors.left: parent.left
                anchors.leftMargin: parent.height * 0.6
                anchors.verticalCenter: parent.verticalCenter
                font.pixelSize: parent.height * 0.5
                text: modelData.subjects.join(", ")
            }
            Text {
                anchors.horizontalCenter: parent.horizontalCenter
                anchors.horizontalCenterOffset: -parent.width * 0.125
                anchors.verticalCenter: parent.verticalCenter
                font.pixelSize: parent.height * 0.5
                color: "#666666"
                text: modelData.teachers.join(", ")
            }
            Text {
                anchors.horizontalCenter: parent.horizontalCenter
                anchors.horizontalCenterOffset: parent.width * 0.125
                anchors.verticalCenter: parent.verticalCenter
                font.pixelSize: parent.height * 0.5
                color: "#666666"
                text: modelData.groups.join(", ")
            }
            Text {
                anchors.right: parent.right
                anchors.rightMargin: parent.height * 0.6
                anchors.verticalCenter: parent.verticalCenter
                font.pixelSize: parent.height * 0.5
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
        headerPositioning: ListView.PullBackHeader
        delegate: dayItem
        spacing: parent.height * 0.02
    }
}
