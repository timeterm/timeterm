#pragma once

#include "device.h"
#include <QJsonObject>
#include <QObject>
#include <QString>

class CreateDeviceRequest
{
    Q_GADGET
    Q_PROPERTY(QString name MEMBER name)

public:
    void write(QJsonObject &json) const;

    QString name;
};

class CreateDeviceResponse {
    Q_GADGET
    Q_PROPERTY(Device device MEMBER device)
    Q_PROPERTY(QString token MEMBER token)

public:
    void read(const QJsonObject &json);

    Device device;
    QString token;
};

Q_DECLARE_METATYPE(CreateDeviceRequest)
Q_DECLARE_METATYPE(CreateDeviceResponse)
