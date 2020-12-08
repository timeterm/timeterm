import QtQuick 2.12

Page {
    property var appointment
    property var textSize: height * 0.05

    GridView {
        width: parent.width
        height: parent.height * 0.8
        anchors.bottom: parent.bottom

        cellWidth: width / 4.5
        cellHeight: height / 3.5

        //model: appointment.

        delegate: choosableAppointment
    }

    Component {
        id: choosableAppointment
        Rectangle {
            width: cellWidth
            height: cellHeight

            Text {
                anchors.top: parent.top
                anchors.left: parent.left
                font.pixelSize: textSize
                text: teachers.join(", ")
            }

            Text {
                anchors.top: parent.top
                anchors.right: parent.right
                font.pixelSize: textSize
                text: locations.join(", ")
            }

            Text {
                anchors.bottom: parent.bottom
                anchors.left: parent.left
                color: "#666666"
                font.pixelSize: textSize * 0.5
                text: subjects.join(", ")
            }
            
            Text {
                anchors.bottom: parent.bottom
                anchors.right: parent.right
                font.pixelSize: textSize * 0.5
                color: "#666666"
                text: availableSpace
            }
        }
    }
}