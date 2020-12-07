#include "natssubscription.h"

#include <QDebug>

namespace MessageQueue
{

NatsSubscription::NatsSubscription(QObject *parent)
    : QObject(parent)
{}

NatsSubscription::~NatsSubscription()
{
    stop();
}

void NatsSubscription::start()
{
    if (m_conn == nullptr) {
        qWarning() << "start() called with null NatsConnection";
        return;
    }
    m_nc = m_conn->getConnection();
    if (m_nc.isNull()) {
        qWarning() << "NatsConnection has no nats.c connection";
        return;
    }
    if (m_sub != nullptr) return;

    auto subj = m_subject.toStdString();

    qDebug() << "Subscribing to" << m_subject;
    auto subStatus = natsConnection_Subscribe(&m_sub, *m_nc, subj.c_str(), &NatsSubscription::handleMessageReceived, this);
    updateStatus(NatsStatus::fromC(subStatus));
    if (subStatus != NATS_OK)
        return;
    qDebug() << "Subscribed to" << m_subject;
}

void NatsSubscription::stop()
{
    if (m_sub != nullptr) {
        natsSubscription_Destroy(m_sub);
        m_sub = nullptr;
    }
}

void NatsSubscription::handleMessageReceived(natsConnection *nc, natsSubscription *sub, natsMsg *msg, void *closure)
{
    static_cast<NatsSubscription *>(closure)->handleMessageReceived(msg);
}

// This function is very thread-unsafe! It gets called by NATS and so should probably only emit a signal
// notifying some other thread that a message has been received.
void NatsSubscription::handleMessageReceived(natsMsg *msg)
{
    auto spMsg = QSharedPointer<natsMsg *>(
        new natsMsg *(msg),
        [](natsMsg **ppMsg) {
            if (ppMsg != nullptr) {
                if (*ppMsg != nullptr) {
                    natsMsg_Destroy(*ppMsg);
                }
                delete ppMsg;
            }
        });

    emit messageReceived(spMsg);
}

void NatsSubscription::connectDecoder(Decoder *decoder) const
{
    connect(this, &NatsSubscription::messageReceived, decoder, &Decoder::decodeMessage);
}

QString NatsSubscription::subject() const
{
    return m_subject;
}

void NatsSubscription::setSubject(const QString &subject)
{
    if (subject != m_subject) {
        m_subject = subject;
        emit subjectChanged();
    }
}

NatsConnection *NatsSubscription::connection() const
{
    return m_conn;
}

void NatsSubscription::setConnection(NatsConnection *connection)
{
    if (connection != m_conn) {
        m_conn = connection;
        emit connectionChanged();
    }
}

NatsStatus::Enum NatsSubscription::lastStatus()
{
    return m_lastStatus;
}

void NatsSubscription::updateStatus(NatsStatus::Enum s)
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

} // namespace MessageQueue
