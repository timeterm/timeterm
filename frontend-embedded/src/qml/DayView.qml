import QtQuick 2.12
import QtQuick.Controls 2.12
import QtQml.Models 2.12
import QtQuick.Layouts 1.12

Page {
    id: dayPage
    anchors.fill: parent

    property int textSize: dayPage.height * 0.04
    property int customMargin: dayPage.height * 0.05

    background: Rectangle {
        color: "#FFFFFF"
    }

    function setTimetable(timetable) {
        // Pretty-print the timetable as JSON
        console.log()

        // Use it
        for (var i = 0; i < timetable.data.length; i++) {
            let appointment = Qt.createComponent("DayViewAppointment.qml")

            var incubator = appointment.incubateObject(appointments, {
                                                           "startTimeSlot": timetable.data[i].startTimeSlot,
                                                           "endTimeSlot": timetable.data[i].endTimeSlot,
                                                           "startTime": timetable.data[i].startTime,
                                                           "endTime": timetable.data[i].endTime,
                                                           "subjects": timetable.data[i].subjects,
                                                           "groups": timetable.data[i].groups,
                                                           "locations": timetable.data[i].locations,
                                                           "teachers": timetable.data[i].teachers
                                                       })
            if (incubator.status !== Component.Ready) {
                incubator.onStatusChanged = function (status) {
                    if (status === Component.Ready) {
                        print("Object", incubator.object, "is now ready!")
                    }
                }
            } else {
                print("Object", incubator.object, "is ready immediately!")
            }

            //            if (appointment.status === Component.Ready) {
            //                console.log("appointment ready")
            //                appointment.createObject(appointments.contentItem, {
            //                                             "startTimeSlot": timetable.data[i].startTimeSlot,
            //                                             "endTimeSlot": timetable.data[i].endTimeSlot,
            //                                             "startTime": timetable.data[i].startTime,
            //                                             "endTime": timetable.data[i].endTime,
            //                                             "subjects": timetable.data[i].subjects,
            //                                             "groups": timetable.data[i].groups,
            //                                             "locations": timetable.data[i].locations,
            //                                             "teachers": timetable.data[i].teachers
            //                                         })
            //                console.log("appointment inserted")
            //            } else {
            //                console.log("appointment wasn't ready")
            //            }
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
        z: 1
        Text {
            text: new Date().toLocaleString(Qt.locale("nl_NL"), "dddd")
            anchors.verticalCenter: parent.verticalCenter
            anchors.centerIn: parent
            font.pixelSize: textSize
        }
    }

    Flickable {
        id: appointments
        width: parent.width * 0.8
        anchors.margins: parent.height * 0.02
        anchors.top: dayHeader.bottom
        anchors.right: parent.right
        anchors.bottom: parent.bottom

        Rectangle {
            width: parent.width
            height: parent.height
            color: "#DDDDFF"
        }
    }

    //    ListView {
    //        id: dayViewList
    //        anchors.top: parent.top
    //        anchors.bottom: parent.bottom
    //        anchors.right: parent.right
    //        anchors.margins: parent.height * 0.02

    //        width: parent.width * 0.8
    //        header: dayHeader
    //        headerPositioning: ListView.OverlayHeader
    //        delegate: DayViewAppointment
    //        spacing: parent.height * 0.02
    //    }
}
