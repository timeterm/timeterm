import QtQuick 2.0
import QtQuick.Controls 2.5

Page {
    anchors.fill: parent

    background: Rectangle {
        color: "#FFFFFF"
    }

    function setTimetable(timetable) {
        console.log(timetable.data[0].locations[0])
    }

    Text {
        id: blabla
        text: qsTr("Nog meer prachtige tekst hierzo")
    }
}
