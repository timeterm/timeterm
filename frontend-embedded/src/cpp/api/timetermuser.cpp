#include "timetermuser.h"

#include <QJsonObject>

void TimetermUser::setCardUid(const QString &cardUid)
{
    if (cardUid != m_cardUid) {
        m_cardUid = cardUid;
    }
}

QString TimetermUser::cardUid() const
{
    return m_cardUid;
}

void TimetermUser::setOrganizationId(const QString &organizationId)
{
    if (organizationId != m_organizationId) {
        m_organizationId = organizationId;
    }
}

QString TimetermUser::organizationId() const
{
    return m_organizationId;
}

void TimetermUser::setName(const QString &name)
{
    if (name != m_name) {
        m_name = name;
    }
}

QString TimetermUser::name() const
{
    return m_name;
}

void TimetermUser::setStudentCode(const QString &studentCode)
{
    if (studentCode != m_studentCode) {
        m_studentCode = studentCode;
    }
}

QString TimetermUser::studentCode() const
{
    return m_studentCode;
}

void TimetermUser::read(const QJsonObject &json)
{
    if (json.contains("cardUid") && json["cardUid"].isString())
        setCardUid(json["cardUid"].toString());

    if (json.contains("name") && json["name"].isString())
        setName(json["name"].toString());

    if (json.contains("organizationId") && json["organizationId"].isString())
        setOrganizationId(json["organizationId"].toString());

    if (json.contains("studentCode") && json["studentCode"].isString())
        setStudentCode(json["studentCode"].toString());
}

void TimetermUser::write(QJsonObject &json) const
{
    json["cardUid"] = cardUid();
    json["name"] = name();
    json["organizationId"] = organizationId();
    json["studentCode"] = studentCode();
}