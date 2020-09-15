import QtQuick.Controls 2.12
import QtQuick 2.0

ApplicationWindow {
    visible: true

    StackView {
        anchors.fill: parent

        initialItem: Page {
            header: ToolBar {
            }

            Image {
                id: pas
                x: parent.width/2
                y: parent.height/2
                source: "qrc:/assets/images/pas.svg"
            }
        }
    }
}