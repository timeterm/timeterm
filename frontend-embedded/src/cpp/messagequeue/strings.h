#ifndef STRINGS_H
#define STRINGS_H

#include <QScopedArrayPointer>
#include <QString>

namespace MessageQueue
{

QScopedArrayPointer<char> asUtf8CString(const QString &str);

}

#endif // STRINGS_H
