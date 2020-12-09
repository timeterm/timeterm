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

        model: appointment.alternatives

        delegate: choosableAppointment
    }

    Component {
        id: choosableAppointment
        Rectangle {
            width: cellWidth
            height: cellHeight
            radius: 5
            color: "#e5e5e5"

            Text {
                anchors.top: parent.top
                anchors.left: parent.left
                anchors.margins: parent.width * 0.06
                font.pixelSize: textSize
                text: teachers.join(", ")
            }

            Text {
                anchors.top: parent.top
                anchors.right: parent.right
                anchors.margins: parent.width * 0.06
                font.pixelSize: textSize
                text: locations.join(", ")
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
                text: content
            }

            Rectangle {
                anchors.horizontalCenter: parent.horizontalCenter
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
                text: subjects.join(", ")
            }
            
            Text {
                anchors.bottom: parent.bottom
                anchors.right: parent.right
                anchors.margins: parent.width * 0.06
                font.pixelSize: textSize * 0.75
                color: "#666666"
                text: availableSpace + " plaatsen vrij"
            }
        }
    }
}