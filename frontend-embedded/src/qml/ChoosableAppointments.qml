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
            Column {
                anchors.fill: parent
                Text {
                    text: subjects.join(", ")
                }
                Text {
                    text: location.join(", ")
                }
                Text {
                    text: teachers.join(", ")
                }
                Text {
                    text: availableSpace
                }
            }
        }
    }
}