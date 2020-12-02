#pragma once

#include <QString>
#include <QObject>
#include <QJsonObject>

class Device
{
    Q_GADGET
    Q_PROPERTY(QString id MEMBER id)
    Q_PROPERTY(QString organizationId MEMBER organizationId)
    Q_PROPERTY(QString name MEMBER name)

public:
    void read(const QJsonObject &json);

    QString id;
    QString organizationId;
    QString name;
};

bool operator==(const Device &a, const Device &b);
bool operator!=(const Device &a, const Device &b);

Q_DECLARE_METATYPE(Device)
