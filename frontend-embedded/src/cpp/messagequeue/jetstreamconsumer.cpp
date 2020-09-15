#include "jetstreamconsumer.h"
#include "strings.h"
#include "util/scopeguard.h"

#include <utility>

#include <QDebug>
#include <QtConcurrent/QtConcurrentRun>
#include <QtCore/QTimer>

namespace MessageQueue
{

JetStreamConsumer::JetStreamConsumer(QObject *parent)
    : QObject(parent)
{
}

JetStreamConsumer::~JetStreamConsumer()
{
    if (m_sub != nullptr)
        natsSubscription_Destroy(m_sub);

    m_workerThread.quit();
    m_workerThread.wait();
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

    auto conn = m_target->getConnection();

    QtConcurrent::run(
        [this, conn](NatsConnection *target, JetStreamConsumerType::Enum type, const QString &stream, const QString &consumer) {
            switch (type) {
            case JetStreamConsumerType::Push:
                qCritical() << "Push consumers are currently not supported";
                break;
            case JetStreamConsumerType::Pull:
                if (m_workerThread.isRunning()) {
                    m_workerThread.quit();
                    m_workerThread.wait();
                }
                auto worker = new JetStreamPullConsumerWorker(conn, stream, consumer);
                worker->moveToThread(&m_workerThread);

                connect(&m_workerThread, &QThread::finished, worker, &QObject::deleteLater);
                connect(&m_workerThread, &QThread::started, worker, &JetStreamPullConsumerWorker::start);
                connect(worker, &JetStreamPullConsumerWorker::messageReceived, this, &JetStreamConsumer::handleMessageSP);

                m_workerThread.start();
                break;
            }
        },
        m_target, m_type, m_stream, m_consumerId);
}

void JetStreamConsumer::handleMessage(natsMsg *msg)
{
    auto subject = QString::fromUtf8(natsMsg_GetSubject(msg));
    qDebug() << "Handling message with subject" << subject;

    if (subject == QString("FEDEV.%1.DISOWN-TOKEN").arg(m_consumerId)) {
        timeterm_proto::messages::DisownTokenMessage m;

        if (m.ParseFromArray(natsMsg_GetData(msg), natsMsg_GetDataLength(msg)))
            handleDisownTokenProto(m);
    } else if (subject == QString("FEDEV.%1.RETRIEVE-NEW-TOKEN").arg(m_consumerId)) {
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

QString JetStreamConsumer::consumerId() const
{
    return m_consumerId;
}

void JetStreamConsumer::setConsumerId(const QString &consumerId)
{
    if (consumerId != m_consumerId) {
        m_consumerId = consumerId;
        emit consumerIdChanged();
    }
}

void JetStreamConsumer::handleMessageSP(const QSharedPointer<natsMsg *> &msg)
{
    handleMessage(*msg);
}

JetStreamPullConsumerWorker::JetStreamPullConsumerWorker(
    const QSharedPointer<natsConnection *> &conn,
    QString stream,
    QString consumerId,
    QObject *parent)
    : QObject(parent)
    , m_timer(new QTimer(this))
    , m_conn(conn)
    , m_stream(std::move(stream))
    , m_consumerId(std::move(consumerId))
{
    m_timer.setInterval(0);
    connect(&m_timer, &QTimer::timeout, this, &JetStreamPullConsumerWorker::getNextMessage);
}

JetStreamPullConsumerWorker::~JetStreamPullConsumerWorker()
{
    m_timer.stop();
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
    QString jsSubj = QString("$JS.API.CONSUMER.MSG.NEXT.%1.FEDEV-%2").arg(m_stream).arg(m_consumerId);
    auto jsSubjCstr = asUtf8CString(jsSubj);
    auto status = natsConnection_RequestString(reply.get(), *m_conn, jsSubjCstr.get(), "", 1000);

    if (status != NATS_OK) {
        // Something went wrong or there are no messages, wait a little bit before asking NATS for new messages.
        m_timer.setInterval(1000);
        // Unset the shared pointer so it doesn't try to free the NATS message which the NATS library
        // frees in case of a failure. When both try to free, we get a segmentation fault and that's not nice.
        *reply = nullptr;

        if (status == NATS_TIMEOUT) {
            // No messages available, move along.
            return;
        }

        const char *err = nats_GetLastError(&status);
        qWarning() << "Could not request JetStream message:" << natsStatus_GetText(status) << "(detail:)" << err;
        nats_PrintLastErrorStack(stderr);
        return;
    }

    emit messageReceived(reply);

    natsMsg *ackReply = nullptr;
    status = natsConnection_RequestString(&ackReply, *m_conn, natsMsg_GetReply(*reply), "", 1000);
    if (status != NATS_OK) {
        // Something went wrong, wait a little bit before asking NATS for new messages.
        m_timer.setInterval(1000);

        qWarning() << "Could not acknowledge JetStream message:" << natsStatus_GetText(status);
    }
    natsMsg_Destroy(ackReply);

    m_timer.setInterval(0);
}

} // namespace MessageQueue