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

    property var map

    background: Rectangle {
        color: "#FFFFFF"
    }

    function setTimetable(timetable) {
        map = new Map()

        weekAppointments.startFirstAppointment = null
        weekAppointments.endLastAppointment = null
        weekAppointments.contentHeight = 0

        for (var childCount = weekAppointments.contentItem.children.length; childCount > 0; childCount--) {
            if (weekAppointments.contentItem.children[childCount-1] instanceof WeekViewAppointment) {
                weekAppointments.contentItem.children[childCount-1].destroy()
            }
        }

        if (!startOfWeek || !endOfWeek) {
            startOfWeek = new Date().startOfWeek()
            endOfWeek = new Date().endOfWeek()
        }

        for (var i = 0; i < timetable.data.length; i++) {
            if (timetable.data[i].startTime.getTime() >= startOfWeek && timetable.data[i].endTime.getTime() < endOfWeek) {
                if (!weekAppointments.startFirstAppointment || timetable.data[i].startTime.getMillisecondsInDay() < weekAppointments.startFirstAppointment) {                                                          // first weekAppointment in the list
                    weekAppointments.startFirstAppointment = timetable.data[i].startTime.getMillisecondsInDay()
                }
                if (!weekAppointments.endLastAppointment || timetable.data[i].endTime.getMillisecondsInDay() > weekAppointments.endLastAppointment) {
                    weekAppointments.endLastAppointment = timetable.data[i].endTime.getMillisecondsInDay()
                }
            }
        }

        if (weekAppointments.startFirstAppointment && weekAppointments.endLastAppointment) {
            weekAppointments.contentHeight = (weekAppointments.endLastAppointment - weekAppointments.startFirstAppointment)
                                                / 1000 * weekPage.secondToPixelRatio - 5         // - 5 because of the spacing between weekAppointments
                                                + weekPage.height * 0.08
        }

        for (var i = 0; i < timetable.data.length; i++) {
            if (timetable.data[i].startTime.getTime() >= startOfWeek && timetable.data[i].endTime.getTime() < endOfWeek) {

                let finishWeekAppointment = function (weekAppointment) {
                    if (weekAppointment.status === Component.Ready) {
                        let incubator = weekAppointment.incubateObject(weekAppointments.contentItem, {
                            appointment: timetable.data[i],
                            startFirstAppointment: weekAppointments.startFirstAppointment,
                            secondToPixelRatio: weekPage.secondToPixelRatio,
                            weekAppointmentWidth: weekAppointments.weekAppointmentWidth
                        })
                        if (incubator.status !== Component.Ready) {
                            let aptId = timetable.data[i].id
                            incubator.onStatusChanged = function(status) {
                                if (status === Component.Ready) {
                                    //print ("Object", incubator.object, "is now ready!");
                                    map.set(aptId, incubator.object)
                                }
                            }
                        } else {
                            //print ("Object", incubator.object, "is ready immediately!");
                            map.set(timetable.data[i].id, incubator.object)
                        }
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
        if (!!weekAppointments.startFirstAppointment) {
            fillWeekTimeLine()
        }
        weekAppointments.visible = true
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
                    console.log("Could not create lineItem:", weekTimeLineItem.errorString())
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

    Timer {
        interval: 100
        running: true
        repeat: true
        onTriggered: {
            let idToDelete = map.keys().next().value
            if (idToDelete !== undefined) {
                console.log("Deleting: " + idToDelete)
                let itemToDelete = map.get(idToDelete)
                itemToDelete.destroy()
                map.delete(idToDelete)
            } else {
                console.log("Size of map: " + map.size)
            }
        }
    }
}
