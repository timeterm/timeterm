import QtQuick 2.12
import QtQml 2.12
import Timeterm.Api 1.0
import Timeterm.Config 1.0
import Timeterm.MessageQueue 1.0
import Timeterm.Rfid 1.0
import Timeterm.Networking 1.0
import Timeterm.Systemd 1.0

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
            apiClient.cardId = uid
            internalsItem.cardRead(uid)
        }
    }

    // FakeApiClient {
    //     id: apiClient
    //
    //     onTimetableReceived: function (timetable) {
    //         internalsItem.timetableReceived(timetable)
    //     }
    // }

    ApiClient {
        id: apiClient

        onDeviceCreated: function (response) {
            console.log("Device registered")
            configManager.deviceConfig.id = response.device.id
            configManager.deviceConfig.name = response.device.name
            configManager.deviceConfig.deviceToken = response.token
            configManager.deviceConfig.deviceTokenOrganizationId = response.device.organizationId
            apiClient.apiKey = configManager.deviceConfig.deviceToken

            console.log("Saving device configuration")
            configManager.saveDeviceConfig()

            console.log("Retrieving NATS credentials")
            apiClient.getNatsCreds(configManager.deviceConfig.id)
        }

        onNatsCredsReceived: function (response) {
            console.log("Writing NATS credentials")
            response.writeToFile()

            natsConn.connect()
        }

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

    Timer {
        id: heartbeat
        repeat: true
        interval: 10000 // Wait 10 seconds before reporting a heartbeat again
        onTriggered: {
            apiClient.doHeartbeat(configManager.deviceConfig.id)
        }
    }

    ConfigManager {
        id: configManager

        Component.onCompleted: {
            configManager.loadConfig()
        }

        onConfigLoaded: {
            console.log("Config loaded, triggering TtNetworkManager")

            apiClient.apiKey = configManager.deviceConfig.deviceToken
            networkManager.configLoaded()
        }
    }

    NetworkManager {
        id: networkManager

        onOnlineChanged: function (online) {
            console.log(`Online changed to ${online}`)
            handleNewOnline(online)
        }

        function handleNewOnline(online) {
            if (!online) heartbeat.stop()

            if (online && configManager.deviceConfig.needsRegistration) {
                console.log("Registering")
                apiClient.apiKey = configManager.deviceConfig.setupToken
                apiClient.createDevice()
            } else {
                console.log("Retrieving NATS credentials")
                apiClient.getNatsCreds(configManager.deviceConfig.id)
            }
            heartbeat.start()
        }
    }

    NatsConnection {
        id: natsConn
        options: NatsOptions {
            id: connOpts
            url: "nats.timeterm.nl"
        }

        onConnected: {
            console.log("Connected to NATS")

            rebootSub.useConnection(natsConn)
            retrieveNewNetworkConfigConsumer.useConnection(natsConn)
            rebootSub.start()
            retrieveNewNetworkConfigConsumer.start()
        }

        onErrorOccurred: function (code, msg) {
            console.log(`An error occurred in the NATS connection: ${msg} (error code ${code})`)
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
            console.log("Connection lost, stopping all subscriptions and consumers")
            rebootSub.stop()
            retrieveNewNetworkConfigConsumer.stop()

            // Try to reconnect
            natsConnReconnectWait.restart()
        }
    }

    Systemd {
        id: systemd
    }

    NatsSubscription {
        id: rebootSub
        subject: `EMDEV.${configManager.deviceConfig.id}.REBOOT`

        onMessageReceived: {
            console.log("Rebooting...")
            systemd.rebootDevice()
        }
    }

    JetStreamConsumer {
        id: retrieveNewNetworkConfigConsumer
        stream: "EMDEV-RETRIEVE-NEW-NETWORKING-CONFIG"
        consumerId: configManager.deviceConfig.id

        onMessageReceived: {
            console.log("Retrieving new networking configuration")
        }
    }
}
