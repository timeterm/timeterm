#include <QGuiApplication>
#include <QQmlApplicationEngine>

#include "apiclient.h"
#include "cardreadercontroller.h"

int main(int argc, char *argv[])
{
    qputenv("QT_IM_MODULE", QByteArray("qtvirtualkeyboard"));

    QCoreApplication::setAttribute(Qt::AA_EnableHighDpiScaling);

    QGuiApplication app(argc, argv);
    QScopedPointer<CardReaderController> cardReader(new CardReaderController());

    qmlRegisterSingletonInstance("Timeterm.Rfid", 1, 0, "CardReader", cardReader.get());
    qmlRegisterUncreatableType<CardReaderController>("Timeterm.Rfid", 1, 0, "CardReaderController", "singleton");
    qmlRegisterType<ZermeloAppointment>("Timeterm.Zermelo", 1, 0, "ZermeloAppointment");
    qmlRegisterType<ZermeloAppointments>("Timeterm.Zermelo", 1, 0, "ZermeloAppointments");
    qmlRegisterType<ApiClient>("Timeterm.Api", 1, 0, "ApiClient");

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

    return QGuiApplication::exec();
}
