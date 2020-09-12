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

    // Timer {
    //     id: natsConnReconnectWait
    //     repeat: false
    //     interval: 10000 // wait 10 seconds for reconnection
    //     onTriggered: {
    //         console.log("Reconnecting after error")
    //         natsConn.connect()
    //     }
    // }

    NatsConnection {
        id: natsConn
        options: NatsOptions {
            id: connOpts
            url: "localhost"
        }

        Component.onCompleted: {
            console.log(`natsConn.lastStatus: ${NatsStatusStringer.stringify(natsConn.lastStatus)}`)
            console.log(`natsConn.connectionOptions.url: ${natsConn.options.url}`)

            natsConn.connect()
        }

        // onConnectionLost: {
        //     console.log("connection lost :(")
        //     console.log("Triggering reconnection lost connection")
        //
        //     // Try to reconnect
        //     natsConnReconnectWait.restart()
        // }

        onConnected: {
            console.log("connected")

            disownSub.subscribe()
        }

        onErrorOccurred: function(code, msg) {
            console.log(`natsConn: Error occurred: code ${code}, message: ${msg}`)
            console.log("Triggering reconnection after error")

            // if (code == NatsStatus.NoServer) {
            //     // Try to reconnect
            //     natsConnReconnectWait.restart()
            // }
        }

        onLastStatusChanged: {
            console.log("status changed")
        }
    }

    NatsSubscription {
        id: disownSub
        target: natsConn
        topic: "timeterm.disown-token"

        onDisownTokenMessage: function(msg) {
            console.log(`device ${msg.deviceId} has to disown their token`)
        }

        onErrorOccurred: function(code, msg) {
            console.log(`error occurred: code ${code}, message: ${msg}`)
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
