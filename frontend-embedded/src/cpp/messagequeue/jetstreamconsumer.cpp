#include "jetstreamconsumer.h"
#include "natscallbackhandlersingleton.h"
#include "strings.h"

#include <QDebug>
#include <QtConcurrent/QtConcurrentRun>
#include <QtCore/QTimer>
#include <src/cpp/util/scopeguard.h>
#include <utility>

namespace MessageQueue
{

JetStreamConsumer::JetStreamConsumer(QObject *parent)
    : QObject(parent)
{
    connect(this, &JetStreamConsumer::updateSubscription, this, &JetStreamConsumer::setSubscription);
}

JetStreamConsumer::~JetStreamConsumer()
{
    if (m_sub != nullptr) {
        natsSubscription_Destroy(m_sub);

        NatsCallbackHandlerSingleton::singleton().removeMsgHandler(m_sub);
    }

    m_workerThread.quit();
    m_workerThread.wait();
}

NatsStatus::Enum JetStreamConsumer::lastStatus() const
{
    return m_lastStatus;
}

void JetStreamConsumer::updateStatus(NatsStatus::Enum s)
{
    if (s != m_lastStatus) {
        m_lastStatus = s;
        emit lastStatusChanged();
    }

    if (s == NatsStatus::Enum::Ok)
        return;

    const char *text = natsStatus_GetText(NatsStatus::asC(s));
    auto statusStr = QString::fromLocal8Bit(text);
    emit errorOccurred(s, statusStr);
}

QString JetStreamConsumer::subject() const
{
    return m_subject;
}

void JetStreamConsumer::setSubject(const QString &subject)
{
    if (subject != m_subject) {
        m_subject = subject;
        emit subjectChanged();
    }
}

NatsConnection *JetStreamConsumer::target() const
{
    return m_target;
}

void JetStreamConsumer::setTarget(NatsConnection *target)
{
    if (target != m_target) {
        m_target = target;
        emit targetChanged();
    }
}

void JetStreamConsumer::start()
{
    if (m_target == nullptr || m_sub != nullptr) return;

    QtConcurrent::run(
        [this](NatsConnection *target, const QString &topic, JetStreamConsumerType::Enum type, const QString &stream, const QString &consumer) {
            switch (type) {
            case JetStreamConsumerType::Push: {
                natsSubscription *pSub = nullptr;
                QSharedPointer<natsConnection *> dontDropConn;
                auto status = target->subscribe(topic, &pSub, dontDropConn);
                updateStatus(status);

                if (status == NatsStatus::Enum::Ok) {
                    auto ppSub = QSharedPointer<natsSubscription *>(new natsSubscription *(pSub));
                    emit updateSubscription(ppSub, dontDropConn, QPrivateSignal());
                }
                break;
            }
            case JetStreamConsumerType::Pull: {
                if (m_worker != nullptr) {
                    m_workerThread.quit();
                    m_workerThread.wait();
                }

                m_dontDropConn = m_target->getConnection();
                m_worker = new JetStreamPullConsumerWorker(m_dontDropConn, stream, consumer);
                m_worker->moveToThread(&m_workerThread);

                connect(&m_workerThread, &QThread::finished, m_worker, &QObject::deleteLater);
                connect(&m_workerThread, &QThread::started, m_worker, &JetStreamPullConsumerWorker::start);
                connect(m_worker, &JetStreamPullConsumerWorker::messageReceived, this, &JetStreamConsumer::handleMessageSP);

                m_workerThread.start();
            }
            }
        },
        m_target, m_subject, m_type, m_stream, m_consumer);
}

void JetStreamConsumer::setSubscription(
    const QSharedPointer<natsSubscription *> &sub,
    const QSharedPointer<natsConnection *> &spConn)
{
    if (m_sub != nullptr)
        natsSubscription_Destroy(m_sub);
    m_sub = *sub;

    // Call to clear is not really needed but useful for making the IDE think we're actually
    // using m_dontDropConn (which we are).
    m_dontDropConn.clear();
    m_dontDropConn = spConn;

    NatsCallbackHandlerSingleton::singleton().setMsgHandler(*sub, [this](natsMsg *msg) {
        qDebug() << "Emitting messageReceived for message on topic" << natsMsg_GetSubject(msg);
        emit handleMessage(msg);
        qDebug() << "Emitted messageReceived for message on topic" << natsMsg_GetSubject(msg);
    });
}

void JetStreamConsumer::handleMessage(natsMsg *msg)
{
    auto topic = QString::fromUtf8(natsMsg_GetSubject(msg));
    if (topic == "timeterm.disown-token") {
        timeterm_proto::messages::DisownTokenMessage m;

        if (m.ParseFromArray(natsMsg_GetData(msg), natsMsg_GetDataLength(msg)))
            handleDisownTokenProto(m);
    } else if (topic == "timeterm.retrieve-new-token") {
        timeterm_proto::messages::RetrieveNewTokenMessage m;

        if (m.ParseFromArray(natsMsg_GetData(msg), natsMsg_GetDataLength(msg)))
            handleRetrieveNewTokenProto(m);
    }
}

void JetStreamConsumer::handleDisownTokenProto(const timeterm_proto::messages::DisownTokenMessage &msg)
{
    DisownTokenMessage m;

    m.setDeviceId(QString::fromStdString(msg.device_id()));
    m.setTokenHash(QString::fromStdString(msg.token_hash()));
    m.setTokenHashAlg(QString::fromStdString(msg.token_hash_alg()));

    emit disownTokenMessage(m);
}

void JetStreamConsumer::handleRetrieveNewTokenProto(const timeterm_proto::messages::RetrieveNewTokenMessage &msg)
{
    RetrieveNewTokenMessage m;

    m.setDeviceId(QString::fromStdString(msg.device_id()));
    m.setCurrentTokenHash(QString::fromStdString(msg.current_token_hash()));
    m.setCurrentTokenHashAlg(QString::fromStdString(msg.current_token_hash_alg()));

    emit retrieveNewTokenMessage(m);
}

JetStreamConsumerType::Enum JetStreamConsumer::type() const
{
    return m_type;
}

void JetStreamConsumer::setType(JetStreamConsumerType::Enum type)
{
    if (type != m_type) {
        m_type = type;
        emit typeChanged();
    }
}

QString JetStreamConsumer::stream() const
{
    return m_stream;
}

void JetStreamConsumer::setStream(const QString &stream)
{
    if (stream != m_stream) {
        m_stream = stream;
        emit streamChanged();
    }
}

QString JetStreamConsumer::consumer() const
{
    return m_consumer;
}

void JetStreamConsumer::setConsumer(const QString &consumer)
{
    if (consumer != m_consumer) {
        m_consumer = consumer;
        emit consumerChanged();
    }
}

void JetStreamConsumer::handleMessageSP(const QSharedPointer<natsMsg *> &msg)
{
    handleMessage(*msg);
}

JetStreamPullConsumerWorker::JetStreamPullConsumerWorker(
    const QSharedPointer<natsConnection *> &conn,
    QString stream,
    QString consumer,
    QObject *parent)
    : QObject(parent)
    , m_timer(new QTimer(this))
    , m_conn(conn)
    , m_stream(std::move(stream))
    , m_consumer(std::move(consumer))
{
    m_timer.setInterval(0);
    connect(&m_timer, &QTimer::timeout, this, &JetStreamPullConsumerWorker::getNextMessage);
}

void JetStreamPullConsumerWorker::start()
{
    m_timer.start();
}

void JetStreamPullConsumerWorker::stop()
{
    m_timer.stop();
}

void JetStreamPullConsumerWorker::getNextMessage()
{
    // Don't fire too often in case of a timeout.
    m_timer.blockSignals(true);
    auto guard = onScopeExit([this]() {
        m_timer.blockSignals(false);
    });

    if (m_conn.isNull()) {
        qWarning() << "Not consuming next message, m_conn is null";
        return;
    }

    auto reply = QSharedPointer<natsMsg *>(
        new natsMsg *(nullptr),
        [](natsMsg **ppMsg) {
            if (ppMsg != nullptr) {
                if (*ppMsg != nullptr) {
                    natsMsg_Destroy(*ppMsg);
                }
                delete ppMsg;
            }
        });
    QString jsSubj = QString("$JS.API.CONSUMER.MSG.NEXT.%1.%2").arg(m_stream).arg(m_consumer);
    auto jsSubjCstr = asUtf8CString(jsSubj);
    auto status = natsConnection_RequestString(reply.get(), *m_conn, jsSubjCstr.get(), "", 1000);

    if (status != NATS_OK) {
        m_timer.setInterval(1000);

        if (status == NATS_TIMEOUT) {
            // No messages available, move along.
            return;
        }

        const char *err = nats_GetLastError(&status);
        qWarning() << "Could not request JetStream message:" << natsStatus_GetText(status) << "(detail:)" << err;
        nats_PrintLastErrorStack(stderr);
        return;
    }
    qDebug() << "Got message from NATS";

    emit messageReceived(reply);

    natsMsg *ackReply = nullptr;
    status = natsConnection_RequestString(&ackReply, *m_conn, natsMsg_GetReply(*reply), "", 1000);
    if (status != NATS_OK) {
        m_timer.setInterval(1000);

        qWarning() << "Could not acknowledge JetStream message:" << natsStatus_GetText(status);
    }
    natsMsg_Destroy(ackReply);

    m_timer.setInterval(0);
}

} // namespace MessageQueue