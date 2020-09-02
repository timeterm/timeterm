#include "stansubscription.h"
#include "stancallbackhandlersingleton.h"

#include <QDebug>

#include <QtConcurrent/QtConcurrentRun>
#include <timeterm_proto/messages.pb.h>

namespace MessageQueue
{

StanSubscription::StanSubscription(QObject *parent)
    : QObject(parent)
{
    connect(this, &StanSubscription::updateSubscription, this, &StanSubscription::setSubscription);
}

void StanSubscription::setSubscription(
    const QSharedPointer<stanSubscription *> &sub,
    const QSharedPointer<stanConnection *> &spConn)
{
    if (m_sub != nullptr)
        stanSubscription_Destroy(m_sub);
    m_sub = *sub;

    // Call to clear is not really needed but useful for making the IDE think we're actually
    // using m_dontDropConn (which we are).
    m_dontDropConn.clear();
    m_dontDropConn = spConn;

    StanCallbackHandlerSingleton::singleton().setMsgHandler(*sub, [this](const char *channel, stanMsg *msg) {
        qDebug() << "Emitting messageReceived for message on channel" << channel;
        emit handleMessage(QString::fromUtf8(channel), msg);
        qDebug() << "Emitted messageReceived for message on channel" << channel;
    });
}

StanSubscription::~StanSubscription()
{
    if (m_sub != nullptr) {
        stanSubscription_Destroy(m_sub);

        StanCallbackHandlerSingleton::singleton().removeMsgHandler(m_sub);
    }
}

StanSubOptions *StanSubscription::options() const
{
    return m_options;
}

void StanSubscription::setOptions(StanSubOptions *options)
{
    if (options != m_options) {
        if (m_options != nullptr)
            m_options->deleteLater();
        options->setParent(this);
        m_options = options;
        emit optionsChanged();
    }
}

void StanSubscription::setTarget(StanConnection *target)
{
    if (target != m_target) {
        m_target = target;
        emit targetChanged();
    }
}

StanConnection *StanSubscription::target() const
{
    return m_target;
}

NatsStatus::Enum StanSubscription::lastStatus() const
{
    return m_lastStatus;
}

void StanSubscription::updateStatus(NatsStatus::Enum s)
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

void StanSubscription::subscribe()
{
    if (m_target == nullptr || m_sub != nullptr) {
        return;
    }

    QtConcurrent::run(
        [this](StanConnection *target, StanSubOptions *options) {
            stanSubscription *pSub = nullptr;
            QSharedPointer<stanConnection *> dontDropConn;
            auto status = target->subscribe(options, &pSub, dontDropConn);
            updateStatus(status);

            if (status == NatsStatus::Enum::Ok) {
                auto ppSub = QSharedPointer<stanSubscription *>(new stanSubscription *(pSub));
                emit updateSubscription(ppSub, dontDropConn, QPrivateSignal());
            }
        },
        m_target, m_options);
}

void StanSubscription::handleMessage(const QString &channel, stanMsg *msg)
{
    if (channel == "timeterm.disown-token") {
        timeterm_proto::messages::DisownTokenMessage m;

        if (m.ParseFromArray(stanMsg_GetData(msg), stanMsg_GetDataLength(msg)))
            handleDisownTokenProto(m);
    } else if (channel == "timeterm.retrieve-new-token") {
        timeterm_proto::messages::RetrieveNewTokenMessage m;

        if (m.ParseFromArray(stanMsg_GetData(msg), stanMsg_GetDataLength(msg)))
            handleRetrieveNewTokenProto(m);
    }
}

void StanSubscription::handleDisownTokenProto(const timeterm_proto::messages::DisownTokenMessage &msg)
{
    DisownTokenMessage m;

    m.setDeviceId(QString::fromStdString(msg.device_id()));
    m.setTokenHash(QString::fromStdString(msg.token_hash()));
    m.setTokenHashAlg(QString::fromStdString(msg.token_hash_alg()));

    emit disownTokenMessage(m);
}

void StanSubscription::handleRetrieveNewTokenProto(const timeterm_proto::messages::RetrieveNewTokenMessage &msg)
{
    RetrieveNewTokenMessage m;

    m.setDeviceId(QString::fromStdString(msg.device_id()));
    m.setCurrentTokenHash(QString::fromStdString(msg.current_token_hash()));
    m.setCurrentTokenHashAlg(QString::fromStdString(msg.current_token_hash_alg()));

    emit retrieveNewTokenMessage(m);
}

} // namespace MessageQueue
