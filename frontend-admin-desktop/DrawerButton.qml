import QtQuick 2.15
import QtQuick.Window 2.15
import QtGraphicalEffects 1.0
import QtQuick.Controls 2.15
import QtQuick.Layouts 1.11
import QtQuick.Controls.Material 2.15
import QtQuick.Controls.Material.impl 2.15

Button {
    implicitHeight: 50
    id: drawerButton

    signal depress()
    signal animateDepress()
    signal animatePress()

    onPressed: {
        ripple.pressed = true
        animatePress()
    }

    onDepress: {
        ripple.pressed = false
        animateDepress()
    }

    onAnimatePress: PropertyAnimation {
        target: rect
        easing.type: Easing.InCirc
        property: "opacity"
        to: 1
    }

    onAnimateDepress: PropertyAnimation {
        target: rect
        easing.type: Easing.InCirc
        property: "opacity"
        to: 0
    }

    background: Rectangle {
        anchors.fill: parent
        id: rect
        color: "#235e91"
        opacity: 0

        Ripple {
            clipRadius: 4
            id: ripple
            color: "#235e91"
            anchor: drawerButton
            anchors.fill: parent
        }
    }
}
