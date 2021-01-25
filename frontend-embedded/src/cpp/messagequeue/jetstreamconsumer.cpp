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
    , m_workerThread(new QThread(parent))
{
}

JetStreamConsumer::~JetStreamConsumer()
{
    m_workerThread->quit();
    m_workerThread->wait();
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

void JetStreamConsumer::start()
{
    if (m_connHolder.isNull()) {
        qCritical("JetStreamConsumer::start() called without a connection, not starting");
        return;
    }
    auto conn = m_connHolder;

    switch (m_type) {
    case JetStreamConsumerType::Push:
        qCritical("Push consumers are currently not supported");
        break;
    case JetStreamConsumerType::Pull:
        stop();
        auto worker = new JetStreamPullConsumerWorker(conn, m_stream, m_consumerId);
        worker->moveToThread(m_workerThread);

        connect(m_workerThread, &QThread::finished, worker, &QObject::deleteLater);
        connect(m_workerThread, &QThread::started, worker, &JetStreamPullConsumerWorker::start);
        connect(worker, &JetStreamPullConsumerWorker::messageReceived, this, &JetStreamConsumer::handleMessage);

        m_workerThread->start();
        break;
    }
}

void JetStreamConsumer::connectDecoder(Decoder *decoder) const
{
    connect(this, &JetStreamConsumer::messageReceived, decoder, &Decoder::decodeMessage);
}

void JetStreamConsumer::stop()
{
    if (m_workerThread->isRunning()) {
        m_workerThread->quit();
        m_workerThread->wait();
    }
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

void JetStreamConsumer::handleMessage(const QSharedPointer<natsMsg *> &msg)
{
    emit messageReceived(msg);
}

void JetStreamConsumer::useConnection(MessageQueue::NatsConnection *connection)
{
    if (!connection)
        stop();
    if (m_connHolder != connection->getHolder()) {
        stop();
        auto newHolder = connection->getHolder();
        m_connHolder.swap(newHolder);
    }
}

JetStreamPullConsumerWorker::JetStreamPullConsumerWorker(
    const QSharedPointer<NatsConnectionHolder> &connHolder,
    QString stream,
    QString consumerId,
    QObject *parent)
    : QObject(parent)
    , m_timer(new QTimer(this))
    , m_connHolder(connHolder)
    , m_stream(std::move(stream))
    , m_consumerId(std::move(consumerId))
{
    m_timer.setInterval(0);
    connect(&m_timer, &QTimer::timeout, this, &JetStreamPullConsumerWorker::getNextMessage);
}

JetStreamPullConsumerWorker::~JetStreamPullConsumerWorker()
{
    stop();
}

void JetStreamPullConsumerWorker::start()
{
    qDebug() << "Starting polling stream" << m_stream;
    m_timer.start();
}

void JetStreamPullConsumerWorker::stop()
{
    qDebug() << "Stopping polling stream" << m_stream;
    m_timer.stop();
}

void JetStreamPullConsumerWorker::getNextMessage()
{
    // Don't fire too often in case of a timeout.
    m_timer.blockSignals(true);
    auto guard = onScopeExit([this]() {
        m_timer.blockSignals(false);
    });

    if (m_connHolder.isNull() || !m_connHolder->getConnection()) {
        qWarning("Not consuming next message, m_connHolder or natsConnection is null");
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
    QString jsSubj = QString("$JS.API.CONSUMER.MSG.NEXT.%1.EMDEV-%2-%3").arg(m_stream).arg(m_consumerId).arg(m_stream);
    auto jsSubjCstr = asUtf8CString(jsSubj);
    auto status = natsConnection_RequestString(reply.get(), m_connHolder->getConnection(), jsSubjCstr.get(), "", 1000);

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
        return;
    }

    qDebug() << "Got JetStream reply with subject" << natsMsg_GetSubject(*reply);
    emit messageReceived(reply);

    // Acknowledge having received the message so JetStream doesn't redeliver it indefinitely.
    natsMsg *ackReply = nullptr;
    status = natsConnection_RequestString(&ackReply, m_connHolder->getConnection(), natsMsg_GetReply(*reply), "", 1000);
    if (status != NATS_OK) {
        // Something went wrong, wait a little bit before asking NATS for new messages.
        m_timer.setInterval(1000);

        qWarning() << "Could not acknowledge JetStream message:" << natsStatus_GetText(status);
        return;
    }
    if (ackReply != nullptr)
        natsMsg_Destroy(ackReply);

    // Everything went well, we can keep on consuming.
    m_timer.setInterval(0);
}

} // namespace MessageQueue
