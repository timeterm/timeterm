#include "strings.h"

namespace MessageQueue
{

QScopedArrayPointer<char> asUtf8CString(const QString &str)
{
    auto bytes = str.toUtf8();
    auto stdString = bytes.toStdString();
    auto src = stdString.c_str();

    auto dst = QScopedArrayPointer<char>(new char[strlen(src) + 1]);
    strcpy(dst.get(), src);

    // For some reason we can't move out the QScopedArrayPointer normally, so using a tiny hack.
    return QScopedArrayPointer<char>(dst.take());
}

}