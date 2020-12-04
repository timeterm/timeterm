#pragma once

#include <QJsonObject>
#include <QObject>

class DeviceConfig: public QObject
{
    Q_OBJECT
    Q_PROPERTY(QString id READ id WRITE setId NOTIFY idChanged)
    Q_PROPERTY(QString name READ name WRITE setName NOTIFY nameChanged)
    Q_PROPERTY(QString setupToken READ setupToken WRITE setSetupToken NOTIFY setupTokenChanged)
    Q_PROPERTY(QString deviceToken READ deviceToken WRITE setDeviceToken NOTIFY deviceTokenChanged)
    Q_PROPERTY(QString deviceTokenSetupTokenHash READ deviceTokenSetupTokenHash WRITE setDeviceTokenSetupTokenHash NOTIFY deviceTokenSetupTokenHashChanged)
    Q_PROPERTY(bool needsRegistration READ needsRegistration)

public:
    explicit DeviceConfig(QObject *parent = nullptr);

    void setId(const QString &id);
    [[nodiscard]] QString id() const;
    void setName(const QString &name);
    [[nodiscard]] QString name() const;
    void setSetupToken(const QString &token);
    [[nodiscard]] QString setupToken() const;
    void setDeviceToken(const QString &token);
    [[nodiscard]] QString deviceToken() const;
    Q_INVOKABLE void setDeviceTokenSetupToken(const QString &token);
    void setDeviceTokenSetupTokenHash(const QString &hash);
    [[nodiscard]] QString deviceTokenSetupTokenHash() const;

    void write(QJsonObject &json) const;
    void read(const QJsonObject &json);

    bool needsRegistration();

signals:
    void idChanged();
    void nameChanged();
    void setupTokenChanged();
    void deviceTokenChanged();
    void deviceTokenSetupTokenHashChanged();

private:
    static QString hashToken(const QString &token);

    QString m_id;
    QString m_name;
    QString m_setupToken;
    QString m_deviceToken;
    QString m_deviceTokenSetupTokenHash;
};
