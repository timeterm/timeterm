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
    QString string = "Tuesday, 3 November 20 13:51:41";
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
    appointment.setSubjects({"nat"});
    appointment.setGroups({"gv6.nat1"});
    appointment.setLocations({"g208"});
    appointment.setTeachers({"mrd"});

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
