#include <QGuiApplication>
#include <QQmlApplicationEngine>

#include "cardreadercontroller.h"
#include "mfrc522device.h"

int main(int argc, char *argv[])
{
    qputenv("QT_IM_MODULE", QByteArray("qtvirtualkeyboard"));

    QCoreApplication::setAttribute(Qt::AA_EnableHighDpiScaling);

    QGuiApplication app(argc, argv);
    QScopedPointer<CardReaderController> cardReader(new CardReaderController(
        CardReaderController::defaultCardReader()));

    qmlRegisterSingletonInstance("Timeterm.Rfid", 1, 0, "CardReader", cardReader.get());

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
