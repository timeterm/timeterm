import QtQuick 2.0
import QtQuick.Controls 2.5
import QtQml.Models 2.3

Page {
    id: dayPage
    anchors.fill: parent

    background: Rectangle {
        color: "#FFFFFF"
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
            color: "#52AEAEAE"
            radius: 5
            Text {
                anchors.left: parent.left
                anchors.leftMargin: parent.height * 0.6
                anchors.verticalCenter: parent.verticalCenter
                font.pixelSize: parent.height * 0.5
                text: subjects.join(", ")
            }
            Text {
                text: teacher
                anchors.horizontalCenter: parent.horizontalCenter
                anchors.horizontalCenterOffset: -parent.width * 0.125
                anchors.verticalCenter: parent.verticalCenter
                font.pixelSize: parent.height * 0.5
                color: "#666666"
            }
            Text {
                text: group
                anchors.horizontalCenter: parent.horizontalCenter
                anchors.horizontalCenterOffset: parent.width * 0.125
                anchors.verticalCenter: parent.verticalCenter
                font.pixelSize: parent.height * 0.5
                color: "#666666"
            }
            Text {
                text: location
                anchors.right: parent.right
                anchors.rightMargin: parent.height * 0.6
                anchors.verticalCenter: parent.verticalCenter
                font.pixelSize: parent.height * 0.5
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
        model: TestdataDayView {}
        delegate: dayItem
        spacing: parent.height * 0.02
    }
}
