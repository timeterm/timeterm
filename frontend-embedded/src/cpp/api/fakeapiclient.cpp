#include "fakeapiclient.h"
#include "timetermuser.h"
#include "zermeloappointment.h"
#include "zermeloappointments.h"

#include <optional>
#include <utility>

#include <QDateTime>
#include <QJsonArray>
#include <QJsonObject>
#include <QJsonParseError>
#include <QNetworkReply>
#include <QString>
#include <QUrlQuery>

FakeApiClient::FakeApiClient(QObject *parent)
    : QObject(parent)
{
}

void FakeApiClient::getAppointments(const QDateTime &start, const QDateTime &end)
{
    //QDateTime testTime = QDateTime::currentDateTime();
    //testTime = testTime.addDays(-3);
    //testTime = testTime.addSecs(-8*60*60);

    auto appointments = ZermeloAppointments();
    auto appointment = ZermeloAppointment();

    for (int i = 0; i < 5; i++) {
        QDateTime testTime = QDateTime::fromMSecsSinceEpoch(1606722000000 + 86400000*i);
        appointment.setId(980803080);
        appointment.setAppointmentInstance(129304801);
        appointment.setStartTimeSlotName("1");
        appointment.setEndTimeSlotName("1");
        appointment.setStartTime(testTime);
        testTime = testTime.addSecs(45*60);
        appointment.setEndTime(testTime);
        appointment.setSubjects({"entl"});
        appointment.setGroups({"gv6.gv6a"});
        appointment.setLocations({"g028"});
        appointment.setTeachers({"dng"});
        appointments.append(appointment);

        appointment.setStartTimeSlotName("2");
        appointment.setEndTimeSlotName("2");
        appointment.setStartTime(testTime);
        testTime = testTime.addSecs(45*60);
        appointment.setEndTime(testTime);
        appointment.setSubjects({"nat"});
        appointment.setGroups({"gv6.nat1"});
        appointment.setLocations({"g208"});
        appointment.setTeachers({"mrd"});
        appointments.append(appointment);

        appointment.setStartTimeSlotName("3");
        appointment.setEndTimeSlotName("3");
        testTime = testTime.addSecs(15*60);
        appointment.setStartTime(testTime);
        testTime = testTime.addSecs(45*60);
        appointment.setEndTime(testTime);
        appointment.setSubjects({"to"});
        appointment.setGroups({"gv6.gv6b"});
        appointment.setLocations({"g045"});
        appointment.setTeachers({"mou"});
        appointment.setIsCanceled(true);
        appointments.append(appointment);

        appointment.setStartTimeSlotName("4");
        appointment.setEndTimeSlotName("4");
        appointment.setStartTime(testTime);
        testTime = testTime.addSecs(45*60);
        appointment.setEndTime(testTime);
        appointment.setSubjects({"gd"});
        appointment.setGroups({"gv6.gv6b"});
        appointment.setLocations({"g045"});
        appointment.setTeachers({"mou"});
        appointment.setIsCanceled(false);
        appointments.append(appointment);

        appointment.setStartTimeSlotName("5");
        appointment.setEndTimeSlotName("5");
        testTime = testTime.addSecs(45*60);
        appointment.setStartTime(testTime);
        testTime = testTime.addSecs(45*60);
        appointment.setEndTime(testTime);
        appointment.setSubjects({"z_uur"});
        appointment.setGroups({});
        appointment.setLocations({"g035"});
        appointment.setTeachers({});
        appointments.append(appointment);

        appointment.setStartTimeSlotName("6");
        appointment.setEndTimeSlotName("6");
        appointment.setStartTime(testTime);
        testTime = testTime.addSecs(45*60);
        appointment.setEndTime(testTime);
        appointment.setSubjects({"netl"});
        appointment.setGroups({"gv6.gv6b"});
        appointment.setLocations({"g137"});
        appointment.setTeachers({"knm"});
        appointments.append(appointment);

        appointment.setStartTimeSlotName("7");
        appointment.setEndTimeSlotName("8");
        testTime = testTime.addSecs(15*60);
        appointment.setStartTime(testTime);
        testTime = testTime.addSecs(90*60);
        appointment.setEndTime(testTime);
        appointment.setSubjects({"wisb"});
        appointment.setGroups({"gv6.wisb6"});
        appointment.setLocations({"g153"});
        appointment.setTeachers({"mlr"});
        appointments.append(appointment);
    }

    emit timetableReceived(appointments);
}

void FakeApiClient::getCurrentUser()
{
    auto user = TimetermUser();

    user.setName("TestUser");
    user.setStudentCode("12345");

    emit currentUserReceived(user);
}
