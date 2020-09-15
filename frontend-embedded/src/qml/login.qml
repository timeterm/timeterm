import QtQuick.Controls 2.14
import QtQuick 2.14

Login {
    visible: true

    StackView {
        anchors.fill: parent

        initialItem: Page {
            header: ToolBar {
            }

            Image {
                id: card
                x: parent.width/2
                y: parent.height/2
                source: "qrc:/assets/images/card.svg"
            }
        }
    }
}