#pragma once

#include <QScopedArrayPointer>
#include <QString>

namespace MessageQueue
{

//! Converts a QString to a UTF-8 char array as a QScopedArrayPointer.
//! Automatically gets destroyed when it goes out of scope.
QScopedArrayPointer<char> asUtf8CString(const QString &str);

}
