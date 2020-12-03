#pragma once

#include <QObject>
#include <QJsonObject>

class DeviceInfo: public QObject
{
    Q_OBJECT
    Q_PROPERTY(QString id READ id WRITE setId NOTIFY idChanged)
    Q_PROPERTY(QString name READ name WRITE setName NOTIFY nameChanged)
    Q_PROPERTY(QString token READ token WRITE setToken NOTIFY tokenChanged)

public:
    explicit DeviceInfo(QObject *parent = nullptr);

    void setId(const QString &id);
    [[nodiscard]] QString id() const;
    void setName(const QString &name);
    [[nodiscard]] QString name() const;
    void setToken(const QString &token);
    [[nodiscard]] QString token() const;

    void write(QJsonObject &json) const;
    void read(const QJsonObject &json);

signals:
    void idChanged();
    void nameChanged();
    void tokenChanged();

private:
    QString m_id;
    QString m_name;
    QString m_token;
};
