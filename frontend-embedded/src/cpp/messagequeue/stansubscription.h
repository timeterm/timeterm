#pragma once

#include <QObject>
#include <QSharedPointer>

#include "stanconnection.h"
#include "stansuboptions.h"

#include <src/cpp/messagequeue/messages/disowntokenmessage.h>
#include <src/cpp/messagequeue/messages/retrievenewtokenmessage.h>
#include <timeterm_proto/messages.pb.h>

namespace MessageQueue
{

using StanSubOptionsDeleter = ScopedPointerDestroyerDeleter<stanSubOptions, void, stanSubOptions_Destroy>;
using StanSubOptionsScopedPointer = QScopedPointer<stanSubOptions, StanSubOptionsDeleter>;

class StanSubscription: public QObject
{
    Q_OBJECT
    Q_PROPERTY(MessageQueue::StanConnection *target READ target WRITE setTarget NOTIFY targetChanged)
    Q_PROPERTY(MessageQueue::StanSubOptions *options READ options WRITE setOptions NOTIFY optionsChanged)
    Q_PROPERTY(MessageQueue::NatsStatus::Enum lastStatus READ lastStatus NOTIFY lastStatusChanged)

public:
    explicit StanSubscription(QObject *parent = nullptr);
    ~StanSubscription() override;

    [[nodiscard]] NatsStatus::Enum lastStatus() const;
    [[nodiscard]] StanSubOptions *options() const;
    void setOptions(StanSubOptions *subOpts);
    [[nodiscard]] StanConnection *target() const;
    void setTarget(StanConnection *target);

    Q_INVOKABLE void subscribe();

signals:
    void optionsChanged();
    void targetChanged();
    void lastStatusChanged();
    void errorOccurred(MessageQueue::NatsStatus::Enum s, const QString &msg);
    void disownTokenMessage(const MessageQueue::DisownTokenMessage &msg);
    void retrieveNewTokenMessage(const MessageQueue::RetrieveNewTokenMessage &msg);

    void updateSubscription(const QSharedPointer<stanSubscription *> &sub, const QSharedPointer<stanConnection*> &spConn, QPrivateSignal);

private slots:
    void setSubscription(const QSharedPointer<stanSubscription *>&sub, const QSharedPointer<stanConnection*> &spConn);

private:
    void updateStatus(NatsStatus::Enum s);
    void handleMessage(const QString &channel, stanMsg *msg);
    void handleRetrieveNewTokenProto(const timeterm_proto::messages::RetrieveNewTokenMessage &msg);
    void handleDisownTokenProto(const timeterm_proto::messages::DisownTokenMessage &msg);

    stanSubscription *m_sub = nullptr;
    QSharedPointer<stanConnection *> m_dontDropConn = nullptr;
    StanSubOptions *m_options = nullptr;
    StanConnection *m_target = nullptr;
    NatsStatus::Enum m_lastStatus = NatsStatus::Enum::Ok;
};

} // namespace MessageQueue
