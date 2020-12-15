import QtQuick 2.15
import QtQuick.Controls 2.15

Popup {
    id: choosableAppointmentView
    x: stackView.width * 0.15 + stackView.height * 0.1
    y: stackView.height * 0.1
    width: stackView.width * 0.85 - stackView.height * 0.2
    height: stackView.height - stackView.height * 0.2

    padding: 0
    modal: true
    focus: true

    background: Rectangle {
        border.color: "#399cf8"
        radius: 5
    }

    property var appointment
    property var cellMargin: height * 0.015
    property var textSize: height * 0.05
    property var headerTextColor: "#e5e5e5"
    property var enrollIntoParticipationAllowedActions
    property var enrollIntoParticipationId
    property var unenrollFromParticipationId

    onAppointmentChanged: function() {
        appointmentName.text = appointment.startTime.toLocaleString(Qt.locale("nl_NL"), "dddd") + " " + appointment.startTimeSlotName + "e"
        appointmentTime.text = appointment.startTime.toLocaleString(Qt.locale("nl_NL"), "h:mm") + "-" + appointment.endTime.toLocaleString(Qt.locale("nl_NL"), "h:mm")

        if (appointment.participationId === 0) {
            grid.header = null
        } else {
            grid.header = choosableAppointment
        }

        grid.model = appointment.alternatives

        enrollIntoParticipationAllowedActions = null
        enrollIntoParticipationId = null
        unenrollFromParticipationId = null

        if (appointment.isStudentEnrolled) {
            subscribe.visible = false
            change.visible = true
            unsubscribe.visible = true
            unenrollFromParticipationId = appointment.participationId
        } else {
            subscribe.visible = true
            change.visible = false
            unsubscribe.visible = false
        }
    }

    Rectangle {
        id: header
        height: parent.height * 0.1
        width: parent.width
        color: "#242424"

        Label {
            anchors.horizontalCenter: parent.horizontalCenter
            anchors.horizontalCenterOffset: parent.width / -3
            anchors.top: parent.top
            anchors.bottom: parent.bottom
            verticalAlignment: "AlignVCenter"
            color: headerTextColor
            text: "Keuze-uren"
            fontSizeMode: Text.Fit
            font.pixelSize: textSize
        }

        Label {
            id: appointmentName
            anchors.horizontalCenter: parent.horizontalCenter
            anchors.top: parent.top
            anchors.bottom: parent.bottom
            verticalAlignment: "AlignVCenter"
            color: headerTextColor
            fontSizeMode: Text.Fit
            font.pixelSize: textSize
            font.capitalization: Font.Capitalize
        }

        Label {
            id: appointmentTime
            anchors.horizontalCenter: parent.horizontalCenter
            anchors.horizontalCenterOffset: parent.width / 3
            anchors.top: parent.top
            anchors.bottom: parent.bottom
            verticalAlignment: "AlignVCenter"
            color: headerTextColor
            fontSizeMode: Text.Fit
            font.pixelSize: textSize
        }
    }

    GridView {
        id: grid
        anchors.margins: choosableAppointmentView.height * 0.02
        anchors.top: header.bottom
        anchors.left: parent.left
        anchors.right: parent.right
        anchors.bottom: footer.top
        clip: true

        cellWidth: width / 3
        cellHeight: height / 2.5

        delegate: choosableAppointment
        focus: true
    }

    Component {
        id: choosableAppointment
        Rectangle {
            id: choosableAppointmentRect
            width: grid.cellWidth
            height: grid.cellHeight
            color: "#00000000"

            property bool isHeader: typeof index === 'undefined'

            Rectangle {
                id: appointmentData
                anchors.fill: parent
                anchors.margins: choosableAppointmentView.height * 0.02
                radius: 5

                Component.onCompleted: setBackgroundColor()

                Connections {
                    target: choosableAppointmentView
                    function onEnrollIntoParticipationIdChanged() {
                        appointmentData.setBackgroundColor()
                    }
                }

                function setBackgroundColor() {
                    if (isHeader) {
                        if (appointment.isStudentEnrolled) {
                            color = "#c4ffab"
                        } else if (appointment.availableSpace <= 0) {
                            color = "#ffabab"
                        }
                    } else if (enrollIntoParticipationId === modelData.participationId) {
                        color = "#d6e6ff"
                    } else if (modelData.availableSpace <= 0) {
                        color = "#ffabab"
                    } else {
                        color = "#e5e5e5"
                    }
                }

                function selectThisAppointment() {
                    if (!isHeader) {
                        enrollIntoParticipationId = modelData.participationId
                        enrollIntoParticipationAllowedActions = modelData.allowedStudentActions
                        //color = "#d6e6ff"
                    }
                }

                Text {
                    anchors.top: parent.top
                    anchors.left: parent.left
                    anchors.margins: cellMargin
                    width: parent.width * 0.5 - anchors.margins * 1.5
                    font.pixelSize: textSize
                    elide: Text.ElideRight
                    text: isHeader ? appointment.teachers.join(", ") : modelData.teachers.join(", ")

                    property bool isPressed: false

                    TapHandler {
                        enabled: parent.truncated
                        onPressedChanged: parent.isPressed = pressed
                    }

                    TapToolTip {
                        visible: parent.truncated && parent.isPressed
                    }
                }

                Text {
                    anchors.top: parent.top
                    anchors.right: parent.right
                    anchors.margins: cellMargin
                    width: parent.width * 0.5 - anchors.margins * 1.5
                    horizontalAlignment: Text.AlignRight
                    font.pixelSize: textSize
                    elide: Text.ElideRight
                    text: isHeader ? appointment.locations.join(", ") : modelData.locations.join(", ")

                    property bool isPressed: false

                    TapHandler {
                        enabled: parent.truncated
                        onPressedChanged: parent.isPressed = pressed
                    }

                    TapToolTip {
                        visible: parent.truncated && parent.isPressed
                    }
                }

                Text {
                    anchors.bottom: parent.bottom
                    anchors.left: parent.left
                    anchors.right: parent.right
                    anchors.leftMargin: cellMargin
                    anchors.rightMargin: cellMargin
                    anchors.bottomMargin: cellMargin * 2 + textSize*0.75
                    font.pixelSize: textSize * 0.55
                    color: "#666666"
                    wrapMode: Text.Wrap
                    maximumLineCount: 3
                    elide: Text.ElideRight
                    text: isHeader ? appointment.content : modelData.content

                    property bool isPressed: false

                    TapHandler {
                        enabled: parent.truncated
                        onPressedChanged: parent.isPressed = pressed
                    }

                    TapToolTip {
                        visible: parent.truncated && parent.isPressed
                    }
                }

                Rectangle {
                    anchors.bottom: parent.bottom
                    anchors.horizontalCenter: parent.horizontalCenter
                    anchors.bottomMargin: cellMargin * 1.5 + textSize*0.75
                    width: parent.width - cellMargin * 2
                    height: 2
                    color: "#666666"
                    border.color: "#666666"
                    border.width: 2
                }

                Text {
                    anchors.bottom: parent.bottom
                    anchors.left: parent.left
                    anchors.right: space.left
                    anchors.leftMargin: cellMargin
                    anchors.rightMargin: cellMargin
                    anchors.bottomMargin: cellMargin
                    font.pixelSize: textSize * 0.75
                    color: "#666666"
                    elide: Text.ElideRight
                    text: isHeader ? appointment.subjects.join(", ") : modelData.subjects.join(", ")

                    property bool isPressed: false

                    TapHandler {
                        enabled: parent.truncated
                        onPressedChanged: parent.isPressed = pressed
                    }

                    TapToolTip {
                        visible: parent.truncated && parent.isPressed
                    }
                }

                Text {
                    id: space
                    anchors.bottom: parent.bottom
                    anchors.right: parent.right
                    anchors.margins: cellMargin
                    font.pixelSize: textSize * 0.75
                    color: "#666666"
                    text: setAvailableSpace() + " plaatsen vrij"

                    function setAvailableSpace() {
                        if (isHeader) {
                            if (appointment.availableSpace <= 0) {
                                return "geen"
                            }
                            return appointment.availableSpace
                        } else {
                            if (modelData.availableSpace <= 0) {
                                return "geen"
                            }
                            return modelData.availableSpace
                        }
                    }
                }

                TapHandler {
                    gesturePolicy: TapHandler.WithinBounds
                    onTapped: parent.selectThisAppointment()
                }
            }
        }
    }

    Rectangle {
        id: footer
        anchors.bottom: parent.bottom
        height: parent.height * 0.1
        width: parent.width
        color: "#242424"

        Button {
            id: subscribe
            anchors.centerIn: parent
            width: grid.cellWidth - choosableAppointmentView.height * 0.04
            height: parent.height * 0.5
            enabled: !!enrollIntoParticipationId
                        && (enrollIntoParticipationAllowedActions === "All" || enrollIntoParticipationAllowedActions === "Switch")
            text: "Inschrijven"

            background: Rectangle {
                radius: 5
            }

            onClicked: internals.updateChoice(null, enrollIntoParticipationId)
        }

        Button {
            id: change
            anchors.horizontalCenter: parent.horizontalCenter
            anchors.horizontalCenterOffset: parent.width / -3
            anchors.verticalCenter: parent.verticalCenter
            width: grid.cellWidth - choosableAppointmentView.height * 0.04
            height: parent.height * 0.5
            enabled: !!unenrollFromParticipationId
                        && !!enrollIntoParticipationId
                        && enrollIntoParticipationId !== unenrollFromParticipationId
                        && (appointment.allowedStudentActions === "All" || appointment.allowedStudentActions === "Switch")
                        && (enrollIntoParticipationAllowedActions === "All" || enrollIntoParticipationAllowedActions === "Switch")
            text: "Wijzigen"

            background: Rectangle {
                radius: 5
            }

            onClicked: internals.updateChoice(unenrollFromParticipationId, enrollIntoParticipationId)
        }

        Button {
            anchors.horizontalCenter: parent.horizontalCenter
            anchors.horizontalCenterOffset: parent.width / 3
            anchors.verticalCenter: parent.verticalCenter
            width: grid.cellWidth - choosableAppointmentView.height * 0.04
            height: parent.height * 0.5
            id: unsubscribe
            enabled: !!unenrollFromParticipationId
                        && appointment.allowedStudentActions === "All"
            text: "Uitschrijven"

            background: Rectangle {
                radius: 5
            }

            onClicked: internals.updateChoice(unenrollFromParticipationId, null)
        }
    }

    enter: Transition {
        NumberAnimation { property: "opacity"; from: 0.0; to: 1.0 }
    }

    exit: Transition {
        NumberAnimation { property: "opacity"; from: 1.0; to: 0.0 }
    }

    Connections {
        target: internals
        function onChoiceUpdateSucceeded() {
            close()
        }
    }

    // Use this to check if there was some action
    MouseArea {
        anchors.fill: parent
        propagateComposedEvents: true
        onPressed: {
            logoutTimer.restart()
            mouse.accepted = false
        }
    }
}
