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

    Connections {
        target: stanConn

        function onConnectionLost() {
            console.log("connection lost :(")
        }

        function onConnected() {
            console.log("connected")

            let sub = stanConn.subscribe("timeterm.disown-token", stanSubOpts)
            sub.messageReceived.connect(sendMessageToBSPC)
        }

        function sendMessageToBSPC(message) {
            console.log("sendMessageToBSPC")
            bspc.handleMessage(message)
        }

        function onErrorOccurred(code, msg) {
            console.log("error occurred: code " + code + ", message: " + msg)
        }

        function onLastStatusChanged() {
            console.log("status changed")
        }
    }

    StanSubOptions {
        id: stanSubOpts

        Component.onCompleted: {
            stanSubOpts.setDurableName("subz")
        }
    }

    Connections {
        target: stanSubOpts

        function onErrorOccurred() {
            console.log("stanSubOpts: error occurred")
        }
    }

    StanConnection {
        id: stanConn
        cluster: "test-cluster"
        clientId: "example"
        connectionOptions: StanConnectionOptions {
            id: connOpts
            url: "localhost"
        }

        Component.onCompleted: {
            console.log("stanConn.lastStatus: " + NatsStatusStringer.stringify(stanConn.lastStatus))
            console.log("stanConn.connectionOptions.url: " + stanConn.connectionOptions.url)

            stanConn.connect()
        }
    }

    BinaryProtoClient {
        id: bspc

        onDisownToken: function(msg) {
            console.log("device " + msg.deviceId + " has to disown their token")
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
