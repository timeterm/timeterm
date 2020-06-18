#include <QGuiApplication>
#include <QQmlApplicationEngine>
#include <src/cpp/messagequeue/enums.h>

#include "api/apiclient.h"
#include "cardreader/cardreadercontroller.h"

int main(int argc, char *argv[])
{
    qputenv("QT_IM_MODULE", QByteArray("qtvirtualkeyboard"));

    QCoreApplication::setAttribute(Qt::AA_EnableHighDpiScaling);

    QGuiApplication app(argc, argv);
    QScopedPointer<CardReaderController> cardReader(new CardReaderController());

    qmlRegisterSingletonInstance("Timeterm.Rfid", 1, 0, "CardReader", cardReader.get());
    qmlRegisterUncreatableType<CardReaderController>("Timeterm.Rfid", 1, 0, "CardReaderController", "singleton");
    qmlRegisterType<ApiClient>("Timeterm.Api", 1, 0, "ApiClient");
    qmlRegisterUncreatableMetaObject(MessageQueue::NatsStatus::staticMetaObject,
                                     "Timeterm.MessageQueue", 1, 0, "NatsStatus",
                                     "cannot create namespace NatsStatus in QML");
    qRegisterMetaType<MessageQueue::NatsStatus::Enum>("NatsStatusEnum");

    QQmlApplicationEngine engine;
    const QUrl url("qrc:/src/qml/main.qml");
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
