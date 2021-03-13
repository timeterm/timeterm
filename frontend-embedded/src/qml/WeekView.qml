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
    property int customMargin: weekPage.height * 0.02
    property var secondToPixelRatio: (weekAppointments.height - weekPage.height * 0.08) * 0.000037
    property var startOfWeek
    property var endOfWeek

    property var map: new Map()
    property var emptyChoiceAppointmentComponents: []
    property var currentAptMap: new Map()

    background: Rectangle {
        color: "#FFFFFF"
    }

    function setTimetable(timetable) {
        weekNumber.text = new Date().addDays(internals.dayOffset).getWeek()

        weekAppointments.startFirstAppointment = null
        weekAppointments.endLastAppointment = null
        weekAppointments.contentHeight = 0

        let newAptMap = new Map()
        let update = []

        for (let apt of timetable.data) {
            if (apt.id !== 0) newAptMap.set(apt.id, apt)
            else update.push(apt)
        }

        for (let apt of newAptMap.values()) {
            if (currentAptMap.has(apt.id)) {
                let old = currentAptMap.get(apt.id)
                if (!old.equals(apt)) update.push(apt)
            } else {
                update.push(apt)
            }
        }

        for (let id of currentAptMap.keys()) {
            if (!newAptMap.has(id)) {
                map.get(id).destroy()
                map.delete(id)
            }
        }

        for (let aptComponent of emptyChoiceAppointmentComponents) {
            aptComponent.destroy()
        }
        emptyChoiceAppointmentComponents = []

        if (!startOfWeek || !endOfWeek) {
            startOfWeek = new Date().startOfWeek()
            endOfWeek = new Date().endOfWeek()
        }

        for (let apt of timetable.data) {
            if (apt.startTime.getTime() >= startOfWeek && apt.endTime.getTime() < endOfWeek) {
                if (!weekAppointments.startFirstAppointment || apt.startTime.getMillisecondsInDay() < weekAppointments.startFirstAppointment) {                                                          // first weekAppointment in the list
                    weekAppointments.startFirstAppointment = apt.startTime.getMillisecondsInDay()
                }
                if (!weekAppointments.endLastAppointment || apt.endTime.getMillisecondsInDay() > weekAppointments.endLastAppointment) {
                    weekAppointments.endLastAppointment = apt.endTime.getMillisecondsInDay()
                }
            }
        }

        if (weekAppointments.startFirstAppointment && weekAppointments.endLastAppointment) {
            weekAppointments.contentHeight = (weekAppointments.endLastAppointment - weekAppointments.startFirstAppointment)
                                                / 1000 * weekPage.secondToPixelRatio - 5         // - 5 because of the spacing between weekAppointments
                                                + weekPage.height * 0.08
        }

        let finishWeekAppointment = function(weekAppointment) {
            if (weekAppointment.status === Component.Ready) {
                for (let apt of update) {
                    if (apt.startTime.getTime() >= startOfWeek && apt.endTime.getTime() < endOfWeek) {
                        let incubator = weekAppointment.incubateObject(weekAppointments.contentItem, {
                            appointment: apt,
                            startFirstAppointment: weekAppointments.startFirstAppointment,
                            secondToPixelRatio: weekPage.secondToPixelRatio,
                            weekAppointmentWidth: weekAppointments.weekAppointmentWidth
                        })

                        let aptId = apt.id
                        let finishIncubation = function(status) {
                            if (status === Component.Ready) {
                                if (aptId !== 0) { 
                                    if (map.has(aptId))
                                        map.get(aptId).destroy()
                                    map.set(aptId, incubator.object)
                                }
                                else emptyChoiceAppointmentComponents.push(incubator.object)
                            }
                        }

                        if (incubator.status !== Component.Ready) {
                            incubator.onStatusChanged = finishIncubation
                        } else {
                            finishIncubation(incubator.status)
                        }
                    }
                }

                if (!!weekAppointments.startFirstAppointment) {
                    fillWeekTimeLine()
                }
                weekAppointments.visible = true
                currentAptMap = newAptMap
            } else if (weekAppointment.status === Component.Error) {
                console.error("Could not create weekAppointment:", weekAppointment.errorString())
            }
        }

        let weekAppointment = Qt.createComponent("WeekViewAppointment.qml")
        if (weekAppointment.status !== Component.Null && weekAppointment.status !== Component.Loading) {
            finishWeekAppointment(weekAppointment)
        } else {
            weekAppointment.statusChanged.connect(finishWeekAppointment)
        }
    }

    function fillWeekTimeLine() {
        for (var childCount = weekTimeLine.children.length; childCount > 0; childCount--) {
            if (weekTimeLine.children[childCount-1] instanceof TimeLineItem) {
                weekTimeLine.children[childCount-1].destroy()
            }
        }

        let currentLineTime = new Date()
        currentLineTime.setTime(weekAppointments.startFirstAppointment)

        if (!currentLineTime.isFullHour()) {
            currentLineTime.addHours(1)
            currentLineTime.setMinutes(0, 0, 0)
        }

        for (currentLineTime;
             weekAppointments.endLastAppointment - currentLineTime.getMillisecondsInDay() > 30 * 60 * 1000; // Last weekTimeLine is at least less than 30 minutes before the end of the last weekAppointment
             currentLineTime.addHours(1)) {
            let finishWeekLineItem = function (weekTimeLineItem) {
                if (weekTimeLineItem.status === Component.Ready) {
                    weekTimeLineItem.incubateObject(weekTimeLine, {
                        y: (currentLineTime.getMillisecondsInDay() - weekAppointments.startFirstAppointment) / 1000 * weekPage.secondToPixelRatio,
                        time: currentLineTime.getHours().toString() + ":"
                              + (currentLineTime.getMinutes().toString() < 10 ? '0' : '')
                              + currentLineTime.getMinutes(),
                        textSize: weekPage.textSize
                    })
                } else if (weekTimeLineItem.status === Component.Error) {
                    console.error("Could not create lineItem:", weekTimeLineItem.errorString())
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
        visible: false // made visible if there are appointments to display

        property var startFirstAppointment
        property var endLastAppointment
        property var weekAppointmentWidth: weekPage.width * 0.85 / 3.5

        contentWidth: weekPage.width * 0.1 + weekAppointmentWidth * 5 + weekPage.height * 0.02 * 5
        flickableDirection: Flickable.HorizontalAndVerticalFlick 
        clip: true

        Rectangle {
            anchors.left: parent.left
            anchors.top: parent.top
            width: weekPage.width * 0.1
            height: weekPage.width * 0.04

            color: "#c4ffab"
            border.color: "#70ff33"
            border.width: 1
            radius: 5

            Button{
                anchors.left: parent.left
                anchors.top: parent.top
                anchors.bottom: parent.bottom
                width: height

                text: "<"
                font.pixelSize: textSize

                background: Rectangle {
                    color: "#00FFFFFF"
                }

                onClicked: function() {
                    internals.dayOffset -= 7
                    const startOfWeek = new Date().addDays(internals.dayOffset).startOfWeek()
                    const endOfWeek = new Date().addDays(internals.dayOffset).endOfWeek()
                    internals.getAppointments(startOfWeek, endOfWeek)
                }
            }

            Text {
                id: weekNumber
                anchors.centerIn: parent
                font.pixelSize: textSize

                Component.onCompleted: function() {
                    text = new Date().addDays(internals.dayOffset).getWeek()
                }
            }

            Button{
                anchors.right: parent.right
                anchors.top: parent.top
                anchors.bottom: parent.bottom
                width: height

                text: ">"
                font.pixelSize: textSize

                background: Rectangle {
                    color: "#00FFFFFF"
                }

                onClicked: function() {
                    internals.dayOffset += 7
                    const startOfWeek = new Date().addDays(internals.dayOffset).startOfWeek()
                    const endOfWeek = new Date().addDays(internals.dayOffset).endOfWeek()
                    internals.getAppointments(startOfWeek, endOfWeek)
                }
            }
        }

        Rectangle {
            id: weekTimeLine
            anchors.topMargin: weekPage.height * 0.08
            anchors.top: parent.top
            anchors.left: parent.left
            anchors.bottom: parent.bottom
            width: weekPage.width * 0.1
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

                    if (!!weekAppointments.startFirstAppointment
                            && !!weekAppointments.endLastAppointment
                            && currTime > startOfWeek
                            && currTime < endOfWeek
                            && currTime.getMillisecondsInDay() > weekAppointments.startFirstAppointment
                            && currTime.getMillisecondsInDay() < weekAppointments.endLastAppointment) {
                        let offset = (currTime.getMillisecondsInDay() - weekAppointments.startFirstAppointment) / 1000
                        offset *= secToPixRatio
                        currentTime.y = offset
                        if (currentTime.y > weekTimeLine.height) {
                            currentTime.visible = false
                        } else {
                            currentTime.visible = true
                        }
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

        Rectangle {
            id: monday
            anchors.top: parent.top
            x: weekTimeLine.width + weekPage.height * 0.02
            width: weekPage.width * 0.85 / 3.5
            height: weekPage.height * 0.06
            color: "#b5b5b5"
            radius: 5
            Text {
                text: "Maandag"
                anchors.verticalCenter: parent.verticalCenter
                anchors.centerIn: parent
                font.pixelSize: textSize
            }
        }

        Rectangle {
            id: tuesday
            anchors.top: parent.top
            x: weekTimeLine.width + weekPage.height * 0.02 * 2 + (weekPage.width * 0.85 / 3.5) * 1
            width: weekPage.width * 0.85 / 3.5
            height: weekPage.height * 0.06
            color: "#b5b5b5"
            radius: 5
            Text {
                text: "Dinsdag"
                anchors.verticalCenter: parent.verticalCenter
                anchors.centerIn: parent
                font.pixelSize: textSize
            }
        }

        Rectangle {
            id: wednesday
            anchors.top: parent.top
            x: weekTimeLine.width + weekPage.height * 0.02 * 3 + (weekPage.width * 0.85 / 3.5) * 2
            width: weekPage.width * 0.85 / 3.5
            height: weekPage.height * 0.06
            color: "#b5b5b5"
            radius: 5
            Text {
                text: "Woensdag"
                anchors.verticalCenter: parent.verticalCenter
                anchors.centerIn: parent
                font.pixelSize: textSize
            }
        }

        Rectangle {
            id: thursday
            anchors.top: parent.top
            x: weekTimeLine.width + weekPage.height * 0.02 * 4 + (weekPage.width * 0.85 / 3.5) * 3
            width: weekPage.width * 0.85 / 3.5
            height: weekPage.height * 0.06
            color: "#b5b5b5"
            radius: 5
            Text {
                text: "Donderdag"
                anchors.verticalCenter: parent.verticalCenter
                anchors.centerIn: parent
                font.pixelSize: textSize
            }
        }

        Rectangle {
            id: friday
            anchors.top: parent.top
            x: weekTimeLine.width + weekPage.height * 0.02 * 5 + (weekPage.width * 0.85 / 3.5) * 4
            width: weekPage.width * 0.85 / 3.5
            height: weekPage.height * 0.06
            color: "#b5b5b5"
            radius: 5
            Text {
                text: "Vrijdag"
                anchors.verticalCenter: parent.verticalCenter
                anchors.centerIn: parent
                font.pixelSize: textSize
            }
        }
    }
}
