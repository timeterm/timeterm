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
    property var cellMargin: height * 0.015
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

            cellWidth: width / 3
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
                    anchors.margins: cellMargin
                    width: parent.width * 0.5 - anchors.margins * 1.5
                    font.pixelSize: textSize
                    elide: Text.ElideRight
                    text: modelData.teachers.join(", ")

                    property bool isPressed: false

                    TapHandler {
                        onPressedChanged: parent.isPressed = pressed
                    }

                    TapToolTip {
                        visible: parent.truncated && parent.isPressed
                    }
                }

                Text {
                    anchors.top: parent.top
                    anchors.right: parent.right
                    anchors.margins: cellMargin
                    width: parent.width * 0.5 - anchors.margins * 1.5
                    horizontalAlignment: Text.AlignRight
                    font.pixelSize: textSize
                    elide: Text.ElideRight
                    text: modelData.locations.join(", ")

                    property bool isPressed: false

                    TapHandler {
                        onPressedChanged: parent.isPressed = pressed
                    }

                    TapToolTip {
                        visible: parent.truncated && parent.isPressed
                    }
                }

                Text {
                    anchors.bottom: parent.bottom
                    anchors.left: parent.left
                    anchors.right: parent.right
                    anchors.leftMargin: cellMargin
                    anchors.rightMargin: cellMargin
                    anchors.bottomMargin: cellMargin * 2 + textSize*0.75
                    font.pixelSize: textSize * 0.55
                    color: "#666666"
                    wrapMode: Text.Wrap
                    maximumLineCount: 4
                    elide: Text.ElideRight
                    text: modelData.content

                    property bool isPressed: false

                    TapHandler {
                        onPressedChanged: parent.isPressed = pressed
                    }

                    TapToolTip {
                        visible: parent.truncated && parent.isPressed
                    }
                }

                Rectangle {
                    anchors.bottom: parent.bottom
                    anchors.horizontalCenter: parent.horizontalCenter
                    anchors.bottomMargin: cellMargin * 1.5 + textSize*0.75
                    width: parent.width - cellMargin * 2
                    height: 2
                    color: "#666666"
                    border.color: "#666666"
                    border.width: 2
                }

                Text {
                    anchors.bottom: parent.bottom
                    anchors.left: parent.left
                    anchors.right: space.left
                    anchors.leftMargin: cellMargin
                    anchors.bottomMargin: cellMargin
                    //width: parent.width * 0.33 - anchors.margins * 1.5
                    font.pixelSize: textSize * 0.75
                    color: "#666666"
                    elide: Text.ElideRight
                    text: modelData.subjects.join(", ")

                    property bool isPressed: false

                    TapHandler {
                        onPressedChanged: parent.isPressed = pressed
                    }

                    TapToolTip {
                        visible: parent.truncated && parent.isPressed
                    }
                }

                Text {
                    id: space
                    anchors.bottom: parent.bottom
                    anchors.right: parent.right
                    anchors.margins: cellMargin
                    font.pixelSize: textSize * 0.75
                    color: "#666666"
                    text: modelData.availableSpace + " plaatsen vrij"
                }
            }
        }
    }
}
