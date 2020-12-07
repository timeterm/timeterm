import QtQuick 2.12
import QtQuick.Controls 2.12
import QtQml.Models 2.12
import QtQuick.Layouts 1.12
import QtGraphicalEffects 1.0

import "../js/TimeFunctions.js" as TimeFunction

Page {
    id: dayPage
    width: stack.width
    height: stack.height

    property int textSize: dayPage.height * 0.04
    property int customMargin: dayPage.height * 0.05
    property var secondToPixelRatio: dayAppointments.height * 0.000037
    property var startOfDay
    property var endOfDay

    background: Rectangle {
        color: "#FFFFFF"
    }

    function setTimetable(timetable) {
        if (!startOfDay || !endOfDay) {
            startOfDay = new Date().setHours(0, 0, 0, 0)
            endOfDay = new Date().setHours(24, 0, 0, 0)
        }
        
        for (var i = 0; i < timetable.data.length; i++) {
            if (timetable.data[i].startTime.getTime() >= startOfDay && timetable.data[i].endTime.getTime() < endOfDay) {
                if (!dayAppointments.startFirstAppointment || timetable.data[i].startTime.getMillisecondsInDay() < dayAppointments.startFirstAppointment.getMillisecondsInDay()) {                                                          // first dayAppointment in the list
                    dayAppointments.startFirstAppointment = timetable.data[i].startTime
                }
                if (!dayAppointments.endLastAppointment || timetable.data[i].endTime.getMillisecondsInDay() > dayAppointments.endLastAppointment.getMillisecondsInDay()) {
                    dayAppointments.endLastAppointment = timetable.data[i].endTime
                    dayAppointments.contentHeight = (dayAppointments.endLastAppointment.getMillisecondsInDay()
                                                - dayAppointments.startFirstAppointment.getMillisecondsInDay())
                                                / 1000 * dayPage.secondToPixelRatio - 5         // - 5 because of the spacing between dayAppointments
                }

                let finishDayAppointment = function (dayAppointment) {
                    if (dayAppointment.status === Component.Ready) {
                        dayAppointment.incubateObject(dayAppointments.contentItem, {
                            appointment: timetable.data[i],
                            startFirstAppointment: dayAppointments.startFirstAppointment,
                            secondToPixelRatio: dayPage.secondToPixelRatio
                        })
                    } else if (dayAppointment.status === Component.Error) {
                        console.log("Could not create dayAppointment:", dayAppointment.errorString())
                    }
                }

                let dayAppointment = Qt.createComponent("DayViewAppointment.qml")
                if (dayAppointment.status !== Component.Null && dayAppointment.status !== Component.Loading) {
                    finishDayAppointment(dayAppointment)
                } else {
                    dayAppointment.statusChanged.connect(finishDayAppointment)
                }
            }
        }
        if (!!dayAppointments.startFirstAppointment) {
            fillDayTimeLine()
        }
        dayAppointments.visible = true
    }

    function fillDayTimeLine() {
        let currentLineTime = new Date()
        currentLineTime.setTime(dayAppointments.startFirstAppointment.getTime())

        if (!currentLineTime.isFullHour()) {
            currentLineTime.addHours(1)
            currentLineTime.setMinutes(0, 0, 0)
        }

        for (currentLineTime;
             dayAppointments.endLastAppointment.getMillisecondsInDay() - currentLineTime.getMillisecondsInDay() > 30 * 60 * 1000; // Last dayTimeLine is at least less than 30 minutes before the end of the last dayAppointment
             currentLineTime.addHours(1)) {
            let finishDayLineItem = function (dayTimeLineItem) {
                if (dayTimeLineItem.status === Component.Ready) {
                    dayTimeLineItem.incubateObject(dayTimeLine, {
                        y: (currentLineTime.getMillisecondsInDay() - dayAppointments.startFirstAppointment.getMillisecondsInDay()) / 1000 * dayPage.secondToPixelRatio,
                        time: currentLineTime.getHours().toString() + ":"
                              + (currentLineTime.getMinutes().toString() < 10 ? '0' : '')
                              + currentLineTime.getMinutes(),
                        textSize: dayPage.textSize
                    })
                } else if (dayTimeLineItem.status === Component.Error) {
                    console.log("Could not create lineItem:",
                                dayTimeLineItem.errorString())
                }
            }

            let dayTimeLineItem = Qt.createComponent("TimeLineItem.qml")
            if (dayTimeLineItem.status !== Component.Null && dayTimeLineItem.status !== Component.Loading) {
                finishDayLineItem(dayTimeLineItem)
            } else {
                dayTimeLineItem.statusChanged.connect(finishDayLineItem)
            }
        }
    }

    Rectangle {
        id: dayHeader
        width: parent.width - dayPage.width * 0.15 - dayPage.height * 0.06
        height: parent.height * 0.06
        anchors.top: parent.top
        anchors.right: parent.right
        anchors.margins: dayPage.height * 0.02
        color: "#b5b5b5"
        radius: 5
        Text {
            text: new Date().toLocaleString(Qt.locale("nl_NL"), "dddd")
            anchors.verticalCenter: parent.verticalCenter
            anchors.centerIn: parent
            font.pixelSize: textSize
            font.capitalization: Font.Capitalize
        }
    }

    Flickable {
        id: dayAppointments
        anchors.margins: dayPage.height * 0.02
        anchors.topMargin: dayPage.height * 0.1
        anchors.fill: parent
        visible: false // made visible if there are appointments to display

        property var startFirstAppointment
        property var endLastAppointment

        contentWidth: width
        flickableDirection: Flickable.VerticalFlick
        clip: true

        Rectangle {
            id: dayTimeLine
            anchors.top: parent.top
            anchors.left: parent.left
            anchors.bottom: parent.bottom
            width: dayPage.width * 0.15
            color: "#D6E6FF"
            radius: 5

            Rectangle { // The red line
                property var secToPixRatio
                id: currentTime
                anchors.left: parent.left
                width: parent.width
                height: 2
                color: "#FF0000"

                function setCurrentdayTimeLine() {
                    let currTime = new Date()
                    if (!!dayAppointments.startFirstAppointment && currTime > dayAppointments.startFirstAppointment && currTime < dayAppointments.endLastAppointment) {
                        let offset = (currTime.getMillisecondsInDay() - dayAppointments.startFirstAppointment.getMillisecondsInDay()) / 1000
                        offset *= secToPixRatio
                        currentTime.y = offset
                        currentTime.visible = true
                    } else {
                        currentTime.visible = false
                    }
                }
            }

            Timer {
                id: dayTimeLineTimer
                interval: 1000 // 60 seconds
                repeat: true
                running: true
                triggeredOnStart: true
                onTriggered: currentTime.setCurrentdayTimeLine()
            }

            Component.onCompleted: currentTime.secToPixRatio = dayPage.secondToPixelRatio
        }
    }
}
