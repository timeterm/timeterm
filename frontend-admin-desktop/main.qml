import QtQuick 2.15
import QtQuick.Window 2.15
import QtGraphicalEffects 1.0
import QtQuick.Controls 2.15
import QtQuick.Layouts 1.11
import QtQuick.Controls.Material 2.15
import QtQuick.Controls.Material.impl 2.15

Window {
    id: window

    visible: true
    width: 640
    height: 480
    title: qsTr("Timeterm Admin")

    Material.theme: Material.Light
    Material.accent: Material.color(Material.Blue, Material.Shade100)

    MouseArea {
        id: globalMouseArea
        anchors.fill: parent
        acceptedButtons: Qt.NoButton
    }

    Rectangle {
        id: topBar
        anchors.top: parent.top
        anchors.topMargin: 0
        anchors.left: navbar.right
        anchors.leftMargin: -navbar.radius
        anchors.right: parent.right
        anchors.rightMargin: 0
        height: 60
        color: "#182e35"

        Item {
            id: element
            anchors.fill: parent
            anchors.leftMargin: navbar.radius*3

            Text {
                text: "Apparaten"
                anchors.verticalCenter: parent.verticalCenter
                font.pointSize: 24
                color: Qt.rgba(255, 255, 255, 0.78)
                font.family: "Roboto"
                font.weight: Font.Bold
            }
        }
    }

    Rectangle {
        id: navbar
//        x: -radius

        radius: 24
        anchors.left: parent.left
        anchors.leftMargin: -radius
        anchors.top: parent.top
        anchors.topMargin: 0
        anchors.bottom: parent.bottom
        anchors.bottomMargin: 0
        width: 200
        gradient: Gradient {
            GradientStop {
                position: 0
                color: "#0f72cd"
            }

            GradientStop {
                position: 1
                color: "#77b5ee"
            }
        }

        Item {
            anchors.left: parent.left
            anchors.leftMargin: parent.radius
            anchors.top: parent.top
            anchors.right: parent.right
            anchors.bottom: parent.bottom

            Image {
                id: image
                x: 20
                y: 10
                width: 90
                height: 90
                fillMode: Image.PreserveAspectFit
                source: "timeterm-logo-white.svg"
                mipmap: true
            }

            Pane {
                padding: 0
                anchors.top: image.bottom
                anchors.right: parent.right
                anchors.left: parent.left
                background: Rectangle {
                    anchors.fill: parent
                    color: "transparent"
                }

                ColumnLayout {
                    anchors.right: parent.right
                    anchors.left: parent.left

                    Text {
                        text:"test"
                    }

                    Rectangle {
                        property bool highlighted: false

                        Layout.preferredWidth: parent.width
                        Layout.preferredHeight: 60
                        width: parent.width
                        color: highlighted ? Qt.rgba(0, 0, 0, 0.20) : "transparent"

                        Button {
                            anchors.fill: parent

                            Text {
                                text: "Apparaten"
                                anchors.verticalCenter: parent.verticalCenter
                                color: "white"
                            }
                            font.family: "Roboto"
                            flat: true


                        }
                    }

                    Text {
                        text: "fdsa"
                        color: "white"
                        font.family: "Roboto"
                    }
                }
            }

            Rectangle {
                anchors.right: parent.right
                anchors.top: parent.top
                anchors.topMargin: parent.radius
                anchors.bottom: parent.bottom
                anchors.bottomMargin: parent.radius

                MouseArea {
                    anchors.right: parent.right
                    anchors.rightMargin: -1
                    anchors.top: parent.top
                    anchors.bottom: parent.bottom
                    width: 14
                    cursorShape: Qt.SizeHorCursor
                    smooth: true

                    drag {
                        target: parent
                        axis: Drag.XAxis
                    }

                    onPressed: {
                        globalMouseArea.cursorShape = Qt.SizeHorCursor
                    }
                    onReleased:  {
                        globalMouseArea.cursorShape = Qt.ArrowCursor
                    }

                    onMouseXChanged: {
                        if (drag.active) {
                            navbar.width += mouseX
                            if (navbar.width < 150) {
                                navbar.width = 150
                            }
                            if (navbar.width > window.width - 100) {
                                navbar.width = window.width - 100
                            }
                        }
                    }
                }
            }
        }
    }

    DropShadow {
        anchors.fill: navbar
        source: navbar
        verticalOffset: 4
        radius: 50
        smooth: true
        samples: 50
        color: Qt.rgba(0, 0, 0, 0.5)
    }
}



/*##^##
Designer {
    D{i:0;formeditorZoom:0.8999999761581421}
}
##^##*/
