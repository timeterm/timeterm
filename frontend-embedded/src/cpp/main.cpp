#include <QGuiApplication>
#include <QQmlApplicationEngine>

#include <src/cpp/messagequeue/enums.h>
#include <src/cpp/messagequeue/jetstreamconsumer.h>
#include <src/cpp/messagequeue/messages/disowntokenmessage.h>
#include <src/cpp/messagequeue/messages/retrievenewtokenmessage.h>
#include <src/cpp/messagequeue/natsconnection.h>
#include <src/cpp/messagequeue/natsoptions.h>
#include <src/cpp/messagequeue/natsstatusstringer.h>
#include <src/cpp/util/signalhandler.h>
#include <timeterm_proto/messages.pb.h>

#include "api/apiclient.h"
#include "cardreader/cardreadercontroller.h"

int runApp(int argc, char *argv[])
{
    QCoreApplication::setAttribute(Qt::AA_EnableHighDpiScaling);
    QGuiApplication app(argc, argv);

    QScopedPointer<CardReaderController> cardReader(new CardReaderController());
    auto natsStatusStringer = MessageQueue::NatsStatusStringer();

    qmlRegisterSingletonInstance("Timeterm.Rfid", 1, 0, "CardReaderController", cardReader.get());
    qmlRegisterUncreatableType<CardReaderController>("Timeterm.Rfid", 1, 0, "CardReaderControllerType", "singleton");
    qmlRegisterType<ApiClient>("Timeterm.Api", 1, 0, "ApiClient");
    qmlRegisterUncreatableMetaObject(MessageQueue::NatsStatus::staticMetaObject,
                                     "Timeterm.MessageQueue", 1, 0, "NatsStatus",
                                     "cannot create namespace NatsStatus in QML");
    qRegisterMetaType<MessageQueue::NatsStatus::Enum>();
    qmlRegisterUncreatableMetaObject(MessageQueue::JetStreamConsumerType::staticMetaObject,
                                     "Timeterm.MessageQueue", 1, 0, "JetStreamConsumerType",
                                     "cannot create namespace JetStreamConsumerType in QML");
    qRegisterMetaType<MessageQueue::JetStreamConsumerType::Enum>();
    qRegisterMetaType<QSharedPointer<natsConnection *>>();
    qRegisterMetaType<QSharedPointer<natsSubscription *>>();
    qRegisterMetaType<MessageQueue::DisownTokenMessage>();
    qRegisterMetaType<MessageQueue::RetrieveNewTokenMessage>();
    qmlRegisterType<MessageQueue::NatsOptions>("Timeterm.MessageQueue", 1, 0, "NatsOptions");
    qmlRegisterType<MessageQueue::NatsConnection>("Timeterm.MessageQueue", 1, 0, "NatsConnection");
    qmlRegisterType<MessageQueue::JetStreamConsumer>("Timeterm.MessageQueue", 1, 0, "JetStreamConsumer");
    qmlRegisterSingletonInstance("Timeterm.MessageQueue", 1, 0, "NatsStatusStringer", &natsStatusStringer);
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

    return teardownAppOnSignal<int>(QGuiApplication::exec);
}

int main(int argc, char *argv[])
{
    qputenv("QT_IM_MODULE", QByteArray("qtvirtualkeyboard"));
    qSetMessagePattern("%{time} %{type}%{if-category}:%{category}%{endif} [%{if-category}%{file}:%{endif}%{function}:%{line}]: %{message}");

    qInfo() << "Starting Timeterm frontend-embedded";

    qInfo() << "Verifying Protobuf library version...";
    GOOGLE_PROTOBUF_VERIFY_VERSION;
    qInfo() << "Protobuf library version OK";

    auto exitCode = runApp(argc, argv);
    qInfo() << "Shutting down...";

    qInfo() << "Shutting down Protobuf library...";
    google::protobuf::ShutdownProtobufLibrary();
    qInfo() << "Protobuf library shut down";

    qInfo() << "Closing NATS library (with timeout of 10s)...";
    nats_Sleep(500);
    nats_CloseAndWait(10000);
    qInfo() << "NATS library closed";

    qInfo("Exiting with code %d, have a nice day!", exitCode);
    return exitCode;
}
