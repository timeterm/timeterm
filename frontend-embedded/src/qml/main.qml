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

    header: Rectangle {
        id: header
        width: parent.width
        height: parent.height * 0.07
        color: "#242424"

        RowLayout {
            anchors.fill: parent
            spacing: 6

            Label {
                color: "#e5e5e5"
                text: "Timeterm"
                fontSizeMode: Text.Fit
                font.pixelSize: 20
                horizontalAlignment: Text.AlignHCenter
            }

            Label {
                color: "#e5e5e5"
                text: new Date().toLocaleString(
                          Qt.locale("nl_NL"), "dddd d MMMM yyyy h:mm") // eg. Donderdag 17 september 2020 13:08
                anchors.centerIn: parent
                fontSizeMode: Text.Fit
                font.pixelSize: 20
            }

            Label {
                color: "#e5e5e5"
                text: "Wifi"
                fontSizeMode: Text.Fit
                font.pixelSize: 20
                transformOrigin: Item.Center
            }
        }

        layer.enabled: true
        layer.effect: DropShadow {
            transparentBorder: true
            verticalOffset: 8
        }
    }

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

    NatsConnection {
        id: natsConn
        options: NatsOptions {
            id: connOpts
            url: "localhost"
        }

        Component.onCompleted: {
            console.log()
            console.log()

            natsConn.connect()
        }

        onConnected: {
            console.log("connected")

            disownSub.start()
        }

        onErrorOccurred: function (code, msg) {
            console.log()
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
