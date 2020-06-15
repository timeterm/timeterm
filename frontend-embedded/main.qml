import QtQuick 2.14
import QtQuick.Window 2.14
import QtQuick.VirtualKeyboard 2.14
import Timeterm.Rfid 1.0
import Timeterm.Api 1.0

Window {
    id: window
    visible: true
    visibility: Qt.WindowFullScreen
    width: 640
    height: 480
    title: qsTr("Timeterm")

    Connections {
        target: CardReader
        function onCardRead(uid) {
            window.title = uid
        }
    }

    Connections {
        target: apiClient
        function onTimetableReceived(timetable) {
            console.log("Timetable received")
            console.log(timetable.data[0].locations[0])
        }
    }

    ApiClient {
        id: apiClient
    }

    Component.onCompleted: {
        apiClient.getAppointments(new Date(), new Date())
    }

    InputPanel {
        id: inputPanel
        z: 99
        x: 0
        y: window.height
        width: window.width

        states: State {
            name: "visible"
            when: inputPanel.active
            PropertyChanges {
                target: inputPanel
                y: window.height - inputPanel.height
            }
        }
        transitions: Transition {
            from: ""
            to: "visible"
            reversible: true
            ParallelAnimation {
                NumberAnimation {
                    properties: "y"
                    duration: 250
                    easing.type: Easing.InOutQuad
                }
            }
        }
    }
}
