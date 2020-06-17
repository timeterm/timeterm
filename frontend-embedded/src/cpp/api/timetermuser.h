#ifndef TIMETERMUSER_H
#define TIMETERMUSER_H

#include <QObject>

class TimetermUser
{
    Q_GADGET
    Q_PROPERTY(QString cardUid READ cardUid WRITE setCardUid)
    Q_PROPERTY(QString organizationId READ organizationId WRITE setOrganizationId)
    Q_PROPERTY(QString name READ name WRITE setName)
    Q_PROPERTY(QString studentCode READ studentCode WRITE setStudentCode)

public:
    void setCardUid(const QString &cardUid);
    [[nodiscard]] QString cardUid() const;
    void setOrganizationId(const QString &organizationId);
    [[nodiscard]] QString organizationId() const;
    void setName(const QString &name);
    [[nodiscard]] QString name() const;
    void setStudentCode(const QString &studentCode);
    [[nodiscard]] QString studentCode() const;

    void read(const QJsonObject &json);
    void write(QJsonObject &json) const;

private:
    QString m_cardUid;
    QString m_organizationId;
    QString m_name;
    QString m_studentCode;
};

Q_DECLARE_METATYPE(TimetermUser)

#endif // TIMETERMUSER_H
