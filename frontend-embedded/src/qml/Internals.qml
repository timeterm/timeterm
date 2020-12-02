import QtQuick 2.12
import QtQml 2.12
import Timeterm.Api 1.0
import Timeterm.Config 1.0
import Timeterm.MessageQueue 1.0
import Timeterm.Rfid 1.0
import Timeterm.Networking 1.0

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

    FakeApiClient {
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

    ConfigLoader {
        id: configLoader

        Component.onCompleted: {
            configLoader.loadConfig()
        }
    }

    Connections {
        target: configLoader

        function onConfigLoaded() {
            console.log("Config loaded, triggering TtNetworkManager")

            networkManager.configLoaded()
        }
    }

    NetworkManager {
        id: networkManager
    }

    NatsConnection {
        id: natsConn
        options: NatsOptions {
            id: connOpts
            url: "nats.timeterm.nl"
            credsFilePath: "EMDEV.creds"
        }

        Component.onCompleted: {
            natsConn.connect()
        }

        onConnected: {
            console.log("Connected to NATS")

            disownSub.start()
            rebootSub.start()
        }

        onErrorOccurred: function (code, msg) {
            console.log(`An error occurred in the NATS connection: ${msg} (error code ${code})`)
            disownSub.stop()
            rebootSub.stop()

            // Try to reconnect
            natsConnReconnectWait.restart()
        }

        onLastStatusChanged: {
            const status = natsConn.lastStatus
            const statusText = NatsStatusStringer.stringify(status)
            console.log(`NATS connection status changed to ${status} (${statusText})`)
        }

        onConnectionLost: {
            console.log("Connection lost")
            disownSub.stop()
            rebootSub.stop()

            // Try to reconnect
            natsConnReconnectWait.restart()
        }
    }

    NatsSubscription {
        id: rebootSub
        subject: "EMDEV.asdfasdfasdf.REBOOT"
        connection: natsConn
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
