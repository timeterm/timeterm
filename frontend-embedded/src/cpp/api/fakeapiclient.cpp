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
    QString string = "Tuesday, 3 November 20 08:40:00";
    QString format = "dddd, d MMMM yy hh:mm:ss";
    QDateTime testTime = QDateTime::fromString(string, format);

    auto appointments = ZermeloAppointments();

    auto appointment = ZermeloAppointment();

    appointment.setId(980803080);
    appointment.setAppointmentInstance(129304801);
    appointment.setStartTimeSlot("1");
    appointment.setEndTimeSlot("1");
    appointment.setStartTime(testTime);
    appointment.setEndTime(testTime.addSecs(45*60));
    appointment.setSubjects({"entl"});
    appointment.setGroups({"gv6.gv6a", "gv6.gv6b"});
    appointment.setLocations({"g028", "g029"});
    appointment.setTeachers({"dng"});
    appointments.append(appointment);

    appointment.setStartTimeSlot("2");
    appointment.setEndTimeSlot("2");
    appointment.setStartTime(testTime);
    appointment.setEndTime(testTime.addSecs(46*60));
    appointment.setSubjects({"nat"});
    appointment.setGroups({"gv6.nat1"});
    appointment.setLocations({"g208"});
    appointment.setTeachers({"mrd"});
    appointments.append(appointment);

    appointment.setStartTimeSlot("3");
    appointment.setEndTimeSlot("3");
    appointment.setStartTime(testTime.addSecs(15*60));
    appointment.setEndTime(testTime.addSecs(46*60));
    appointment.setSubjects({"to"});
    appointment.setGroups({"gv6.gv6b"});
    appointment.setLocations({"g045"});
    appointment.setTeachers({"mou"});
    appointments.append(appointment);

    appointment.setStartTimeSlot("4");
    appointment.setEndTimeSlot("4");
    appointment.setStartTime(testTime);
    appointment.setEndTime(testTime.addSecs(46*60));
    appointment.setSubjects({"gd"});
    appointment.setGroups({"gv6.gv6b"});
    appointment.setLocations({"g045"});
    appointment.setTeachers({"mou"});
    appointments.append(appointment);

    appointment.setStartTimeSlot("5");
    appointment.setEndTimeSlot("5");
    appointment.setStartTime(testTime.addSecs(45*60));
    appointment.setEndTime(testTime.addSecs(46*60));
    appointment.setSubjects({"z_uur"});
    appointment.setGroups({});
    appointment.setLocations({"g035"});
    appointment.setTeachers({});
    appointments.append(appointment);

    appointment.setStartTimeSlot("6");
    appointment.setEndTimeSlot("6");
    appointment.setStartTime(testTime);
    appointment.setEndTime(testTime.addSecs(46*60));
    appointment.setSubjects({"netl"});
    appointment.setGroups({"gv6.gv6b"});
    appointment.setLocations({"g137"});
    appointment.setTeachers({"knm"});
    appointments.append(appointment);

    appointment.setStartTimeSlot("7");
    appointment.setEndTimeSlot("7");
    appointment.setStartTime(testTime.addSecs(15*60));
    appointment.setEndTime(testTime.addSecs(46*60));
    appointment.setSubjects({"wisb"});
    appointment.setGroups({"gv6.wisb6"});
    appointment.setLocations({"g153"});
    appointment.setTeachers({"mlr"});
    appointments.append(appointment);

    emit timetableReceived(appointments);
}

void FakeApiClient::getCurrentUser()
{
    auto user = TimetermUser();

    user.setName("TestUser");
    user.setStudentCode("12345");

    emit currentUserReceived(user);
}
