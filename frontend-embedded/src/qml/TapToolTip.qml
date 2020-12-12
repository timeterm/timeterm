import QtQuick 2.15
import QtQuick.Controls 2.15

ToolTip {
    delay: Qt.styleHints.mousePressAndHoldInterval
    text: parent.text

    background: Rectangle {
        border.color: "#399cf8"
        radius: 5
    }

    enter: Transition {
        NumberAnimation { property: "opacity"; from: 0.0; to: 1.0 }
    }

    exit: Transition {
        NumberAnimation { property: "opacity"; from: 1.0; to: 0.0 }
    }
}
