import QtQuick.Controls 2.12
import QtQuick 2.12
import QtGraphicalEffects 1.0

Page {
    id: loginPage
    width: stackView.width
    height: stackView.height

    Image {
        id: card
        x: parent.width * 0.75 - paintedWidth
        width: parent.width * 0.4
        anchors.verticalCenter: parent.verticalCenter
        source: "qrc:/assets/images/card.svg"
        fillMode: Image.PreserveAspectFit
        visible: false
    }

    DropShadow {
        id: cardShadow
        anchors.fill: card
        transparentBorder: true
        horizontalOffset: card.paintedWidth * 0.02
        verticalOffset: card.paintedWidth * 0.02
        radius: card.paintedWidth * 0.08
        samples: 32
        color: "#40000000"
        source: card
    }
}
