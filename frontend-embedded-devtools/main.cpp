#include "fakecardreaderclient.h"
#include <QGuiApplication>
#include <QQmlApplicationEngine>
#include <QQuickStyle>

int main(int argc, char *argv[])
{
    QCoreApplication::setAttribute(Qt::AA_EnableHighDpiScaling);

    QGuiApplication app(argc, argv);

    QQuickStyle::setStyle("Material");

    QQmlApplicationEngine engine;
    const QUrl url(QStringLiteral("qrc:/main.qml"));
    QObject::connect(
        &engine, &QQmlApplicationEngine::objectCreated,
        &app, [url](QObject *obj, const QUrl &objUrl) {
            if (!obj && url == objUrl)
                QCoreApplication::exit(-1);
        },
        Qt::QueuedConnection);
    engine.load(url);

    auto *window = engine.rootObjects().first()->findChild<QObject *>("button");
    auto *client = new FakeCardReaderClient(window);

    QObject::connect(window, SIGNAL(sendCardUid(QString, QString)), client, SLOT(sendCardUid(QString, QString)));

    return QGuiApplication::exec();
}
