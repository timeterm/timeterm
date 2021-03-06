#include "api/fakeapiclient.h"
#include "cardreader/cardreadercontroller.h"
#include "devcfg/connmanserviceconfig.h"
#include "messagequeue/enums.h"
#include "messagequeue/jetstreamconsumer.h"
#include "messagequeue/natsconnection.h"
#include "messagequeue/natsoptions.h"
#include "messagequeue/natsstatusstringer.h"
#include "util/teardown.h"

#include <QFontDatabase>
#include <QGuiApplication>
#include <QQmlApplicationEngine>

#include <api/apiclient.h>
#include <devcfg/configmanager.h>
#include <logs/logmanager.h>
#include <messagequeue/natssubscription.h>
#include <networking/networkmanager.h>
#include <systemd/systemd.h>
#include <timeterm_proto/mq/mq.pb.h>
#include <util/scopeguard.h>

void installDefaultFont()
{
    qint32 fontId = QFontDatabase::addApplicationFont(":/assets/fonts/Roboto/Roboto-Regular.ttf");
    QStringList fontList = QFontDatabase::applicationFontFamilies(fontId);

    const QString &family = fontList.at(0);
    QGuiApplication::setFont(QFont(family));
}

int runApp(int argc, char *argv[])
{
    QCoreApplication::setAttribute(Qt::AA_EnableHighDpiScaling);
    QGuiApplication app(argc, argv);

    QScopedPointer<CardReaderController> cardReader(new CardReaderController());
    auto natsStatusStringer = MessageQueue::NatsStatusStringer();

    qmlRegisterSingletonInstance("Timeterm.Rfid", 1, 0, "CardReaderController", cardReader.get());
    qmlRegisterUncreatableType<CardReaderController>("Timeterm.Rfid", 1, 0, "CardReaderControllerType", "singleton");
    qmlRegisterType<ApiClient>("Timeterm.Api", 1, 0, "ApiClient");
    qmlRegisterType<FakeApiClient>("Timeterm.Api", 1, 0, "FakeApiClient");
    qmlRegisterUncreatableMetaObject(MessageQueue::NatsStatus::staticMetaObject,
                                     "Timeterm.MessageQueue", 1, 0, "NatsStatus",
                                     "cannot create namespace NatsStatus in QML");
    qRegisterMetaType<MessageQueue::NatsStatus::Enum>();
    qmlRegisterUncreatableMetaObject(MessageQueue::JetStreamConsumerType::staticMetaObject,
                                     "Timeterm.MessageQueue", 1, 0, "JetStreamConsumerType",
                                     "cannot create namespace JetStreamConsumerType in QML");
    qRegisterMetaType<MessageQueue::JetStreamConsumerType::Enum>();
    qRegisterMetaType<QSharedPointer<natsMsg *>>();
    qRegisterMetaType<QSharedPointer<natsSubscription *>>();
    qRegisterMetaType<QSharedPointer<MessageQueue::NatsConnectionHolder>>();
    qmlRegisterType<MessageQueue::NatsOptions>("Timeterm.MessageQueue", 1, 0, "NatsOptions");
    qmlRegisterType<MessageQueue::NatsConnection>("Timeterm.MessageQueue", 1, 0, "NatsConnection");
    qmlRegisterType<MessageQueue::JetStreamConsumer>("Timeterm.MessageQueue", 1, 0, "JetStreamConsumer");
    qmlRegisterType<MessageQueue::NatsSubscription>("Timeterm.MessageQueue", 1, 0, "NatsSubscription");
    qmlRegisterType<MessageQueue::Decoder>("Timeterm.MessageQueue", 1, 0, "Decoder");
    qmlRegisterSingletonInstance("Timeterm.MessageQueue", 1, 0, "NatsStatusStringer", &natsStatusStringer);
    qmlRegisterUncreatableType<MessageQueue::NatsStatusStringer>("Timeterm.MessageQueue", 1, 0, "NatsStatusStringerType", "singleton");
    qmlRegisterType<ConfigManager>("Timeterm.Config", 1, 0, "ConfigManager");
    qmlRegisterSingletonType<QObject>("Timeterm.Logging", 1, 0, "LogManager", [](QQmlEngine *e, QJSEngine *se) {
        auto logMgr = LogManager::singleton();

        // The QML should not destroy the singleton as that would cause an invalid free.
        QQmlEngine::setObjectOwnership(logMgr, QQmlEngine::CppOwnership);

        return logMgr;
    });
    qmlRegisterType<NetworkManager>("Timeterm.Networking", 1, 0, "NetworkManager");
    qRegisterMetaType<NetworkState>();
    qmlRegisterType<Systemd>("Timeterm.Systemd", 1, 0, "Systemd");

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

    installDefaultFont();

    return tearDownAppOnSignal<int>(QGuiApplication::exec);
}

int main(int argc, char *argv[])
{
    qputenv("QT_IM_MODULE", QByteArray("qtvirtualkeyboard"));
    qSetMessagePattern("%{time yyyy-MM-dd HH:mm:ss.zzz} %{type}%{if-category}:%{category}%{endif} [%{if-category}%{file}:%{endif}%{function}:%{line}]: %{message}");
    qInstallMessageHandler(LogManager::handleMessage);

    qInfo() << "Starting Timeterm frontend-embedded";

    qInfo() << "Verifying Protobuf and NATS library versions...";
    GOOGLE_PROTOBUF_VERIFY_VERSION;
    nats_CheckCompatibility();
    qInfo() << "Protobuf and NATS library versions OK";

    qInfo() << "Initializing the NATS library...";
    // Open the NATS library manually so OpenSSL memory isn't double freed after nats_Close
    // by OpenSSL which seems to register an atexit() function (could also be Qt Networking though).
    nats_Open(-1);
    qInfo() << "NATS library initialized";

    auto exitCode = runApp(argc, argv);
    qInfo() << "Shutting down...";

    qInfo() << "Shutting down Protobuf library...";
    google::protobuf::ShutdownProtobufLibrary();
    qInfo() << "Protobuf library shut down";

    qInfo() << "Shutting down NATS library (with a timeout of 10s)...";
    nats_CloseAndWait(10000);
    qInfo() << "NATS library shut down";

    qInfo("Exiting with code %d, have a nice day!", exitCode);
    return exitCode;
}
