import QtQuick 2.14
import QtQml 2.3
import Timeterm.Rfid 1.0
import Timeterm.Api 1.0
import Timeterm.MessageQueue 1.0

Item {
    id: internalsItem

    signal cardRead(string uid)
    signal timetableReceived(var timetable)

    function getAppointments(start, end) {
        apiClient.getAppointments(start, end)
    }

    Connections {
        target: CardReaderController
        function onCardRead(uid) {
            internalsItem.cardRead(uid)
        }
    }

    ApiClient {
        id: apiClient

        onTimetableReceived: function (timetable) {
            internalsItem.timetableReceived(timetable)
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
}