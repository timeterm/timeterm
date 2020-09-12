#include "natssubscription.h"
#include "natscallbackhandlersingleton.h"
#include <QtConcurrent/QtConcurrentRun>

namespace MessageQueue
{

NatsSubscription::NatsSubscription(QObject *parent)
    : QObject(parent)
{
    connect(this, &NatsSubscription::updateSubscription, this, &NatsSubscription::setSubscription);
}

NatsSubscription::~NatsSubscription()
{
    if (m_sub != nullptr) {
        natsSubscription_Destroy(m_sub);

        NatsCallbackHandlerSingleton::singleton().removeMsgHandler(m_sub);
    }
}

NatsStatus::Enum NatsSubscription::lastStatus() const
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

QString NatsSubscription::topic() const
{
    return m_topic;
}

void NatsSubscription::setTopic(const QString &topic)
{
    if (topic != m_topic) {
        m_topic = topic;
        emit topicChanged();
    }
}

NatsConnection *NatsSubscription::target() const
{
    return m_target;
}

void NatsSubscription::setTarget(NatsConnection *target)
{
    if (target != m_target) {
        m_target = target;
        emit targetChanged();
    }
}

void NatsSubscription::subscribe()
{
    if (m_target == nullptr || m_sub != nullptr) return;

    QtConcurrent::run(
        [this](NatsConnection *target, const QString &topic) {
            natsSubscription *pSub = nullptr;
            QSharedPointer<natsConnection *> dontDropConn;
            auto status = target->subscribe(topic, &pSub, dontDropConn);
            updateStatus(status);

            if (status == NatsStatus::Enum::Ok) {
                auto ppSub = QSharedPointer<natsSubscription *>(new natsSubscription *(pSub));
                emit updateSubscription(ppSub, dontDropConn, QPrivateSignal());
            }
        },
        m_target, m_topic);
}

void NatsSubscription::setSubscription(
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

void NatsSubscription::handleMessage(natsMsg *msg)
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

void NatsSubscription::handleDisownTokenProto(const timeterm_proto::messages::DisownTokenMessage &msg)
{
    DisownTokenMessage m;

    m.setDeviceId(QString::fromStdString(msg.device_id()));
    m.setTokenHash(QString::fromStdString(msg.token_hash()));
    m.setTokenHashAlg(QString::fromStdString(msg.token_hash_alg()));

    emit disownTokenMessage(m);
}

void NatsSubscription::handleRetrieveNewTokenProto(const timeterm_proto::messages::RetrieveNewTokenMessage &msg)
{
    RetrieveNewTokenMessage m;

    m.setDeviceId(QString::fromStdString(msg.device_id()));
    m.setCurrentTokenHash(QString::fromStdString(msg.current_token_hash()));
    m.setCurrentTokenHashAlg(QString::fromStdString(msg.current_token_hash_alg()));

    emit retrieveNewTokenMessage(m);
}

} // namespace MessageQueue