import QtQuick 2.14
import QtQuick.Window 2.14
import QtQuick.VirtualKeyboard 2.14
import QtQuick.Controls 2.14
import Timeterm.Rfid 1.0
import Timeterm.Api 1.0
import Timeterm.MessageQueue 1.0

Window {
    id: window
    visible: true
    visibility: Qt.WindowFullScreen
    width: 640
    height: 480
    title: qsTr("Timeterm")

    Button {
        text: "blabla"
    }

    Connections {
        target: CardReaderController
        function onCardRead(uid) {
            window.title = uid
        }
    }

    ApiClient {
        id: apiClient

        onTimetableReceived: function(timetable) {
            console.log("Timetable received")
            console.log(timetable.data[0].locations[0])
        }
    }

    NatsConnection {
        id: natsConn
        options: NatsOptions {
            id: connOpts
            url: "localhost"
        }

        Component.onCompleted: {
            console.log(`natsConn.lastStatus: ${NatsStatusStringer.stringify(natsConn.lastStatus)}`)
            console.log(`natsConn.options.url: ${natsConn.options.url}`)

            natsConn.connect()
        }

        onConnected: {
            console.log("connected")

            disownSub.start()
        }

        onErrorOccurred: function(code, msg) {
            console.log(`natsConn: Error occurred: code ${code}, message: ${msg}`)
        }

        onLastStatusChanged: {
            console.log("status changed")
        }
    }

    JetStreamConsumer {
        id: disownSub
        target: natsConn
        stream: "DISOWN-TOKEN"
        consumerId: "ozuhLrexlBa4p50INjihAl"
        type: JetStreamConsumerType.Pull

        onDisownTokenMessage: function(msg) {
            console.log(`device ${msg.deviceId} has to disown their token`)
        }
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
