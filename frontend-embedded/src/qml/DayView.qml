import QtQuick 2.12
import QtQuick.Controls 2.12
import QtQml.Models 2.12
import QtQuick.Layouts 1.12
import QtGraphicalEffects 1.0

import "../js/TimeFunctions.js" as TimeFunction

Page {
    id: dayPage
    anchors.fill: parent

    property int textSize: dayPage.height * 0.04
    property int customMargin: dayPage.height * 0.05
    property var secondToPixelRatio: appointments.height * 0.000037
    property var startOfDay
    property var endOfDay

    background: Rectangle {
        color: "#FFFFFF"
    }

    function setTimetable(timetable) {
        for (var i = 0; i < timetable.data.length; i++) {
            if (typeof startOfDay === 'undefined' || typeof endOfDay === 'undefined') {
                startOfDay = new Date().setHours(0, 0, 0, 0)
                endOfDay = new Date().setHours(24, 0, 0, 0)
            }

            if (timetable.data[i].startTime.getTime() >= startOfDay && timetable.data[i].endTime.getTime() < endOfDay) {
                if (typeof appointments.startFirstAppointment === 'undefined' || timetable.data[i].startTime.getMillisecondsInDay() < appointments.startFirstAppointment.getMillisecondsInDay()) {                                                          // first appointment in the list
                    appointments.startFirstAppointment = timetable.data[i].startTime
                }
                if (typeof appointments.endLastAppointment === 'undefined' || timetable.data[i].endTime.getMillisecondsInDay() > appointments.endLastAppointment.getMillisecondsInDay()) {
                    appointments.endLastAppointment = timetable.data[i].endTime
                    appointments.contentHeight = (appointments.endLastAppointment.getMillisecondsInDay()
                                                - appointments.startFirstAppointment.getMillisecondsInDay())
                                                / 1000 * secondToPixelRatio - 5         // - 5 because of the spacing between appointments

                    fillTimeLine()
                }

                let finishCreation = function (appointment) {
                    if (appointment.status === Component.Ready) {
                        appointment.incubateObject(appointments.contentItem, {
                            appointment: timetable.data[i],
                            startFirstAppointment: appointments.startFirstAppointment,
                            secondToPixelRatio: secondToPixelRatio
                        })
                    } else if (appointment.status === Component.Error) {
                        console.log("Could not create appointment:", appointment.errorString())
                    }
                }

                let appointment = Qt.createComponent("DayViewAppointment.qml")
                if (appointment.status !== Component.Null && appointment.status !== Component.Loading) {
                    finishCreation(appointment)
                } else {
                    appointment.statusChanged.connect(finishCreation)
                }
            }
        }
    }

    function fillTimeLine() {
        let currentLineTime = new Date()
        currentLineTime.setTime(appointments.startFirstAppointment.getTime())

        if (!currentLineTime.isFullHour()) {
            currentLineTime.addHours(1)
            currentLineTime.setMinutes(0, 0, 0)
        }

        for (currentLineTime;
             appointments.endLastAppointment.getMillisecondsInDay() - currentLineTime.getMillisecondsInDay() > 30 * 60 * 1000; // Last timeLine is at least less than 30 minutes before the end of the last appointment
             currentLineTime.addHours(1)) {
            let finishLineItem = function (timeLineItem) {
                if (timeLineItem.status === Component.Ready) {
                    timeLineItem.incubateObject(timeLine, {
                        y: (currentLineTime.getMillisecondsInDay() - appointments.startFirstAppointment.getMillisecondsInDay()) / 1000 * secondToPixelRatio,
                        time: currentLineTime.getHours().toString() + ":"
                              + (currentLineTime.getMinutes().toString() < 10 ? '0' : '')
                              + currentLineTime.getMinutes()
                    })
                } else if (timeLineItem.status === Component.Error) {
                    console.log("Could not create lineItem:",
                                timeLineItem.errorString())
                }
            }

            let timeLineItem = Qt.createComponent("DayViewTimeLineItem.qml")
            if (timeLineItem.status !== Component.Null && timeLineItem.status !== Component.Loading) {
                finishLineItem(timeLineItem)
            } else {
                timeLineItem.statusChanged.connect(finishLineItem)
            }
        }
    }

    Rectangle {
        id: dayHeader
        width: parent.width * 0.8
        height: parent.height * 0.06
        anchors.top: parent.top
        anchors.right: parent.right
        anchors.margins: parent.height * 0.02
        color: "#b5b5b5"
        radius: 5
        Text {
            text: new Date().toLocaleString(Qt.locale("nl_NL"), "dddd")
            anchors.verticalCenter: parent.verticalCenter
            anchors.centerIn: parent
            font.pixelSize: textSize
        }
    }

    Flickable {
        id: appointments
        anchors.margins: parent.height * 0.02
        anchors.top: dayHeader.bottom
        anchors.left: parent.left
        anchors.right: parent.right
        anchors.bottom: parent.bottom

        contentWidth: width
        flickableDirection: Flickable.VerticalFlick
        clip: true

        property var startFirstAppointment
        property var endLastAppointment

        Rectangle {
            id: timeLine
            anchors.top: parent.top
            anchors.left: parent.left
            anchors.bottom: parent.bottom
            width: parent.width - dayHeader.width - dayHeader.anchors.margins
            color: "#D6E6FF"
            radius: 5

            Rectangle { // The red line
                property var secToPixRatio
                id: currentTime
                anchors.left: parent.left
                width: parent.width
                height: 2
                color: "#FF0000"

                function setCurrentTimeLine() {
                    let currTime = new Date()
                    if (typeof appointments.startFirstAppointment !== 'undefined' && currTime.getTime() > appointments.startFirstAppointment.getTime() && currTime.getTime() < appointments.endLastAppointment.getTime()) {
                        let offset = (currTime.getMillisecondsInDay() - appointments.startFirstAppointment.getMillisecondsInDay()) / 1000
                        offset *= secToPixRatio
                        currentTime.y = offset
                        currentTime.visible = true
                    } else {
                        currentTime.visible = false
                    }
                }
            }

            Timer {
                id: timeLineTimer
                interval: 1000 // 60 seconds
                repeat: true
                running: true
                triggeredOnStart: true
                onTriggered: currentTime.setCurrentTimeLine()
            }

            Component.onCompleted: currentTime.secToPixRatio = secondToPixelRatio
        }

        Component.onCompleted: appointments.visible = typeof appointments.startFirstAppointment !== 'undefined' // If there are no appointments, don't show the appointments and timeLine
    }
}
