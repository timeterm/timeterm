#pragma once

#include <QScopedArrayPointer>
#include <QString>

namespace MessageQueue
{

QScopedArrayPointer<char> asUtf8CString(const QString &str);

}
