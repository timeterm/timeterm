#include "natscreds.h"
#include <QDir>

void NatsCredsResponse::read(const QJsonObject &json)
{
    if (json.contains("credentials") && json["credentials"].isString())
        credentials = json["credentials"].toString();
}

void NatsCredsResponse::writeToFile() const
{
    auto path = createNatsCredsPath();
    auto f = QFile(path);
    if (!f.open(QIODevice::WriteOnly | QIODevice::Truncate)) {
        qCritical() << "Could not open NATS credentials file (for writing)";
        return;
    }

    auto bytes = credentials.toUtf8();
    f.write(bytes);
    f.close();
}

QString createNatsCredsPath()
{
    auto filename = "EMDEV.creds";
    auto relative = QStringLiteral("nats/");

#if TIMETERMOS
    QString dir = "/opt/frontend-embedded/" + relative;
#else
    const QString &dir = relative;
#endif

    QDir(dir).mkpath(".");

    return dir + filename;
}
