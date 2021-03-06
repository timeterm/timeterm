import QtQuick 2.14
import QtQuick.Window 2.14
import QtQuick.Controls 2.14
import Timeterm.Rfid 1.0

Window {
    id: window
    visible: true
    width: 640
    height: 480
    title: qsTr("Timeterm devtools")

    FakeCardReaderClient {
        id: fakeCardReaderClient
    }

    Column {
        id: column
        spacing: 20
        padding: 20
        anchors.fill: parent

        TextField {
            id: cardUid
            width: 320
            height: 50
            placeholderText: qsTr("Card UID")
        }

        TextField {
            id: serverName
            width: 320
            height: 50
            placeholderText: qsTr("Server Name")
        }

        Button {
            id: button
            text: qsTr("Send")
            objectName: "button"
            onClicked: fakeCardReaderClient.sendCardUid(serverName.text, cardUid.text)
        }
    }
}

/*##^##
Designer {
    D{i:1;anchors_height:400;anchors_width:200;anchors_x:302;anchors_y:93}
}
##^##*/
