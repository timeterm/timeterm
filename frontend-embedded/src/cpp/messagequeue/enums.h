#ifndef ENUMS_H
#define ENUMS_H

#include <QObject>

namespace MessageQueue
{
Q_NAMESPACE

/// Status returned by most of the APIs
enum class NatsStatus
{
    /// Success
    Ok = 0,

    /// Generic error
    Err,

    /// Error when parsing a protocol message,
    /// or not getting the expected message.
    ProtocolError,

    /// IO Error (network communication).
    IoError,

    /// The protocol message read from the socket
    /// does not fit in the read buffer.
    LineTooLong,

    /// Operation on this connection failed because
    /// the connection is closed.
    ConnectionClosed,

    /// Unable to connect, the server could not be
    /// reached or is not running.
    NoServer,

    /// The server closed our connection because it
    /// did not receive PINGs at the expected interval.
    StaleConnection,

    /// The client is configured to use TLS, but the
    /// server is not.
    SecureConnectionWanted,

    /// The server expects a TLS connection.
    SecureConnectionRequired,

    /// The connection was disconnected. Depending on
    /// the configuration, the connection may reconnect.
    ConnectionDisconnected,

    /// The connection failed due to authentication error.
    ConnectionAuthFailed,

    /// The action is not permitted.
    NotPermitted,

    /// An action could not complete because something
    /// was not found. So far, this is an internal error.
    NotFound,

    /// Incorrect URL. For instance no host specified in
    /// the URL.
    AddressMissing,

    /// Invalid subject, for instance NULL or empty string.
    InvalidSubject,

    /// An invalid argument is passed to a function. For
    /// instance passing NULL to an API that does not
    /// accept this value.
    InvalidArg,

    /// The call to a subscription function fails because
    /// the subscription has previously been closed.
    InvalidSubscription,

    /// Timeout must be positive numbers.
    InvalidTimeout,

    /// An unexpected state, for instance calling
    /// #natsSubscription_NextMsg() on an asynchronous
    /// subscriber.
    IllegalState,

    /// The maximum number of messages waiting to be
    /// delivered has been reached. Messages are dropped.
    SlowConsumer,

    /// Attempt to send a payload larger than the maximum
    /// allowed by the NATS Server.
    MaxPayload,

    /// Attempt to receive more messages than allowed, for
    /// instance because of #natsSubscription_AutoUnsubscribe().
    MaxDeliveredMsg,

    /// A buffer is not large enough to accommodate the data.
    InsufficientBuffer,

    /// An operation could not complete because of insufficient
    /// memory.
    NoMemory,

    /// Some system function returned an error.
    SysError,

    /// An operation timed-out. For instance
    /// #natsSubscription_NextMsg().
    Timeout,

    /// The library failed to initialize.
    FailedToInitialize,

    /// The library is not yet initialized.
    NotInitialized,

    /// An SSL error occurred when trying to establish a
    /// connection.
    SslError,

    /// The server does not support this action.
    NoServerSupport,

    /// A connection could not be immediately established and
    /// #natsOptions_SetRetryOnFailedConnect() specified
    /// a connected callback. The connect is retried asynchronously.
    NotYetConnected,

    /// A connection and/or subscription entered the draining mode.
    /// Some operations will fail when in that mode.
    Draining,

    /// An invalid queue name was passed when creating a queue subscription.
    InvalidQueueName,
};

Q_ENUM_NS(NatsStatus)

} // namespace MessageQueue

#endif // ENUMS_H
