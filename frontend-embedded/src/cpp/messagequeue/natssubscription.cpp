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
    if (m_connHolder == nullptr) {
        qWarning() << "start() called with no set connection holder";
        return;
    }
    auto nc = m_connHolder->getConnection();
    if (!nc) {
        qWarning() << "NatsConnectionHolder contains no valid natsConnection";
        return;
    }
    if (m_sub) return;

    auto subj = m_subject.toStdString();

    qDebug() << "Subscribing to" << m_subject;
    auto subStatus = natsConnection_Subscribe(&m_sub, nc, subj.c_str(), &NatsSubscription::handleMessageReceived, this);
    updateStatus(NatsStatus::fromC(subStatus));
    if (subStatus != NATS_OK)
        return;
    qDebug() << "Subscribed to" << m_subject;
}

void NatsSubscription::stop()
{
    if (m_sub) {
        qDebug() << "Destroying subscription to" << m_subject;
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

void NatsSubscription::useConnection(NatsConnection *connection)
{
    if (!connection)
        stop();
    if (m_connHolder != connection->getHolder()) {
        stop();
        auto newHolder = connection->getHolder();
        m_connHolder.swap(newHolder);
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
