import QtQuick 2.14
import QtQuick.Window 2.14
import QtQuick.VirtualKeyboard 2.14
import QtQuick.Controls 2.14
import QtQuick.Layouts 1.3
import QtGraphicalEffects 1.0
import QtQml 2.3
import Timeterm.Rfid 1.0
import Timeterm.Api 1.0
import Timeterm.MessageQueue 1.0

ApplicationWindow {
    id: mainWindow
    visible: true
    visibility: Qt.WindowFullScreen
    width: 640
    height: 480
    title: qsTr("Timeterm")

    HeaderComponent {}

    Button {
        y: 200
        text: "blabla"
    }

    Connections {
        target: CardReaderController
        function onCardRead(uid) {
            mainWindow.title = uid
        }
    }

    ApiClient {
        id: apiClient

        onTimetableReceived: function (timetable) {
            console.log("Timetable received")
            console.log(timetable.data[0].locations[0])
        }
    }

    Timer {
        id: natsConnReconnectWait
        repeat: false
        interval: 10000 // wait 10 seconds for reconnection
        onTriggered: {
            console.log("Reconnecting after error")
            natsConn.connect()
        }
    }

    NatsConnection {
        id: natsConn
        options: NatsOptions {
            id: connOpts
            url: "localhost"
        }

        Component.onCompleted: {
            natsConn.connect()
        }

        onConnected: {
            console.log("Connect to NATS")

            disownSub.start()
        }

        onErrorOccurred: function (code, msg) {
            console.log(`Error occurred in NATS connection: ${msg}`)
            disownSub.stop()

            if (code == NatsStatus.NoServer || code == NatsStatus.ConnectionClosed) {
                // Try to reconnect
                natsConnReconnectWait.restart()
            }
        }

        onLastStatusChanged: {
            console.log("NATS connection status changed")
        }

        onConnectionLost: {
            console.log("Connection lost")
            disownSub.stop()

            // Try to reconnect
            natsConnReconnectWait.restart()
        }
    }

    JetStreamConsumer {
        id: disownSub
        connection: natsConn
        stream: "DISOWN-TOKEN"
        consumerId: "ozuhLrexlBa4p50INjihAl"
        type: JetStreamConsumerType.Pull

        onDisownTokenMessage: function (msg) {
            console.log()
        }
    }

    Component.onCompleted: {
        apiClient.getAppointments(new Date(), new Date())
    }

    //    InputPanel {
    //        id: inputPanel
    //        z: 99
    //        x: 0
    //        y: mainWindow.height
    //        width: mainWindow.width

    //        states: State {
    //            name: "visible"
    //            when: inputPanel.active
    //            PropertyChanges {
    //                target: inputPanel
    //                y: mainWindow.height - inputPanel.height
    //            }
    //        }
    //        transitions: Transition {
    //            from: ""
    //            to: "visible"
    //            reversible: true
    //            ParallelAnimation {
    //                NumberAnimation {
    //                    properties: "y"
    //                    duration: 250
    //                    easing.type: Easing.InOutQuad
    //                }
    //            }
    //        }
    //    }
}
