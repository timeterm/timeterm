import QtQuick 2.15
import QtQuick.Controls 2.15

Page {
    id: choosableAppointmentView
    width: stackView.width
    height: stackView.height

    background: Rectangle {
        color: "#AA000000"
    }

    property var appointment
    property var customMargin: height * 0.1
    property var textSize: height * 0.04

    MouseArea {
        anchors.fill: parent
        onClicked: stackView.pop()
    }

    Rectangle {
        anchors.fill: parent
        anchors.topMargin: customMargin
        anchors.leftMargin: stackView.width * 0.15 + customMargin
        anchors.rightMargin: customMargin
        anchors.bottomMargin: customMargin
        radius: 5

        GridView {
            id: grid
            anchors.fill: parent
            clip: true

            cellWidth: width / 4
            cellHeight: height / 3

            model: appointment.alternatives
            delegate: choosableAppointment
        }
    }

    Component {
        id: choosableAppointment
        Rectangle {
            width: grid.cellWidth
            height: grid.cellHeight
            color: "#00000000"

            Rectangle {
                anchors.fill: parent
                anchors.margins: choosableAppointmentView.height * 0.02
                radius: 5
                color: "#e5e5e5"

                Text {
                    anchors.top: parent.top
                    anchors.left: parent.left
                    anchors.margins: parent.width * 0.06
                    font.pixelSize: textSize
                    text: modelData.teachers.join(", ")
                }

                Text {
                    anchors.top: parent.top
                    anchors.right: parent.right
                    anchors.margins: parent.width * 0.06
                    font.pixelSize: textSize
                    text: modelData.locations.join(", ")
                }

                Text {
                    anchors.bottom: parent.bottom
                    anchors.left: parent.left
                    anchors.right: parent.right
                    anchors.leftMargin: parent.width * 0.06
                    anchors.rightMargin: parent.width * 0.06
                    anchors.bottomMargin: parent.width * 0.12 + textSize*0.75
                    font.pixelSize: textSize * 0.75
                    color: "#666666"
                    wrapMode: Text.Wrap
                    text: modelData.content
                }

                Rectangle {
                    anchors.bottom: parent.bottom
                    anchors.horizontalCenter: parent.horizontalCenter
                    anchors.bottomMargin: parent.width * 0.09 + textSize*0.75
                    width: parent.width * 0.88
                    height: 2
                    color: "#666666"
                    border.color: "#666666"
                    border.width: 2
                }

                Text {
                    anchors.bottom: parent.bottom
                    anchors.left: parent.left
                    anchors.margins: parent.width * 0.06
                    font.pixelSize: textSize * 0.75
                    color: "#666666"
                    text: modelData.subjects.join(", ")
                }

                Text {
                    anchors.bottom: parent.bottom
                    anchors.right: parent.right
                    anchors.margins: parent.width * 0.06
                    font.pixelSize: textSize * 0.75
                    color: "#666666"
                    text: modelData.availableSpace + " plaatsen vrij"
                }
            }
        }
    }
}
