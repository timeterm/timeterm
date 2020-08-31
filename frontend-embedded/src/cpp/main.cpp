#include <QGuiApplication>
#include <QQmlApplicationEngine>
#include <src/cpp/messagequeue/binaryprotoclient.h>
#include <src/cpp/messagequeue/enums.h>
#include <src/cpp/messagequeue/natsoptions.h>
#include <src/cpp/messagequeue/natsstatusstringer.h>
#include <src/cpp/messagequeue/stanconnection.h>
#include <src/cpp/messagequeue/stanconnectionoptions.h>

#include "api/apiclient.h"
#include "cardreader/cardreadercontroller.h"

int main(int argc, char *argv[])
{
    GOOGLE_PROTOBUF_VERIFY_VERSION;

    qputenv("QT_IM_MODULE", QByteArray("qtvirtualkeyboard"));

    QCoreApplication::setAttribute(Qt::AA_EnableHighDpiScaling);

    QGuiApplication app(argc, argv);
    QScopedPointer<CardReaderController> cardReader(new CardReaderController());
    QScopedPointer<MessageQueue::NatsStatusStringer> natsStatusStringer(new MessageQueue::NatsStatusStringer());

    qmlRegisterSingletonInstance("Timeterm.Rfid", 1, 0, "CardReaderController", cardReader.get());
    qmlRegisterUncreatableType<CardReaderController>("Timeterm.Rfid", 1, 0, "CardReaderControllerType", "singleton");
    qmlRegisterType<ApiClient>("Timeterm.Api", 1, 0, "ApiClient");
    qmlRegisterUncreatableMetaObject(MessageQueue::NatsStatus::staticMetaObject,
                                     "Timeterm.MessageQueue", 1, 0, "NatsStatus",
                                     "cannot create namespace NatsStatus in QML");
    qRegisterMetaType<MessageQueue::NatsStatus::Enum>();
    qRegisterMetaType<QSharedPointer<stanConnection*>>();
    qRegisterMetaType<MessageQueue::StanMessage>();
    qRegisterMetaType<MessageQueue::DisownTokenMessage>();
    qRegisterMetaType<MessageQueue::RetrieveNewTokenMessage>();
    qmlRegisterType<MessageQueue::BinaryProtoClient>("Timeterm.MessageQueue", 1, 0, "BinaryProtoClient");
    qmlRegisterType<MessageQueue::NatsOptions>("Timeterm.MessageQueue", 1, 0, "NatsOptions");
    qmlRegisterType<MessageQueue::StanConnection>("Timeterm.MessageQueue", 1, 0, "StanConnection");
    qmlRegisterType<MessageQueue::StanConnectionOptions>("Timeterm.MessageQueue", 1, 0, "StanConnectionOptions");
    qmlRegisterType<MessageQueue::StanSubOptions>("Timeterm.MessageQueue", 1, 0, "StanSubOptions");
    qmlRegisterType<MessageQueue::StanSubscription>("Timeterm.MessageQueue", 1, 0, "StanSubscription");
    qmlRegisterSingletonInstance("Timeterm.MessageQueue", 1, 0, "NatsStatusStringer", natsStatusStringer.get());
    qmlRegisterUncreatableType<MessageQueue::NatsStatusStringer>("Timeterm.MessageQueue", 1, 0, "NatsStatusStringerType", "singleton");

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

    auto exitCode = QGuiApplication::exec();

    google::protobuf::ShutdownProtobufLibrary();

    nats_Sleep(500);
    nats_Close();

    return exitCode;
}
