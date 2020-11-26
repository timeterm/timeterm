import QtQuick 2.12
import QtQuick.Controls 2.12
import QtQml.Models 2.12
import QtQuick.Layouts 1.12
import QtGraphicalEffects 1.0

import "../js/TimeFunctions.js" as TimeFunction

Page {
    id: weekPage
    width: stack.width
    height: stack.height

    property int textSize: weekPage.height * 0.04
    property int customMargin: weekPage.height * 0.05
    property var secondToPixelRatio: weekAppointments.height * 0.000037
    property var startOfWeek
    property var endOfWeek

    background: Rectangle {
        color: "#FFFFFF"
    }

    function setTimetable(timetable) {
        for (var i = 0; i < timetable.data.length; i++) {
            if (typeof startOfWeek === 'undefined' || typeof endOfWeek === 'undefined') {
                startOfWeek = new Date().setHours(0, 0, 0, 0)
                endOfWeek = new Date().setHours(24, 0, 0, 0)
            }

            if (timetable.data[i].startTime.getTime() >= startOfWeek && timetable.data[i].endTime.getTime() < endOfWeek) {
                if (typeof weekAppointments.startFirstAppointment === 'undefined' || timetable.data[i].startTime.getMillisecondsInDay() < weekAppointments.startFirstAppointment.getMillisecondsInDay()) {                                                          // first weekAppointment in the list
                    weekAppointments.startFirstAppointment = timetable.data[i].startTime
                }
                if (typeof weekAppointments.endLastAppointment === 'undefined' || timetable.data[i].endTime.getMillisecondsInDay() > weekAppointments.endLastAppointment.getMillisecondsInDay()) {
                    weekAppointments.endLastAppointment = timetable.data[i].endTime
                    weekAppointments.contentHeight = (weekAppointments.endLastAppointment.getMillisecondsInDay()
                                                - weekAppointments.startFirstAppointment.getMillisecondsInDay())
                                                / 1000 * weekPage.secondToPixelRatio - 5         // - 5 because of the spacing between weekAppointments
                                                + weekPage.height * 0.08
                }

                let finishWeekAppointment = function (weekAppointment) {
                    if (weekAppointment.status === Component.Ready) {
                        weekAppointment.incubateObject(weekAppointments.contentItem, {
                            appointment: timetable.data[i],
                            startFirstAppointment: weekAppointments.startFirstAppointment,
                            secondToPixelRatio: weekPage.secondToPixelRatio,
                            weekAppointmentWidth: weekAppointments.weekAppointmentWidth
                        })
                    } else if (weekAppointment.status === Component.Error) {
                        console.log("Could not create weekAppointment:", weekAppointment.errorString())
                    }
                }

                let weekAppointment = Qt.createComponent("WeekViewAppointment.qml")
                if (weekAppointment.status !== Component.Null && weekAppointment.status !== Component.Loading) {
                    finishWeekAppointment(weekAppointment)
                } else {
                    weekAppointment.statusChanged.connect(finishWeekAppointment)
                }
            }
        }
        if (typeof weekAppointments.startFirstAppointment !== 'undefined') {
            fillWeekTimeLine()
        }
    }

    function fillWeekTimeLine() {
        let currentLineTime = new Date()
        currentLineTime.setTime(weekAppointments.startFirstAppointment.getTime())

        if (!currentLineTime.isFullHour()) {
            currentLineTime.addHours(1)
            currentLineTime.setMinutes(0, 0, 0)
        }

        for (currentLineTime;
             weekAppointments.endLastAppointment.getMillisecondsInDay() - currentLineTime.getMillisecondsInDay() > 30 * 60 * 1000; // Last weekTimeLine is at least less than 30 minutes before the end of the last weekAppointment
             currentLineTime.addHours(1)) {
            let finishWeekLineItem = function (weekTimeLineItem) {
                if (weekTimeLineItem.status === Component.Ready) {
                    weekTimeLineItem.incubateObject(weekTimeLine, {
                        y: (currentLineTime.getMillisecondsInDay() - weekAppointments.startFirstAppointment.getMillisecondsInDay()) / 1000 * weekPage.secondToPixelRatio,
                        time: currentLineTime.getHours().toString() + ":"
                              + (currentLineTime.getMinutes().toString() < 10 ? '0' : '')
                              + currentLineTime.getMinutes(),
                        textSize: weekPage.textSize
                    })
                } else if (weekTimeLineItem.status === Component.Error) {
                    console.log("Could not create lineItem:",
                                weekTimeLineItem.errorString())
                }
            }

            let weekTimeLineItem = Qt.createComponent("TimeLineItem.qml")
            if (weekTimeLineItem.status !== Component.Null && weekTimeLineItem.status !== Component.Loading) {
                finishWeekLineItem(weekTimeLineItem)
            } else {
                weekTimeLineItem.statusChanged.connect(finishWeekLineItem)
            }
        }
    }

    Flickable {
        id: weekAppointments
        anchors.margins: parent.height * 0.02
        anchors.fill: parent

        property var startFirstAppointment
        property var endLastAppointment
        property var weekAppointmentWidth: weekPage.width * 0.8 / 3.5

        contentWidth: weekPage.width * 0.15 + weekAppointmentWidth * 5 + weekPage.height * 0.02 * 4
        flickableDirection: Flickable.HorizontalAndVerticalFlick 
        clip: true

        Rectangle {
            id: weekTimeLine
            anchors.topMargin: weekPage.height * 0.08
            anchors.top: parent.top
            anchors.left: parent.left
            anchors.bottom: parent.bottom
            width: weekPage.width * 0.15
            color: "#D6E6FF"
            radius: 5

            Rectangle { // The red line
                property var secToPixRatio
                id: currentTime
                anchors.left: parent.left
                width: parent.width
                height: 2
                color: "#FF0000"

                function setCurrentweekTimeLine() {
                    let currTime = new Date()
                    if (typeof weekAppointments.startFirstAppointment !== 'undefined' && currTime.getTime() > weekAppointments.startFirstAppointment.getTime() && currTime.getTime() < weekAppointments.endLastAppointment.getTime()) {
                        let offset = (currTime.getMillisecondsInDay() - weekAppointments.startFirstAppointment.getMillisecondsInDay()) / 1000
                        offset *= secToPixRatio
                        currentTime.y = offset
                        currentTime.visible = true
                    } else {
                        currentTime.visible = false
                    }
                }
            }

            Timer {
                id: weekTimeLineTimer
                interval: 1000 // 60 seconds
                repeat: true
                running: true
                triggeredOnStart: true
                onTriggered: currentTime.setCurrentweekTimeLine()
            }

            Component.onCompleted: currentTime.secToPixRatio = weekPage.secondToPixelRatio
        }

        Component.onCompleted: weekAppointments.visible = typeof weekAppointments.startFirstAppointment !== 'undefined' // If there are no weekAppointments, don't show the weekAppointments and weekTimeLine
    }
}
